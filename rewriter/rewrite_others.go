package rewriter

import (
	"fmt"
	"go/ast"
)

// 整体可以看成, 按操作语义, 求值的顺序进行抽象解释的过程

// ===================================================================================================

// rewriteExpr / rewriteExprList 调用时候, 需要根据求值顺序, 手动控制 ctx.withTry()
// rewriteStmt / rewriteStmtList 不涉及求值顺序问题, 每个单独的 stmt 自动 resetTry()

// golang spec
// https://go.dev/ref/spec#Order_of_evaluation

// At package level, initialization dependencies determine the evaluation order of individual
// initialization expressions in variable declarations.
// Otherwise, when evaluating the operands of an expression, assignment, or return statement,
// all function calls, method calls, receive operations, and binary logical operations are evaluated
// in lexical left-to-right order.
//
// 1. 包级别, 初始化依赖决定初始化表达式的求值顺序
// 2. 表达式的操作数 / 赋值语句 / 返回语句 / 函数调用 / 方法调用 / recv 表达式 / 二元逻辑表达式
// 都是按照词法顺序(源码顺序)从左到右进行求值
// 第二条是改写 try 关注的, 因为 try 不能在 pkg level 的 init stmt 中出现
//
// y[f()], ok = g(z || h(), i()+x[j()], <-c), k()  (函数体内赋值, 非包级别赋值语句)
// the function calls and communication happen in the order f(), h() (if z evaluates to false), i(), j(), <-c, g(), and k().
// However, the order of those events compared to the evaluation and indexing of x and the evaluation of y and z is not specified,
// except as required lexically. For instance, g cannot be called before its arguments are evaluated.
// x[..], z, y 顺序与函数调用, <- 之间只有词法顺序的依赖, 其他未定义, 比如 x[j()求值的结果] 与 <- 之间
// 未定义的行为统一按照书写顺序进行

// ===================================================================================================

