package rewriter

import (
	"fmt"
	"github.com/goghcrow/go-loader"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ast/astutil"
)

func (r *fileRewriter) rewriteCall(ctx *rewriteCtx, n *ast.CallExpr) (ast.Expr, []ast.Stmt) {
	tryFn := r.tryCallee(n)
	if tryFn == "" {
		return r.rewriteNonTryCall(ctx, n)
	} else {
		return r.rewriteTryCall(ctx, n, tryFn)
	}
}

func (r *fileRewriter) rewriteNonTryCall(ctx *rewriteCtx, n *ast.CallExpr) (ast.Expr, []ast.Stmt) {
	var xs, ys []ast.Stmt

	// 先"求值" fun, 再"求值" args
	n1 := &ast.CallExpr{}
	n1.Fun, xs = r.rewriteExpr(ctx, n.Fun)
	n1.Args, ys = r.rewriteExprList(ctx, n.Args)

	cnt, isFnCall := retCntOfFnExpr(r.pkg, n.Fun)
	if !isFnCall {
		// 类型转换
		// e.g. int(a)
		assert(len(xs) == 0)
		assert(len(n.Args) == 1)
	}

	// 两个特殊场景,  go 和 defer
	// go func(){}() , go func() int { ... }()
	// 不能把 funLit 变成 define, 语义就变了
	if isGoCall(ctx, n) || isDeferCall(ctx, n) {
		return n1, append(xs, ys...)
	}

	// 调用无返回值函数
	if cnt == 0 {
		p := ctx.unparenParentNode()
		assert(instanceof[*ast.ExprStmt](p))
		return n1, append(xs, ys...)
		// return r.emptyExpr(), append(append(xs, ys...), &ast.ExprStmt{X: n1})
	}

	// 简化多返回值的改写
	// 因为
	// 		consume2 := func(int, bool) {}
	//		produce2 := func() (int, bool) { return 42, true }
	//		consume2(produce2())
	// 所以可以将多返回值处理成单返回值
	// 		consume2 := func(int, bool) {}
	//		produce2 := func() (int, bool) { return 42, true }
	//		consume2(produce2())
	// ==>
	//		i, b := produce2()
	//		consume2(tuple2(i, b))
	x, assign := r.tupleAssign(ctx, n1, cnt)
	return x, concat(xs, ys, sliceOf(assign))
}

func (r *fileRewriter) rewriteTryCall(
	ctx *rewriteCtx,
	n *ast.CallExpr, // callSite
	tryFn string,
) (v ast.Expr, extract []ast.Stmt) {
	fn := r.checkAndGetEnclosingFn(ctx.fun, n)

	r.checkTryCall(fn, n, tryFn)

	rhs, xs := r.rewriteTryCallArgs(ctx, n, tryFn)

	vCnt := retCntOfTryFn(tryFn) - 1
	lhs := make([]ast.Expr, vCnt+1 /*err*/)
	switch {
	case vCnt == 0:
		v = r.emptyExpr()
	case instanceof[*ast.ExprStmt](ctx.unparenParentNode()):
		//	(Try(ret1Err[B]()))
		// -----------------------------------
		// 		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟰 := ret1Err[B]()
		//		if 𝗲𝗿𝗿𝟰 != nil {
		//			return 𝗲𝗿𝗿𝟰
		//		}
		//		(𝘃𝗮𝗹𝟮) // 语法错误, 所以需要 unparen
		// ==>
		// 		_, 𝗲𝗿𝗿𝟰 := ret1Err[B]()
		//		if 𝗲𝗿𝗿𝟰 != nil {
		//			return 𝗲𝗿𝗿𝟰
		//		}
		//		(Ø())
		// case parentIsExprStmt:
		for i := 0; i < vCnt; i++ {
			lhs[i] = ast.NewIdent("_")
		}
		v = r.emptyExpr()
	case vCnt == 1:
		valId := r.genValId(ctx.fun)
		lhs[0] = valId
		v = valId
	default:
		for i := 0; i < vCnt; i++ {
			lhs[i] = r.genValId(ctx.fun)
		}
		v = r.tupleExpr(lhs[:vCnt])
	}
	errId := r.genErrId(ctx.fun)
	lhs[vCnt] = errId

	// 简化多返回值的改写
	// 因为
	// 		consume2 := func(int, bool) {}
	//		produce2 := func() (int, bool) { return 42, true }
	//		consume2(produce2())
	// 所以可以将多返回值处理成单返回值
	// 		f1 := func() (int, bool, error) { }
	// 		f2 := func(int, bool) { }
	// 		f2(try(f1()))
	// 		i,b,e := f1()
	// 		if e != nil {}
	// 		f2(tuple2(i, b))

	var (
		retIds    []ast.Expr
		beforeRet ast.Stmt = &ast.EmptyStmt{}
	)
	if fn.errId == nil {
		// if err != nil {
		// 		return ...,err
		// }
		zeroDecls := r.mkRetZeroVarDecl(ctx.fun, fn.sig)
		zeroIds := sliceMap(zeroDecls, func(i int, x *ast.ValueSpec) ast.Expr {
			return x.Names[0]
		})
		retIds = append(zeroIds, errId)
	} else {
		// if err != nil {
		//		errR = err
		// 		return
		// }
		beforeRet = r.simpleAssign(fn.errId, errId, token.ASSIGN)
	}

	return v, append(xs,
		// Try(f()) 	=> a, err := f()
		// Try(v, e) 	=> a, err = v, e
		&ast.AssignStmt{
			Lhs: lhs,
			Tok: token.DEFINE,
			Rhs: rhs,
		},
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X:  errId,
				Op: token.NEQ,
				Y:  ast.NewIdent("nil"), // shadow 已经 check 过
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					beforeRet,
					&ast.ReturnStmt{
						Results: retIds,
					},
				},
			},
		},
	)
}

