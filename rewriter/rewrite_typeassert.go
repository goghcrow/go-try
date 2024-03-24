package rewriter

import "go/ast"

func (r *fileRewriter) rewriteTypeAssert(ctx *rewriteCtx, n *ast.TypeAssertExpr) (ast.Expr, []ast.Stmt) {
	var xs []ast.Stmt

	n1 := &ast.TypeAssertExpr{Type: n.Type}
	n1.X, xs = r.rewriteExpr(ctx, n.X)

	// 需要处理处理 x, y = v.(T) 形式, 避免
	// b, ok := Try(ret1Err[error]()).(error)
	// 被改写成
	// -------------------------------------
	// 	𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[error]()
	//	if 𝗲𝗿𝗿𝟭 != nil {
	//		return 𝗲𝗿𝗿𝟭
	//	}
	//	此处语义发生变化 !!!
	//	𝘃𝗮𝗹𝟮 := 𝘃𝗮𝗹𝟭.(error)
	//	v, ok := 𝘃𝗮𝗹𝟮
	// -------------------------------------
	// golang spec:
	// -------------------------------------
	// v, ok = x.(T)
	// v, ok := x.(T)
	// var v, ok = x.(T)
	// var v, ok interface{} = x.(T) // dynamic types of v and ok are T and bool
	// -------------------------------------
	// "A type assertion used in an assignment statement or initialization of the special form
	// yields an additional untyped boolean value.
	// The value of ok is true if the assertion holds.
	// Otherwise it is false and the value of v is the zero value for type T.
	// No run-time panic occurs in this case."
	// assignment 或者 init 形式 不能因为求值顺序原因, 把 v.(T) 提出出来
	// -------------------------------------
	// 另外, 为什么需要 unparen
	// e.g. v, ok := (𝘃𝗮𝗹𝟮)
	// 𝘃𝗮𝗹𝟮 的 parent 是 parenExpr, 而不是 assignStmt
	if isTuple2Assign(ctx, n) {
		return n1, xs
	}

	// n.Type == nil 代表 type switch X.(type), 也不能改写
	if n.Type == nil {
		return n1, xs
	}

	// 同理, TypeAssertExpr 也涉及到求值顺序
	// 需要按照求值顺序进行抽象解释
	// e.g. 非 _,_ = v.(T) 形式使用类型断言
	// e,g, v.(int) + try()
	// v.(int) 先于 try() 求值
	// v.(int) 可能会 panic
	// 需要先求值 v.(int) (提取 define) 再展开 try
	id, assign := r.assign(ctx, n1)
	return id, append(xs, assign)
}