// 会产生 runtime panics 的 stmt 或者 expr 都代表副作用
// 都需要严格遵循求值的顺序
// 总结下 golang spec, 可能发生 runtime panic 节点
// -------------------------------------
// 1. [Map Types]The comparison operators == and != must be fully defined for
// 	operands of the key type; thus the key type must not be a function,
// 	map, or slice. If the key type is an interface type, these comparison
// 	operators must be defined for the dynamic key values; failure will cause a run-time panic.
// 		m := map[any]int{}
//		m[func() {}] = 1
//		m[map[int]int{}] = 1
//		m[[]int{}] = 1
// 2. [Selectors] If x is of pointer type and has the value nil and x.f denotes a struct field,
// 	assigning to or evaluating x.f causes a run-time panic.
//	 	var p *Point
//		p.X = 1 // panic
//		consume(p.X) // panic
//		赋值 (作为左值) / 求值 (比如作为右值, 作为参数)
// 3. [Selectors]If x is of interface type and has the value nil, calling or evaluating the method x.f causes a run-time panic.
//		var a ast.Node
//		a.Pos // panic
//		a.Pos() // panic
// 4. [Index a[x]] For a of array type A: if x is out of range at run time, a run-time panic occurs
// 5. [Index a[x]] For a of slice type S: if x is out of range at run time, a run-time panic occurs
// 6. [Index a[x]] For a of string type S: if x is out of range at run time, a run-time panic occurs
// 7. [Index expressions] Assigning to an element of a nil map causes a run-time panic.
// 	注意这里是赋值, 从 nil map 读不会panic, 只会返回 zero value
// 8. [Simple slice expressions] a[low : high], If the indices are out of range at run time, a run-time panic occurs.
// 9. [Full slice expressions] a[low : high : max], If the indices are out of range at run time, a run-time panic occurs.
// 10. [Type assertions] x.(T),  If the type assertion is false, a run-time panic occurs.
// 11. [Calls] Calling a nil function value causes a run-time panic.
// 12. [Integer operators] If the divisor is zero at run time, a run-time panic occurs.
// 13. [Integer operators] If the shift count is negative at run time, a run-time panic occurs.
// 14. [Comparison operators] A comparison of two interface values with identical dynamic types causes a run-time panic
// 							 if that type is not comparable.
//		 					This behavior applies not only to direct interface value comparisons but also when comparing
//		 					arrays of interface values or structs with interface-valued fields.
// 15. [Address operators] If the evaluation of x would cause a run-time panic, then the evaluation of &x does too.
// 16. [Address operators] If x is nil, an attempt to evaluate *x will cause a run-time panic.
// 17. [Conversions from slice to array or array pointer] if the length of the slice is less than the length of the array,
// 							a run-time panic occurs.
// 18. [Send statements] A send on a closed channel proceeds by causing a run-time panic. A send on a nil channel blocks forever.
// 19. [Close] Sending to or closing a closed channel causes a run-time panic.
// 20. [Close] Closing the nil channel also causes a run-time panic.
// 21. [Making slices, maps and channels] make(T, n, m), make(T, n), For slices and channels,
// 							if n is negative or larger than m at run time, a run-time panic occurs.
// 22. [Panics] an explicit call to panic
// 23. [Panics] calling panic with a nil interface value (or an untyped nil) causes a run-time panic.
// 24. [unsafe] (*[len]ArbitraryType)(unsafe.Pointer(ptr))[:], At run time, if len is negative,
// 						or if ptr is nil and len is not zero, a run-time panic occurs [Go 1.17].
// -------------------------------------
// 1. [ast.IndexExpr] Index.X 必须是 map
// 2. [ast.SelectorExpr] ptrtype, Selector 必须代表 struct 访问
// 3. [ast.SelectorExpr] interface type | [ast.CallExpr] 并且 call.Fun 必须是 ast.SelectorExpr 且代表 struct 访问
// 4. [ast.IndexExpr] Index.X 必须是 array
// 5. [ast.IndexExpr] Index.X 必须是 slice
// 6. [ast.IndexExpr] Index.X 必须是 string
// 7. [ast.IndexExpr] Index.X 必须是 map
// 8. [ast.SliceExpr] Slice.X 必须是 array/slice
// 9. [ast.SliceExpr] Slice.X 必须是 array/slice
// 10. [ast.TypeAssertExpr]
// 11. [ast.CallExpr] 函数调用而非 TypeConv
// 12. [ast.BinaryExpr] /
// 13. [ast.BinaryExpr] <<
// 14. [ast.BinaryExpr] == !=
// 15. [ast.UnaryExpr] &
// 16. [ast.StarExpr] *
// 17. [ast.CallExpr] TypeConv
// 18. [ast.SendStmt]
// 19. [ast.CallExpr] close
// 20. [ast.CallExpr] close
// 21. [ast.CallExpr] make
// 22/23. [ast.CallExpr] panic
// 24. 不处理了
// -------------------------------------
// 节点
// [x] ast.IndexExpr		 Index.X map/array/slice/string
// [x] ast.SelectorExpr		 ptrtype / interface type
// [x] ast.SliceExpr		 array / slice
// [x] ast.TypeAssertExpr
// [x] ast.CallExpr			call / typeconv
// [x] ast.BinaryExpr		/ <<
// [x] ast.UnaryExpr		&x, x会 panic, x 已经处理了, &x 不需要处理
// [x] ast.StarExpr
// [x] ast.SendStmt		本身是 Stmt, 不是 Expr, 不需要处理求值顺序问题

// ===================================================================================================

