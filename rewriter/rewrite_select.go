package rewriter

import (
	"go/ast"
)

// golang spec
// https://go.dev/ref/spec#SelectStmt
// A case with a RecvStmt may assign the result of a RecvExpr to one or two variables,
// which may be declared using a short variable declaration.
// The RecvExpr must be a (possibly parenthesized) receive operation.
// There can be at most one default case and it may appear anywhere in the list of cases.
//
// Execution of a "select" statement proceeds in several steps:
//
//  1. For all the cases in the statement, the channel operands of receive operations and the channel
//     and right-hand-side expressions of send statements are evaluated exactly once, in source order,
//     upon entering the "select" statement. The result is a set of channels to receive from or send to,
//     and the corresponding values to send. Any side effects in that evaluation will occur irrespective
//     of which (if any) communication operation is selected to proceed. Expressions on the left-hand side
//     of a RecvStmt with a short variable declaration or assignment are not yet evaluated.
//  2. If one or more of the communications can proceed, a single one that can proceed is chosen via a
//     uniform pseudo-random selection. Otherwise, if there is a default case, that case is chosen.
//     If there is no default case, the "select" statement blocks until at least one of the communications
//     can proceed.
//  3. Unless the selected case is the default case, the respective communication operation is executed.
//  4. If the selected case is a RecvStmt with a short variable declaration or an assignment,
//     the left-hand side expressions are evaluated and the received value (or values) are assigned.
//  5. The statement list of the selected case is executed.
//
// RecvStmt 可以包含 一个 AssignStmt
// 最多一个 default, 位置无所谓 (跟 switch 规则一样)
// 进入 select 语句时, 每个分支的
//  1. recv 的 rhs
//  2. send 语句 lhs, rhs
//
// 会按源码顺序求值仅且一次
// 其中, recv 的 lhs 只有在分支被选中时才会求值
// 注意: select case recv 的 assign 跟普通 assign 的求值顺序是不一样
// assign stmt 的 lhs 会在 rhs 之前求值,
// 而 select case recv 的 lhs 只有在分支被选中时才会求值 (可以看成先 rhs 再 lhs)
//
// e.g.
//
//	select {
//		case *getvar(&i1) = <-getRecvCh("c1", 1):
//		case mkSndCh() <- mkValToSend():
//		case a[f()] = <-getRecvCh("c2", 0):
//		default:
//		}
//
// output:
//
//	getRecvCh called: c1
//	mkSndCh called
//	mkValToSend called
//	getRecvCh called: c2
//	getvar called // 选中第一个分支执行
func (r *fileRewriter) rewriteSelect(ctx *rewriteCtx, n *ast.SelectStmt) (ast.Stmt, []ast.Stmt) {
	var xs, ys []ast.Stmt
	n1 := &ast.SelectStmt{
		Body: &ast.BlockStmt{
			List: make([]ast.Stmt, len(n.Body.List)),
		},
	}
	// 额外生成的 stmt 需要提取到 select 的外头,
	// 而不是合并到 n.Body.List 中, 因为 select 的 body 中 只允许 CommClause
	// 所以不能用复用 rewriteStmtList
	// n1.Body.List = r.rewriteStmtList(ctx, n.Body.List)
	for i, it := range n.Body.List { // CommClauses only
		n1.Body.List[i], xs = r.rewriteStmt(ctx.enter(it), it)
		ys = append(ys, xs...)
	}
	return n1, ys
}

