package rewriter

import (
	"go/ast"
	"go/token"
)

// golang spec
// https://go.dev/ref/spec#Assignment_statements
// A tuple assignment assigns the individual elements of a multi-valued operation to a list of variables.
// There are two forms.
//
// In the first, the right hand operand is a single multi-valued expression
// such as a function call, a channel or map operation, or a type assertion.
//
// In the second form, the number of operands on the left must equal the number
// of expressions on the right, each of which must be single-valued, and the nth
// expression on the right is assigned to the nth operand on the left
//
// 1. 第一种, rhs 必须有且仅有一个多值表达式: 函数调用, channel recv, map index, 类型断言
// 		x, err = f()
// 		v, ok = a[x]
// 		x, ok = <-ch
// 		x, ok = a.(T)
//		其中, 赋值可以是 define/assign/valueSpec 形式, 或者可以说是短定义/赋值/声明
// 2. rhs 和 lhs 数量相等, 类型匹配
//
// The assignment proceeds in two phases.
// First, the operands of index expressions and pointer indirections (including implicit pointer indirections in selectors)
// on the left and the expressions on the right are all evaluated in the usual order.
// Second, the assignments are carried out in left-to-right order.
//
// 三种可选返回 bool 值的特殊情况: channel recv, map index, type assert
//
// An index expression on a map a of type map[K]V used in an assignment statement or initialization of the special form
// 		v, ok = a[x]
// 		v, ok := a[x]
// 		var v, ok = a[x]
// yields an additional untyped boolean value. The value of ok is true if the key x is present in the map, and false otherwise.
// Assigning to an element of a nil map causes a run-time panic.
//
// A type assertion used in an assignment statement or initialization of the special form
// 		v, ok = x.(T)
// 		v, ok := x.(T)
// 		var v, ok = x.(T)
// 		var v, ok interface{} = x.(T) // dynamic types of v and ok are T and bool
// yields an additional untyped boolean value. The value of ok is true if the assertion holds.
// Otherwise it is false and the value of v is the zero value for type T. No run-time panic occurs in this case.
//
// A receive expression used in an assignment statement or initialization of the special form
// 		x, ok = <-ch
// 		x, ok := <-ch
// 		var x, ok = <-ch
// 		var x, ok T = <-ch
// yields an additional untyped boolean result reporting whether the communication succeeded.
// The value of ok is true if the value received was delivered by a successful send operation to the channel,
// or false if it is a zero value generated because the channel is closed and empty.

func (r *fileRewriter) rewriteAssign(ctx *rewriteCtx, n *ast.AssignStmt) (ast.Stmt, []ast.Stmt) {
	var xs, ys []ast.Stmt
	n1 := &ast.AssignStmt{
		Lhs: n.Lhs,
		Tok: n.Tok,
	}

	// 左值也可以包含 try
	// e.g. Try(ret1Err[*X]()).v = 1
	//
	// 求值顺序: 先左值再右值
	// ---------------------------
	// type X struct{ x string }
	//
	//	lhs := func() *X {
	//		println("lhs")
	//		return &X{}
	//	}
	//	rhs := func() string {
	//		println("rhs")
	//		return "hello"
	//	}
	//	lhs().x = rhs()
	//
	// ---------------------------
	// output:
	//	lhs
	//	rhs
	// ---------------------------
	// lhs 与 rhs 之间以及各自内部元素都存在求值顺序依赖
	// lhs 任何一个元素包含 try 或者 rhs 包含 try, 则都要展开, 满足求值顺序
	// (其实可以按求值顺序处理到最后一个 try 停止)
	// ---------------------------
	// 三个可选返回 bool 的特殊情况在各自 rewrite 函数内部处理
	// 通过 parent node 是否是赋值语句进行判断

	ctx = ctx.withTry(r.tryNodes[n])
	// lhs & rhs 之前有求值顺序依赖, lhs+rhs 需要放到一起裁剪定义提取, 这里只处理了 rhs
	n1.Lhs, xs = r.rewriteExprList(ctx.withLhs(), n.Lhs)
	n1.Rhs, ys = r.rewriteExprList(ctx, n.Rhs)
	return n1, append(xs, ys...)
}