func (r *fileRewriter) rewriteExpr(ctx *rewriteCtx, n ast.Expr) (ast.Expr, []ast.Stmt) {
	ctx = ctx.enter(n)
	if !ctx.try {
		return n, nil
	}

	var xs, ys, zs []ast.Stmt

	switch n := n.(type) {
	case *ast.BadExpr, *ast.Ident, *ast.Ellipsis, *ast.BasicLit:
		return n, nil

	case *ast.FuncLit:
		fn, _ := r.rewriteFun(ctx, n)
		return fn.(ast.Expr), nil

	case *ast.CompositeLit:
		n1 := &ast.CompositeLit{
			Type:       n.Type,
			Incomplete: n.Incomplete,
		}
		n1.Elts, xs = r.rewriteExprList(ctx, n.Elts)
		return n1, xs

	case *ast.ParenExpr:
		n1 := &ast.ParenExpr{}
		n1.X, xs = r.rewriteExpr(ctx, n.X)
		return n1, xs

	case *ast.SelectorExpr:
		return r.rewriteSelector(ctx, n)

	case *ast.IndexExpr:
		return r.rewriteIndex(ctx, n)

	case *ast.IndexListExpr:
		// X[type parameter...] 不需要处理求值顺序
		n1 := &ast.IndexListExpr{}
		n1.X, xs = r.rewriteExpr(ctx, n.X)
		n1.Indices, ys = r.rewriteExprList(ctx, n.Indices)
		assert(allTypeNames(r.pkg.TypesInfo, n.Indices...))
		assert(len(ys) == 0) // 类型参数, 不应该改写出来 stmt
		return n1, xs

	case *ast.SliceExpr:
		// 同理, SliceExpr 也涉及到求值顺序
		// 需要按照求值顺序进行抽象解释
		// e.g. x[1:2:3] + try()
		// x[1:2:3] 先于 try() 求值
		// x[1:2:3] 可能会 panic
		// 需要先求值 x[1:2:3] (提取 define) 再展开 try
		var ws []ast.Stmt
		n1 := &ast.SliceExpr{Slice3: n.Slice3}
		n1.X, ws = r.rewriteExpr(ctx, n.X)
		if n.Low != nil {
			n1.Low, xs = r.rewriteExpr(ctx, n.Low)
		}
		if n.High != nil {
			n1.High, ys = r.rewriteExpr(ctx, n.High)
		}
		if n.Max != nil {
			n1.Max, zs = r.rewriteExpr(ctx, n.Max)
		}
		id, assign := r.assign(ctx, n1)
		return id, concat(ws, xs, ys, zs, sliceOf(assign))

	case *ast.TypeAssertExpr:
		return r.rewriteTypeAssert(ctx, n)

	case *ast.CallExpr:
		return r.rewriteCall(ctx, n)

	case *ast.StarExpr:
		if isTypeName(r.pkg.TypesInfo, n) {
			return n, nil
		}

		// 同理, StarExpr 也涉及到求值顺序
		// 需要按照求值顺序进行抽象解释
		// e.g. *ptr + try()
		// *ptr 先于 try() 求值
		// *ptr 可能会 panic
		// 需要先求值 *ptr (提取 define) 再展开 try
		n1 := &ast.StarExpr{}
		n1.X, xs = r.rewriteExpr(ctx, n.X)
		// star 左值不能提取表达式
		// e.g. *a = try(...)
		if ctx.lhs {
			return n1, xs
		}
		id, assign := r.assign(ctx, n1)
		return id, append(xs, assign)

	case *ast.UnaryExpr:
		n1 := &ast.UnaryExpr{Op: n.Op}
		n1.X, xs = r.rewriteExpr(ctx, n.X)

		// <- 不会发生 runtime panic, 所以不会发生求值顺序语义的问题
		// 所以不需要因为求值顺序原因展开
		// 所以不需要特殊处理 `v, ok := <-ch`
		// if false {
		// 	isChRecv := n.Op == token.ARROW
		// 	if isChRecv && isTuple2Assign(ctx, n) {
		// 		return n1, xs
		// 	}
		// }

		return n1, xs

	case *ast.BinaryExpr:
		return r.rewriteBinaryExpr(ctx, n)

	case *ast.KeyValueExpr:
		n1 := &ast.KeyValueExpr{}
		n1.Key, xs = r.rewriteExpr(ctx, n.Key)
		n1.Value, ys = r.rewriteExpr(ctx, n.Value)
		return n1, append(xs, ys...)

	// Types
	case *ast.ArrayType, *ast.StructType, *ast.FuncType, *ast.InterfaceType, *ast.MapType, *ast.ChanType:
		return n, nil

	default:
		panic(fmt.Sprintf("unexpected node type %T", n))
	}
}

