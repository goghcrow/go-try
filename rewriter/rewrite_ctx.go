package rewriter

import (
	"go/ast"
)

type rewriteCtx struct {
	// 用 try 来标记当前表达式是否包含 try 函数调用
	// 一个表达式的任意子节点 (函数作为边界, 同时排除函数自身) 包含 try 调用
	// 则为了保持求值顺序(语义), 整个表达式均需要改写,
	// try 的两个 scope
	// 1. fun: 每个 FuncDecl|FuncLit 的 try 状态是独立的, 遇到 fun 边界, 重置 try withFunc()
	// 2. stmt: 每个 Stmt 的 try 状态是独立的 (Init Stmt除外), 遇到 Stmt 重置 try withTry(false)
	//
	// e.g.
	// a() + try(b()), 处理 a() 的 ctx try=true
	// f(a(), try(b)) 类似
	// x, y = a(), try(b()) 类似
	//
	// **所有有求值顺序依赖的 stmt 或 expr 都会在一个 try 上下文中处理**
	// 具体看 withTry(true)
	try bool

	// 处于左值, selection, index 不能展开
	lhs bool

	// 当前节点所处的最近的函数定义
	// 不是必须, 通过 parent 向上遍历也可以找到
	fun fnNode // enclosing func

	node   ast.Node // 当前改写的节点
	parent *rewriteCtx
}

func newWalkCtx(n ast.Node) *rewriteCtx {
	return &rewriteCtx{node: n}
}

func (c *rewriteCtx) enter(n ast.Node) *rewriteCtx {
	return &rewriteCtx{try: c.try, lhs: c.lhs, fun: c.fun, node: n, parent: c}
}

func (c *rewriteCtx) withFun(n fnNode) *rewriteCtx {
	// 进入函数同时会重置 try, 子树是否在 try 上下文, 内部自己判断
	return &rewriteCtx{try: false, lhs: false, fun: n, node: c.node, parent: c.parent}
}

func (c *rewriteCtx) resetTry() *rewriteCtx {
	return c.withTry(false)
}

func (c *rewriteCtx) withTry(try bool) *rewriteCtx {
	return &rewriteCtx{try: try, lhs: c.lhs, fun: c.fun, node: c.node, parent: c.parent}
}

// 标记当前表达式为左值
// 1. assign.lhs
// 2. range.{k v}
// 3. select case recv 的 assign.lhs (这个 1 已经包含)
// 注意: 不需要标记 type switch 的 assign.lhs
// TypeSwitchStmt  = "switch" [ SimpleStmt ";" ] TypeSwitchGuard "{" { TypeCaseClause } "}" .
// TypeSwitchGuard = [ identifier ":=" ] PrimaryExpr "." "(" "type" ")" .
// 因为 只支持 id := ?.(T) 形式
func (c *rewriteCtx) withLhs() *rewriteCtx {
	return &rewriteCtx{try: c.try, lhs: true, fun: c.fun, node: c.node, parent: c.parent}
}

func (c *rewriteCtx) unparenParentNode() ast.Node {
	it := c.parent
	for it != nil {
		_, ok := it.node.(*ast.ParenExpr)
		if !ok {
			return it.node
		}
		it = it.parent
	}
	return nil
}