// 改写 range 和 select case recv 的 tupleAssign
func (r *fileRewriter) rewriteTuple2Assign(
	ctx *rewriteCtx,
	outFst, outSndOpt *ast.Expr,
	outTok *token.Token,
	fstGen, sndGen func() *ast.Ident,
) []ast.Stmt {
	oriTok := *outTok
	fst, snd := *outFst, *outSndOpt

	// <<对于 range 而言>>
	// =================================================================
	// 副作用: 包含 try 调用, 或者可能 panic
	// 注意: 因为求值顺序, k v 中如果包含副作用需要展开到 body
	// 否则会出现把 k 展开到 body, v 中副作用(比如函数调用, 比如可能 panic 的 selector/index) 没展开
	// 从而导致, 每次循环会先求值 v 的左值, 然后再求值 k 的左值, 打破了求值顺序的语义
	// 所以, 只要 v 需要展开, 在 v 之前求值的 k (无论是否有副作用) 也必须展开到 body
	// 如果 v 没有副作用, 可以只展开 k(有副作用), 这种情况提前求值 v 不影响语义
	// RangeStmt.Token 的处理
	// e.g.
	// var v int
	// for try(f()).k, v = range []int{} { }
	// ---------------------------
	// var v int
	// var k int // 补充 k 前置声明
	// for k, v = range []int{} { // 这里保留 ASSIGN
	//	r,err := f(); if err ...
	// 	r.k = k
	//  ...
	// }
	// 只需要生成临时变量 i, 不需为 v 生成临时变量
	// 但是 range.Token 不能是 DEFINE, 因为 v 是外部变量
	// 所以, 需要前置补充一个 value decl, 但是类型字面量需要自己构造, 比如 range map 时候
	// ---------------------------
	// 所以, 简单的策略是只有 k,v 任意包含 try, 无脑生成 k v 两个迭代变量
	// var v int
	// for k_, v_ := range []int{} { // 无脑 DEFINE
	//	r,err := f(); if err ...
	// 	r.k = K_
	//	v = v_ // 坏处是, 这里会多出来冗余的赋值语句
	//  ...
	// }

	// <<对于 select 的 case 而言>>
	// =================================================================
	// 因为只有选中的分支才会对 lhs 求值, 所以, 一旦 lhs 中包含 try 调用
	// 则需要把 lhs 展开到 body (而不是 select 外), 与 range 的处理逻辑相同
	// e.g.
	// case xs[f()] = <-ch:
	//	same as:
	// =>
	// case t := <-c4
	//	xs[f()] = t

	var xs, ys []ast.Stmt
	tryWith := r.tryNodes[fst] || r.tryNodes[snd]
	if tryWith {
		// 只要有 try 调用, 就意味着 lhs 肯定不是两个 ident,
		// token 就不可能是 define,
		// e.g.
		// 	try().xxx, =
		// 	*try() =
		assert(oriTok == token.ASSIGN)
		if fst != nil {
			var fst1 ast.Expr
			fst1, xs = r.rewriteExpr(ctx.withTry(tryWith).withLhs(), fst)
			id := fstGen()
			xs = append(xs, &ast.AssignStmt{
				Lhs: []ast.Expr{fst1},
				Tok: oriTok, // 提取的 assignStmt 保留原 token
				Rhs: []ast.Expr{id},
			})
			*outFst = id
			*outTok = token.DEFINE
		}
		if snd != nil {
			var snd1 ast.Expr
			snd1, ys = r.rewriteExpr(ctx.withTry(r.tryNodes[snd]).withLhs(), snd)
			// 无论 snd 是否被改写, 都无脑 genSym
			// 这样即可统一将 tok 处理成 DEFINE
			id := sndGen()
			ys = append(ys, &ast.AssignStmt{
				Lhs: []ast.Expr{snd1},
				Tok: oriTok, // 提取的 assignStmt 保留原 token
				Rhs: []ast.Expr{id},
			})
			*outSndOpt = id
			*outTok = token.DEFINE
		}
	}
	return append(xs, ys...)
}
