package rewriter

import (
	"go/ast"
	"go/token"
)

func (r *fileRewriter) rewriteFun(ctx *rewriteCtx, fn fnNode) (x fnNode, xs []ast.Stmt) {
	var (
		zeroRetDecl = func(fn fnNode) ast.Stmt {
			z := r.fnZero[fn]
			if z == nil {
				return &ast.EmptyStmt{}
			}
			return &ast.DeclStmt{
				Decl: &ast.GenDecl{
					Tok: token.VAR,
					Specs: groupNamesByType[ast.Spec](z, func(x, y ast.Expr) bool {
						return x == y
					}),
				},
			}
		}
		rewriteBody = func(body *ast.BlockStmt) *ast.BlockStmt {
			if body == nil {
				return nil
			}
			// 注意这里必须先处理 Body, 因为处理完 body, 才能取到 zeroRet
			block := &ast.BlockStmt{}
			block.List = r.rewriteStmtList(ctx.withFun(fn), body.List)
			block.List = prepend(block.List, zeroRetDecl(fn))
			return block
		}
	)

	switch n := fn.(type) {
	case *ast.FuncDecl:
		n1 := &ast.FuncDecl{
			Recv: n.Recv,
			Name: n.Name,
			Type: n.Type,
			Body: rewriteBody(n.Body),
		}
		return n1, nil
	case *ast.FuncLit:
		n1 := &ast.FuncLit{
			Type: n.Type,
			Body: rewriteBody(n.Body),
		}
		return n1, nil
	default:
		panic("illegal state")
	}
}
