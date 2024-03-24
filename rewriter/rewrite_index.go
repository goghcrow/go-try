package rewriter

import (
	"go/ast"
)

// 1. array index
// 2. slice index
// 3. string index
// 4. map index
// 5. type parameter index
func (r *fileRewriter) rewriteIndex(ctx *rewriteCtx, n *ast.IndexExpr) (ast.Expr, []ast.Stmt) {
	var xs, ys []ast.Stmt
	n1 := &ast.IndexExpr{}
	n1.X, xs = r.rewriteExpr(ctx, n.X)
	n1.Index, ys = r.rewriteExpr(ctx, n.Index)
	// isTypeName 分支可以直接原样返回, 这样保持代码统一
	// StarExpr 中已经判断当前是类型还是 deref
	// n1.Index = n.Index

	if ctx.lhs {
		// 如果 selector 是左值, 不能展开
		// e.g. xs[0] = 1 !=> v = xs[0]; v = 1
		return n1, append(xs, ys...)
	}

	if isTypeName(r.pkg.TypesInfo, n.Index) {
		// X[type parameter]
		// e.g. Try(ret1Err[*Struct]())
		assert(len(ys) == 0)
		return n1, xs
	}

	// 👇🏻 map index 一定是读取, 写(assign) 在 ctx.lhs 分支提前 return 了

	if mt, ok := isMapIndex(r.pkg.TypeInfo(), n); ok {
		// <rewrite map index>
		// 1. map 的 key 是 interface{} 可能会因为 key 不可比较而发生运行时异常
		// 2. map 是 nil, 会因为对 nil map 赋值而产生运行时异常
		// [golang spec] A nil map behaves like an empty map when reading,
		// but attempts to write to a nil map will cause a runtime panic; don't do that.
		// 3. nil map 读不会 panic, 写会, 但写本身是 stmt, 没有求值顺序问题, 且在 ctx.lhs 分支提前返回了
		if isAny(mt.Key(), true) {
			if isTuple2Assign(ctx, n) {
				// e.g. i, ok := map[any]int{}[?] 不能展开, 否则语法错误
				return n1, append(xs, ys...)
			} else {
				// 可能 panic, 则改写需要遵守求值顺序, 不能将 try 展开到 map index 之前
				// e.g. map[any]int{}[?] + try()
			}
		} else {
			// 非 map[~any]T 的 index 读取一定不会 panic, 不需要提取
			return n1, append(xs, ys...)
		}
	}

	// others, array/slice/string index
	// IndexExpr 涉及到求值顺序
	// 需要按照求值顺序进行抽象解释
	// e.g. x[1] + try()
	// x[1] 先于 try() 求值
	// x[1] 可能会 panic: out of index
	// 需要先求值 x[1] (并提取 assign) 再展开 try
	id, assign := r.assign(ctx, n1)
	return id, concat(xs, ys, sliceOf(assign))
}
