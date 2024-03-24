package rewriter

import (
	"go/ast"
)

func (r *fileRewriter) rewriteLabeled(ctx *rewriteCtx, n *ast.LabeledStmt) (
	*ast.LabeledStmt,
	[]ast.Stmt,
) {
	ys, xs := r.rewriteLabeledEx(ctx, n)
	assert(len(ys) == 1)
	// 只有 StmtList 允许 labeledStmt
	// labeled for/range/switch/select 也返回 len(ys) == 1
	return ys[0].(*ast.LabeledStmt), xs
}

func (r *fileRewriter) rewriteLabeledEx(ctx *rewriteCtx, n *ast.LabeledStmt) (
	[]ast.Stmt, // *ast.LabeledStmt,
	[]ast.Stmt,
) {
	var xs []ast.Stmt
	n1 := &ast.LabeledStmt{Label: n.Label}
	assert(n.Stmt != nil)
	n1.Stmt, xs = r.rewriteStmt(ctx, n.Stmt)
	if len(xs) == 0 {
		return sliceOf[ast.Stmt](n1), xs
	}

	// 1. insert before 的 stmt 不能越过 label !!!!
	// 		否则会导致 goto 跳转, skip xs 中的 stmt
	// 2. labeld for/range/switch/select, 可能跟 break 和 continue 绑定
	// 		不能修改, 否则破坏语法
	// 3. 除了 labeled for/range/switch/select, 与 stmtList 其他位置也不允许出现 labeld
	// 		所以只需要在 rewriteStmtList 特殊处理 labeledStmt 即可
	switch n1.Stmt.(type) {
	case *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.SelectStmt:
		// 一定不会有 goto 跳转到该 label, 因为 prewriteLabeled(...) 已经处理了
		jTbl := r.jmpTbl(ctx.fun)
		for _, l := range jTbl {
			assert(l != n)
		}
		return sliceOf[ast.Stmt](n1), xs
	default:
		// L:
		//	var a = Try(ret1Err[int]())
		//	goto L
		//	println(a)
		// ==>
		// L:
		//	𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[int]()
		//	if 𝗲𝗿𝗿𝟭 != nil {
		//		return 𝗲𝗿𝗿𝟭
		//	}
		//	var a = 𝘃𝗮𝗹𝟭
		//	goto L
		//	println(a)
		st := n1.Stmt
		n1.Stmt = xs[0]
		return concat(sliceOf[ast.Stmt](n1), xs[1:], sliceOf(st)), nil
	}
}