func (r *fileRewriter) rewriteExprList(ctx *rewriteCtx, list []ast.Expr) (ys []ast.Expr, ext []ast.Stmt) {
	// 原来是 nil, 仍需要返回 nil, e.g.
	// 	ast.CaseClause.List []Expr, nil means default case
	if list == nil {
		return
	}

	// list 元素不能根据自身是否是 tryNode 决定是否改写
	// 比如 x = f(a(), b()) + Try(c()),
	// c() 不能提出出来判断 err, 因为求值顺序是 a() -> b() -> f(...) -> c() -> +
	// 所以, rhs binary expr 包含 try, 整个 expr (stmt) 都需要展开
	// 准确说,
	// **所有有求值顺序依赖的 stmt / expr 都会在一个 try 上下文中处理**
	// 以下表达式列表:
	// - 实参列表
	// - CompositeLit.Elts
	// - AssignStmt.Rhs
	// - ValueSpec.Values
	// - ReturnStmt.Results
	// 都有求值顺序问题, 如果其中包含 try 调用, 需要全部展开

	ys = make([]ast.Expr, len(list))
	for i, it := range list {
		x, xs := r.rewriteExpr(ctx, it)
		ys[i] = x
		ext = append(ext, xs...)
	}
	return ys, ext
}

func (r *fileRewriter) rewriteStmt(ctx *rewriteCtx, n ast.Stmt) (ast.Stmt, []ast.Stmt) {
	// 每一个 stmt 的 ctx.try (是否包含 try 调用) 都需要单独判断
	// stmt 之间的 ctx.try 是独立, 需要 reset ctx
	ctx = ctx.enter(n).resetTry()

	var xs, ys, zs []ast.Stmt

	switch n := n.(type) {
	case *ast.BadStmt, *ast.EmptyStmt:
		return n, nil

	case *ast.DeclStmt:
		n1 := &ast.DeclStmt{}
		n1.Decl, xs = r.rewriteDecl(ctx, n.Decl)
		return n1, xs

	// 在 rewriteStmtList 中已经处理
	// 理论上不应该走到这个分支
	// case *ast.LabeledStmt:
	// 	ys, xs = r.rewriteLabeled(ctx, n)
	// 	assert(len(ys) == 1)
	// 	// 只有 StmtList 允许 labeledStmt
	// 	// labeled for/range/switch/select 也返回 len(ys) == 1
	// 	return ys[0].(*ast.LabeledStmt), xs

	case *ast.ExprStmt:
		n1 := &ast.ExprStmt{}
		// exprStmt 是否需要改写由 root 节点是否包含 try 调用决定
		ctx = ctx.withTry(r.tryNodes[n])
		n1.X, xs = r.rewriteExpr(ctx, n.X)
		return n1, xs

	case *ast.SendStmt:
		n1 := &ast.SendStmt{Arrow: n.Arrow}
		// n.Chan 与 n.Value 有求值顺序依赖
		ctx = ctx.withTry(r.tryNodes[n])
		n1.Chan, xs = r.rewriteExpr(ctx, n.Chan)
		n1.Value, ys = r.rewriteExpr(ctx, n.Value)
		return n1, append(xs, ys...)

	case *ast.IncDecStmt:
		n1 := &ast.IncDecStmt{Tok: n.Tok}
		ctx = ctx.withTry(r.tryNodes[n])
		n1.X, xs = r.rewriteExpr(ctx, n.X)
		return n1, xs

	case *ast.AssignStmt:
		return r.rewriteAssign(ctx, n)

	case *ast.GoStmt:
		n1 := &ast.GoStmt{}
		ctx = ctx.
			withTry(r.tryNodes[n])
		e, ex1 := r.rewriteExpr(ctx, n.Call)
		call, ok := e.(*ast.CallExpr)
		r.assert(n.Call, len(ex1) == 0 && ok, "try is not allowed in go stmt")
		n1.Call = call
		return n1, ex1

	case *ast.DeferStmt:
		n1 := &ast.DeferStmt{}
		ctx = ctx.withTry(r.tryNodes[n])
		e, ex1 := r.rewriteExpr(ctx, n.Call)
		call, ok := e.(*ast.CallExpr)
		r.assert(n.Call, len(ex1) == 0 && ok, "defer is not allowed in go stmt")
		n1.Call = call
		return n1, ex1

	case *ast.ReturnStmt:
		n1 := &ast.ReturnStmt{}
		// results 元素存在求值顺序依赖
		ctx = ctx.withTry(r.tryNodes[n])
		n1.Results, xs = r.rewriteExprList(ctx, n.Results)
		return n1, xs

	case *ast.BranchStmt:
		return n, nil

	case *ast.BlockStmt:
		n1 := &ast.BlockStmt{}
		n1.List = r.rewriteStmtList(ctx, n.List)
		return n1, nil

	case *ast.IfStmt:
		fixScope := false
		n1 := &ast.IfStmt{Body: &ast.BlockStmt{}}
		if n.Init != nil {
			initCtx := ctx.withTry(r.tryNodes[n.Init]).enter(n.Init)
			n1.Init, xs = r.rewriteInitStmt(initCtx, n.Init)
			if r.tryNodes[n.Cond] {
				// init 和 cond 之间有求值顺序依赖
				// cond 中包含 try, cond 则需要展开
				// 而 cond 可以依赖 init, 所以
				// init 必须先展开 (提取到 if 外)
				// 又 init 可能是 assign(define) stmt,
				// 且 init 是独立的作用域, 所以 提取 init
				// (独立作用域 if a := 1; true {  a := 2 } )
				// 必须显式样再用 blockStmt 恢复作用域
				// 反例:
				// 	a:=1
				// 	if a:=1; Try(a); {}
				// 	a:=1 从 init 提到外头, 符号会冲突
				// 所有 init 处理逻辑一致
				xs = append(xs, n1.Init)
				n1.Init = nil
				fixScope = true
			}
		}
		n1.Cond, ys = r.rewriteExpr(ctx.withTry(r.tryNodes[n.Cond]), n.Cond)
		n1.Body.List = r.rewriteStmtList(ctx, n.Body.List)
		if n.Else != nil {
			n1.Else, zs = r.rewriteStmt(ctx, n.Else)
		}
		if fixScope {
			return &ast.BlockStmt{
				List: concat(xs, ys, zs, sliceOf[ast.Stmt](n1)),
			}, nil
		}
		return n1, concat(xs, ys, zs)

	case *ast.CaseClause:
		n1 := &ast.CaseClause{}
		// 理论上 n.List 不需要改写
		// 因为 case.List 如果包含 call, 一定会被预处理成 if 形式
		n1.List, xs = r.rewriteExprList(ctx, n.List)
		// 只有存在 try 被改写, xs 才会有值
		assert(len(xs) == 0)
		n1.Body = r.rewriteStmtList(ctx, n.Body)
		return n1, xs

	case *ast.SwitchStmt:
		fixScope := false
		n1 := &ast.SwitchStmt{Body: &ast.BlockStmt{}}
		if n.Init != nil {
			initCtx := ctx.withTry(r.tryNodes[n.Init]).enter(n.Init)
			n1.Init, xs = r.rewriteInitStmt(initCtx, n.Init)
			// init 和 tag 之间有求值顺序依赖, 参见改写 ifStmt 逻辑
			// e.g. switch i := 42; Try(func1[int, B](i)) { }
			if r.tryNodes[n.Tag] {
				xs = append(xs, n1.Init)
				n1.Init = nil
				fixScope = true
			}
		}
		if n.Tag != nil {
			n1.Tag, ys = r.rewriteExpr(ctx.withTry(r.tryNodes[n.Tag]), n.Tag)
		}
		// CaseClauses only
		n1.Body.List = r.rewriteStmtList(ctx, n.Body.List)
		if fixScope {
			return &ast.BlockStmt{
				List: concat(xs, ys, sliceOf[ast.Stmt](n1)),
			}, nil
		}
		return n1, append(xs, ys...)

	case *ast.TypeSwitchStmt:
		fixScope := false
		n1 := &ast.TypeSwitchStmt{Body: &ast.BlockStmt{}}
		if n.Init != nil {
			initCtx := ctx.withTry(r.tryNodes[n.Init]).enter(n.Init)
			n1.Init, xs = r.rewriteInitStmt(initCtx, n.Init)
			// init 和 assign 之间有求值顺序依赖, 参见改写 ifStmt 逻辑
			if r.tryNodes[n.Assign] {
				xs = append(xs, n1.Init)
				n1.Init = nil
				fixScope = true
			}
		}
		n1.Assign, ys = r.rewriteStmt(ctx, n.Assign)
		n1.Body.List = r.rewriteStmtList(ctx, n.Body.List)
		if fixScope {
			return &ast.BlockStmt{
				List: concat(xs, ys, sliceOf[ast.Stmt](n1)),
			}, nil
		}
		return n1, append(xs, ys...)

	case *ast.CommClause:
		return r.rewriteCommClause(ctx, n)

	case *ast.SelectStmt:
		return r.rewriteSelect(ctx, n)

	case *ast.ForStmt:
		n1 := &ast.ForStmt{Body: &ast.BlockStmt{}}
		// cond 和 post 已经前置改写到 forStmt 的 body 了,
		// 所以, init 和 cond 和 post 已经不是兄弟关系,
		// 天然通过父子节点构建了求值顺序, 不需要额外处理 ctx,
		// 不会将 init 提取到 for 外层, 则不会导致 init 的 define 和外层作用域符号冲突
		// 所以, 也不需要额外处理作用域
		// 具体参见改写 ifStmt 逻辑
		if n.Init != nil {
			n1.Init, xs = r.rewriteStmt(
				ctx.withTry(r.tryNodes[n.Init]),
				n.Init,
			)
		}
		assert(!r.tryNodes[n.Cond])
		// n.Cond, ys = rewriteExpr(ctx, n.Cond)
		n1.Cond = n.Cond
		assert(!r.tryNodes[n.Post])
		// if n.Post != nil {
		// 	n.Post, zs = rewriteStmt(ctx, n.Post)
		// }
		n1.Post = n.Post
		n1.Body.List = r.rewriteStmtList(ctx, n.Body.List)
		// return sliceOf[ast.Stmt](n1), concat(xs, ys, zs)
		return n1, xs

	case *ast.RangeStmt:
		return r.rewriteRange(ctx, n)

	default:
		panic(fmt.Sprintf("unexpected node type %T", n))
	}
}

