package rewriter

import (
	"go/ast"
)

func (r *fileRewriter) rewriteRange(ctx *rewriteCtx, n *ast.RangeStmt) (ast.Stmt, []ast.Stmt) {
	n1 := &ast.RangeStmt{
		Key:   n.Key,
		Value: n.Value,
		Tok:   n.Tok,
		Body:  &ast.BlockStmt{},
	}

	//  range 的 k v 是支持函数调用的
	//  e.g.
	//	getvar := func(p *int) *int { return p }
	//	var i, v int
	//	for *getvar(&i), *getvar(&v) = range [2]int{1, 2} {
	//		println(i, v)
	//	}
	// output: 0 1 \n  1 2

	// golang spec
	// "As with an assignment, if present the operands on the left
	// must be addressable or map index expressions; they denote the iteration variables."
	// 所以
	// 		type X struct {
	//			i int
	//			v int
	//		}
	//		x := X{}
	//		for x.i, x.v = range []int{} { }
	//		f := func() *X { return &X{} }
	//		for f().i, f().v = range []int{} { }
	//		for slice[0], slice[1] = range []int{} { }
	//		for mapX[k], mapX[v] = range []int{} { }
	// 都是合法的

	// The range expression x is evaluated once before beginning the loop,
	// with one exception: if at most one iteration variable is present and len(x)
	// is constant, the range expression is not evaluated.
	// Function calls on the left are evaluated once per iteration.

	// e.g.
	// type X struct{ i int }
	// lhs := func(s string) *X {
	// 	print(s)
	// 	return &X{}
	// }
	// for lhs("K").i, lhs("V").i = range []int{1, 2, 3} {
	// 	print("X")
	// }
	// ---------------------------
	// output: KVXKVXKVX

	var xs []ast.Stmt
	n1.X, xs = r.rewriteExpr(ctx.withTry(r.tryNodes[n.X]), n.X)

	// 参见 rewriteTuple2Assign 注释

	kGen := func() *ast.Ident { return r.genSym(ctx.fun, "k") }
	vGen := func() *ast.Ident { return r.genSym(ctx.fun, "v") }
	insertBefore := r.rewriteTuple2Assign(ctx, &n1.Key, &n1.Value, &n1.Tok, kGen, vGen)

	// 都是生成的 sym, 没有 shadow 的问题, 不用套 block
	// n1.Body.List = append(insertBefore, &ast.BlockStmt{ List: r.rewriteStmtList(ctx, n.Body.List) })
	n1.Body.List = append(insertBefore, r.rewriteStmtList(ctx, n.Body.List)...)
	return n1, xs
}