// CommCase   = "case" ( SendStmt | RecvStmt ) | "default" .
// RecvStmt   = [ ExpressionList "=" | IdentifierList ":=" ] RecvExpr .
// SendStmt   = Channel "<-" Expression .
// RecvExpr   = Expression .
// Channel    = Expression .
func (r *fileRewriter) rewriteCommClause(ctx *rewriteCtx, n *ast.CommClause) (ast.Stmt, []ast.Stmt) {
	var xs []ast.Stmt
	switch s := n.Comm.(type) {
	case nil: // default branch
		n1 := &ast.CommClause{}
		n1.Body = r.rewriteStmtList(ctx, n.Body)
		return n1, nil
	case *ast.SendStmt:
		// send 先对 lhs 求值, 再对 rhs 求值
		// 参见 处理 *ast.SendStmt 分支
		n1 := &ast.CommClause{}
		n1.Comm, xs = r.rewriteStmt(ctx, s)
		n1.Body = r.rewriteStmtList(ctx, n.Body)
		return n1, xs
	case *ast.ExprStmt: // recv expr
		// 非 assign 的 recv, 直接对 rhs 求值
		assert(isRecvOperation(s.X))
		n1 := &ast.CommClause{}
		n1.Comm, xs = r.rewriteStmt(ctx, s)
		n1.Body = r.rewriteStmtList(ctx, n.Body)
		return n1, xs
	case *ast.AssignStmt: // assign recv stmt
		// golang spec
		// The RecvExpr must be a (possibly parenthesized) receive operation.
		assert(len(s.Rhs) == 1 && isRecvOperation(s.Rhs[0]))
		n1 := &ast.CommClause{}
		var insertBefore []ast.Stmt
		n1.Comm, insertBefore, xs = r.rewriteAssignInRecvCommLhs(ctx, s)
		// 都是生成的 sym, 没有 shadow 的问题, 不用套 block
		// n1.Body = concat(insertBefore, sliceOf[ast.Stmt](&ast.BlockStmt{ List: r.rewriteStmtList(ctx, body) }))
		n1.Body = concat(insertBefore, r.rewriteStmtList(ctx, n.Body))
		return n1, xs
	default:
		panic("clause.Comm must be a SendStmt, RecvStmt, or default case")
	}
}

func (r *fileRewriter) rewriteAssignInRecvCommLhs(ctx *rewriteCtx, n *ast.AssignStmt) (*ast.AssignStmt, []ast.Stmt, []ast.Stmt) {
	var xs []ast.Stmt
	n1 := &ast.AssignStmt{
		Lhs: n.Lhs[:],
		Tok: n.Tok,
	}

	// 前面解释过 comm 中 assign 与 普通 assign 求值顺序不同,
	// 进入 select , rhs 先求值, 分支命中时 lhs 才求值
	// 所以, rhs 跟 lhs 的求值顺序无关
	// 所以, 不能用 r.tryNodes[n/*assign*/] 来标记 ctx 中 try 状态
	// 假设 lhs 有 try, rhs 没有 try
	// ctx = ctx.withTry(r.tryNodes[n]) 则会导致 rhs 多余的展开
	// 其次, 按求值顺序改写可以保证生成的符号 id 符合书写顺序
	// 所以先改写 n.Rhs, 再改写 n.Lhs
	tryInRhs := findFirst(n.Rhs, func(it ast.Expr) bool { return r.tryNodes[it] }) != nil
	n1.Rhs, xs = r.rewriteExprList(ctx.withTry(tryInRhs), n.Rhs)

	// 参见 rewriteTuple2Assign 注释

	recv, recvOK := unpackTestTupleAssign(n.Lhs)
	vGen := func() *ast.Ident { return r.genValId(ctx.fun) }
	bGen := func() *ast.Ident { return r.genSym(ctx.fun, "ok") }
	insertBefore := r.rewriteTuple2Assign(ctx, &recv, &recvOK, &n1.Tok, vGen, bGen)
	n1.Lhs[0] = recv
	if recvOK != nil {
		n1.Lhs[1] = recvOK
	}
	return n1, insertBefore, xs
}

// 三种可选返回 bool 值的特殊情况:
// channel recv, map index, type assert
// x, ok = ...
func unpackTestTupleAssign(lhs []ast.Expr) (v ast.Expr, ok ast.Expr) {
	// v [:]= <-ch
	// v, ok [:]= <-ch
	assert(len(lhs) == 1 || len(lhs) == 2)
	v = lhs[0]
	if len(lhs) == 2 {
		ok = lhs[1]
	}
	return
}
