package rewriter

import (
	"go/ast"
)

func (r *fileRewriter) rewriteSelector(ctx *rewriteCtx, n *ast.SelectorExpr) (ast.Expr, []ast.Stmt) {
	var xs []ast.Stmt

	// pkg.Ident 不需要处理
	if isPkgSel(r.pkg.TypesInfo, n) {
		return n, nil
	}
	// 这里应该已经被 pkgSel 覆盖了
	//if isTypeName(r.pkg.TypesInfo, n) {
	//	return n, nil
	//}

	// SelectorExpr 涉及到求值顺序
	// 需要按照求值顺序进行抽象解释
	// e.g. v.x + try()
	// v.x 先于 try() 求值
	// v.x 可能会 panic
	// 需要先求值 v.x (提取 define) 再展开 try

	n1 := &ast.SelectorExpr{Sel: n.Sel}
	n1.X, xs = r.rewriteExpr(ctx, n.X)
	// n.Sel 是 ident, 不需要处理

	// 如果 selector 是左值, 不能展开
	// e.g. a.b = 1 !=> v = a.b; v = 1
	if ctx.lhs {
		return n1, xs
	}

	// golang spec
	// If x is of pointer type and has the value nil and x.f denotes a struct field,
	// 		assigning to or evaluating x.f causes a run-time panic.
	//		注意这里说 x.f 代表字段, 而不是方法!!!
	//	 	var p *Point
	//		p.X = 1 // panic
	//		consume(p.X) // panic
	//		p.Method // 不会 panic
	// If x is of interface type and has the value nil, calling or evaluating the method x.f causes a run-time panic.
	//		var a ast.Node
	//		a.Pos // panic
	//		a.Pos() // panic

	// e.g.
	// consume := func(...any) {}
	//	var x ast.Node
	//	consume(x.Pos, func() int { println("snd"); return 42 }())
	//	consume(x.Pos(), func() int { println("snd"); return 42 }())
	// 都会直接panic 不会输出 snd

	// e.g.
	// consume := func(...any) {}
	//	var x *ast.CallExpr
	//	consume(x.Pos(), func() int { println("snd"); return 42 }())
	// 取决于 CallExpr.Pos 实现, 实际是 CallExpr.Pos(x)
	//	consume(x.Pos, func() int { println("snd"); return 42 }())
	// 不会 panic

	xt := r.pkg.TypesInfo.TypeOf(n.X)
	st, isStPtr := typeOfStructPtr(xt, true)
	switch {
	// 只有 isStPtr | iface 才会触发 runtime panic
	case isStPtr || isIface(xt, true):
		// 特例: iface_or_ptr.m() 不需要处理成
		// x := iface_or_ptr.m
		// y := x()
		if isCallFun(ctx, n) {
			return n1, xs
		}

		// 特例: nil struct ptr func eval (not call)
		// var x *ast.CallExpr
		// consume2(x.Pos, try(...))
		// x.Pos 索引到函数, 而不是字段, 不会 panic
		// 注意: 不是 consume2(x.Pos(), try(...))
		// x.Pos() 可能 panic, 取决于 Pos 内部实现
		if isStPtr /*&& not call fun*/ {
			// isFunc
			sel := r.pkg.ObjectOf(n.Sel)
			if !isStructField(sel, st) {
				return n1, xs
			}
		}

		id, assign := r.assign(ctx, n1)
		return id, append(xs, assign)
	default:
		return n1, xs
	}
}