func (r *fileRewriter) rewriteStmtList(ctx *rewriteCtx, list []ast.Stmt) (zs []ast.Stmt) {
	for _, it := range list {
		// 参见 rewriteLabeled(...) 注释
		if l, ok := it.(*ast.LabeledStmt); ok {
			xs, ys := r.rewriteLabeled(ctx, l)
			zs = concat(zs, ys, xs)
		} else {
			x, ys := r.rewriteStmt(ctx, it)
			zs = concat(zs, ys, sliceOf(x))
		}
	}
	return
}

// IfStmt | SwitchStmt | TypeSwitchStmt | ForStmt 均包含 Init Stmt
// IfStmt 与 ForStmt 中 Init 与 Cond 存在求值顺序依赖 (ForStmt 已经前置处理)
// SwitchStmt 中 Init 与 Tag 存在求值顺序依赖
// TypeSwitchStmt 中 Init 与 Assign 存在求值顺序依赖
// 所以, Init 与依赖 Init 不分要一起改写, 不能只展开依赖 Init 不分, 会破坏语义
// Init 是 SimpleStmt,
// SimpleStmt = EmptyStmt | ExpressionStmt | SendStmt | IncDecStmt | Assignment | ShortVarDecl .
// 所以 rewriteInitStmt 对 SimpleStmt 定制处理, 与被依赖部分 共享 ctx try? 状态, 而不是重置
func (r *fileRewriter) rewriteInitStmt(ctx *rewriteCtx, n ast.Stmt) (ast.Stmt, []ast.Stmt) {
	assert(ctx.node == n)

	var xs, ys []ast.Stmt
	switch n := n.(type) {
	case *ast.BadStmt, *ast.EmptyStmt:
		return n, nil

	case *ast.ExprStmt:
		n1 := &ast.ExprStmt{}
		n1.X, xs = r.rewriteExpr(ctx, n.X)
		return n1, xs

	case *ast.SendStmt:
		n1 := &ast.SendStmt{Arrow: n.Arrow}
		n1.Chan, xs = r.rewriteExpr(ctx, n.Chan)
		n1.Value, ys = r.rewriteExpr(ctx, n.Value)
		return n1, append(xs, ys...)

	case *ast.IncDecStmt:
		n1 := &ast.IncDecStmt{Tok: n.Tok}
		n1.X, xs = r.rewriteExpr(ctx, n.X)
		return n1, xs

	case *ast.AssignStmt:
		// 参照非 init ast.AssignStmt 的改写过程
		n1 := &ast.AssignStmt{
			Lhs: n.Lhs,
			Tok: n.Tok,
		}

		// golang spec
		// https://go.dev/ref/spec#Assignment_statements
		// Each left-hand side operand must be addressable, a map index expression, or (for = assignments only) the blank identifier.
		// Operands may be parenthesized.
		//
		// lhs = addr | array | slice | map | blank
		//
		// The assignment proceeds in two phases.
		// First, the operands of index expressions and pointer indirections
		// (including implicit pointer indirections in selectors) on the left
		// and the expressions on the right are all evaluated in the usual order.
		// Second, the assignments are carried out in left-to-right order.
		//
		// 复制语句的求值顺序:
		// e.g.
		// func identity[T any](t T, msg string) T {
		//	println(msg)
		//	return t
		// }
		//
		// 	var i int
		//	*identity(&i, "lhs") = identity(42, "rhs")
		//	println(i)
		//
		// output:
		// 	lhs \n rhs \n 42
		//
		// eval(lhs) -> eval(rhs) -> assign
		//
		// RangeStmt 中的 K,V 同样遵循 Assign 的求值顺序

		n1.Lhs, xs = r.rewriteExprList(ctx.withLhs(), n.Lhs)
		n1.Rhs, ys = r.rewriteExprList(ctx, n.Rhs)
		return n1, append(xs, ys...)
	default:
		panic(fmt.Sprintf("unexpected node type %T", n))
	}
}