func (r *fileRewriter) rewriteTryCallArgs(
	ctx *rewriteCtx,
	n *ast.CallExpr, // callSite
	tryFn string,
) ([]ast.Expr, []ast.Stmt) {
	// fast routine
	tryWithOneArgNonTryFn := func() (bool, *ast.CallExpr) {
		if len(n.Args) == 1 {
			// e.g. _ = Try((ret1Err[int]()))
			if c, ok := astutil.Unparen(n.Args[0]).(*ast.CallExpr); ok {
				return r.tryCallee(c) == "", c
			}
		}
		return false, nil
	}

	is, nonTryFnCall := tryWithOneArgNonTryFn()
	if !is {
		// 正常路径
		return r.rewriteExprList(ctx, n.Args)
	}

	// 特化路径:
	// 处理 Try(non_try_fun()) 形式, 避免生成过多临时变量
	// e.g. Try(non_try_fun()) =>
	// 正常路径, non_try_fun() 会通过 rewriteNonTryCall 生成一遍临时变量
	// 	a,b = non_try_fun()
	// doRewriteTryFnCall 又会生成一遍
	// 	c,d = T2(a,b)
	// 特化处理的思路: non_try_fun() 不生成 define, 直接 inline 处理成 callExpr
	cnt, isFnCall := retCntOfFnExpr(r.pkg, nonTryFnCall.Fun)
	if isFnCall {
		// non_try_fun()函数一定有返回值 至少一个 err 返回值
		assert(cnt > 0)
		assert(retCntOfTryFn(tryFn) == cnt) // cnt - err
	} else { // type cast
		// assert(len(xs) == 0)
		assert(len(n.Args) == 1)
	}

	// 注意这里不生成 define, 原地修改, 从而避免生成过多局部变量
	var xs, ys []ast.Stmt
	n1 := &ast.CallExpr{}
	n1.Fun, xs = r.rewriteExpr(ctx, nonTryFnCall.Fun)
	n1.Args, ys = r.rewriteExprList(ctx, nonTryFnCall.Args)
	args := []ast.Expr{n1}
	return args, append(xs, ys...)
}

func (r *fileRewriter) checkAndGetEnclosingFn(
	enFn fnNode,
	callSitePos loader.Positioner,
) (fn *enclosingFn) {
	fn = r.fnSig[enFn]
	if fn != nil {
		return
	}

	var (
		fnTy *ast.FuncType
		sig  *types.Signature
	)
	switch n := enFn.(type) {
	case *ast.FuncLit:
		sig, _ = r.pkg.TypeOf(n).(*types.Signature) // underlying?
		fnTy = n.Type
	case *ast.FuncDecl:
		sig, _ = r.pkg.TypeOf(n.Name).(*types.Signature) // underlying?
		fnTy = n.Type
	case nil:
		r.assert(callSitePos, false, "try must be called in the tryable func (return error in last position)")
	}

	var retErr types.Type
	// 检查函数最后一个返回值必须为 error
	{
		retCnt := sig.Results().Len()
		r.assert(fnTy, retCnt > 0, "expect at least one error return")
		retErr = sig.Results().At(retCnt - 1).Type()
		// r.assert(callSitePos, types.AssignableTo(retErr, errTy), "the last return type MUST assignable to error")
		r.assert(fnTy, types.Identical(retErr, r.errTy), "the last result must be error, but %s", retErr)
	}

	var errId *ast.Ident
	{
		xs := fnTy.Results.List
		assert(len(xs) > 0)
		last := xs[len(xs)-1]
		if len(last.Names) > 0 {
			err := last.Names[len(last.Names)-1]
			if err.Name != "_" {
				errId = err
			}
		}
	}

	fn = &enclosingFn{
		sig:   sig,
		errTy: retErr,
		errId: errId,
	}
	r.fnSig[enFn] = fn
	return fn
}

func (r *fileRewriter) checkTryCall(
	fn *enclosingFn,
	n *ast.CallExpr, // callSite
	tryFn string,
) {
	// 检查参数列表最后一个必须为 error
	// Try(err)    				OR		func F() error ; 			Try(F())
	// Try(v1, err) 			OR 		func F() (A, error); 		Try(F())
	// Try(v1, v2, err) 		OR 		func F() (A, B, error); 	Try(F())
	// Try(v1, v2, v3, err) 	OR 		func F() (A, B, C, error); 	Try(F())
	var argErr types.Type
	for i, name := range tryFnNames {
		lastParam, paramCnt := i, i+1
		if tryFn == name {
			argCnt := len(n.Args)
			r.assert(n, argCnt > 0, "at least one arg required")
			switch argCnt {
			case paramCnt:
				lastArg := n.Args[argCnt-1]
				argErr = r.pkg.TypeOf(lastArg)
			case 1:
				tup, _ := r.pkg.TypeOf(n.Args[0]).(*types.Tuple)
				r.assert(n, tup != nil && tup.Len() == paramCnt, "invalid args, expect %d params", paramCnt)
				argErr = tup.At(lastParam).Type()
			}
			break
		}
	}
	r.assert(n, argErr != nil, "invalid try args")
	// 因为需要用 if err != nil { } 所以要求返回的 err 必须可以与 nil 判等
	r.assert(n, isValidNil(argErr), "error <%s> must can be equals by nil", argErr.String())
	// r.assert(n, types.AssignableTo(argErr, retErr), "type mismatch, Try(..., ?) expect %v but %v", retErr, argErr)
	r.assert(n, types.AssignableTo(argErr, r.errTy), "type mismatch, Try(..., ?) expect %v but %v", fn.errTy, argErr)
}

func (r *fileRewriter) mkRetZeroVarDecl(fn fnNode, sig *types.Signature) []*ast.ValueSpec {
	// 一个函数内部会有多个 try 调用, 只处理一次返回值声明
	if xs, ok := r.fnZero[fn]; ok {
		return xs
	}

	retCnt := sig.Results().Len()
	if retCnt <= 1 /*error result*/ {
		r.fnZero[fn] = nil
		return nil
	}

	specs := make([]*ast.ValueSpec, 0, retCnt-1)

	fnTy := flattenFuncResults(fn)
	for i := 0; i < retCnt-1; i++ {
		zero := r.mkSym(fn, fmt.Sprintf("%s%d", valZeroIdentPrefix, i))
		ty := fnTy[i].Type
		specs = append(specs, &ast.ValueSpec{
			Names: []*ast.Ident{zero},
			Type:  ty,
		})
	}

	r.fnZero[fn] = specs
	return specs
}

func flattenFuncResults(f fnNode) []*ast.Field {
	ty, _ := unpackFunc(f)
	return flattenNamesByType(ty.Results.List, func(n *ast.Ident, ty ast.Expr) *ast.Field {
		return &ast.Field{Names: []*ast.Ident{n}, Type: ty}
	})
}

func retCntOfFnExpr(pkg loader.Pkg, f ast.Expr) (int, bool) {
	// e.g. (*int)(nil)
	// e.g. (f)()
	f = astutil.Unparen(f)
	switch sig := pkg.TypeOf(f).Underlying().(type) {
	case *types.Signature:
		return sig.Results().Len(), true
	default: // T(x) 类型转换
		assert(isTypeName(pkg.TypesInfo, f))
		return 1, false
	}
}