func (r *fileRewriter) rewriteSpec(ctx *rewriteCtx, n ast.Spec) (ast.Spec, []ast.Stmt) {
	ctx = ctx.enter(n)

	switch n := n.(type) {
	case *ast.ImportSpec:
		return n, nil

	case *ast.TypeSpec:
		return n, nil

	case *ast.ValueSpec:
		n1 := &ast.ValueSpec{
			Names: n.Names,
			Type:  n.Type,
		}
		// e.g.
		// var ( // GenDecl
		// 	a,b = f1(), Try(f2()) // ValueSpec
		//	c = f3() // ValueSpec
		// )
		// ValueSpec 相当于 Stmt, 有独立的 ctx.try
		// e.g.
		// var a,b = f1(), Try(f2())
		// f1(), f2() 求值有先后关系, 不能直接把 f2() 提到前面
		withTry := findFirst(n.Values, func(s ast.Expr) bool { return r.tryNodes[s] }) != nil
		ys, ext := r.rewriteExprList(ctx.withTry(withTry), n.Values)
		n1.Values = ys
		return n1, ext

	default:
		panic(fmt.Sprintf("unexpected node type %T", n))
	}
}

func (r *fileRewriter) rewriteDecl(ctx *rewriteCtx, n ast.Decl) (x ast.Decl, xs []ast.Stmt) {
	ctx = ctx.enter(n)

	switch n := n.(type) {
	case *ast.BadDecl:
		return n, nil

	case *ast.FuncDecl:
		fnDecl, xs := r.rewriteFun(ctx, n)
		return fnDecl.(ast.Decl), xs

	case *ast.GenDecl:
		n1 := &ast.GenDecl{
			Tok:   n.Tok,
			Specs: make([]ast.Spec, len(n.Specs)),
		}
		var ext []ast.Stmt
		for i, s := range n.Specs {
			n1.Specs[i], ext = r.rewriteSpec(ctx, s)
			xs = append(xs, ext...)
		}
		return n1, xs

	default:
		panic(fmt.Sprintf("unexpected node type %T", n))
	}
}

func (r *fileRewriter) rewriteDeclList(ctx *rewriteCtx, list []ast.Decl) (decls []ast.Decl, xs []ast.Stmt) {
	decls = make([]ast.Decl, len(list))
	var ext []ast.Stmt
	for i, it := range list {
		decls[i], ext = r.rewriteDecl(ctx, it)
		xs = append(xs, ext...)
	}
	return
}
