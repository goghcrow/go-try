package rewriter

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

func (r *fileRewriter) emptyExpr() ast.Expr {
	return r.tupleExpr(nil)
}

// 把多个表达式和成一个多值表达式
func (r *fileRewriter) tupleExpr(xs []ast.Expr) ast.Expr {
	switch len(xs) {
	case 1:
		return xs[0]
	default:
		r.importRT = true
		return &ast.CallExpr{
			Fun:  ast.NewIdent(tupleNames[len(xs)]),
			Args: xs,
		}
	}
}

func (r *fileRewriter) simpleAssign(lhs *ast.Ident, rhs ast.Expr, tok token.Token) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Tok: tok,
		Rhs: []ast.Expr{astutil.Unparen(rhs)},
	}
}

func (r *fileRewriter) assign(ctx *rewriteCtx, rhs ast.Expr) (ast.Expr, ast.Stmt) {
	return r.tupleAssign(ctx, rhs, 1)
}

// rhs: 多值表达式
// n: lhs var 数量
func (r *fileRewriter) tupleAssign(ctx *rewriteCtx, rhs ast.Expr, n int) (ast.Expr, ast.Stmt) {
	assert(n > 0)
	switch ctx.parent.node.(type) {
	case *ast.ExprStmt:
		return rhs, &ast.EmptyStmt{}
	default:
		lhs := make([]ast.Expr, n)
		for i := 0; i < n; i++ {
			lhs[i] = r.genValId(ctx.fun)
		}
		assign := &ast.AssignStmt{
			Lhs: lhs,
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				astutil.Unparen(rhs),
			},
		}
		return r.tupleExpr(lhs), assign
	}
}

func unpackFunc(n fnNode) (*ast.FuncType, *ast.BlockStmt) {
	switch n := n.(type) {
	case *ast.FuncLit:
		return n.Type, n.Body
	case *ast.FuncDecl:
		return n.Type, n.Body
	default:
		panic("illegal state")
	}
}

func trimTrailingEmptyStmts(xs []ast.Stmt) []ast.Stmt {
	for i := len(xs); i > 0; i-- {
		if _, ok := xs[i-1].(*ast.EmptyStmt); !ok {
			return xs[:i]
		}
	}
	return nil
}

func constFalse() ast.Expr {
	// false 可能被 shadow, 用 42 != 42 来表示 false
	x := &ast.BasicLit{Kind: token.INT, Value: "42"}
	return &ast.BinaryExpr{X: x, Op: token.NEQ, Y: x}
}

func constTrue() ast.Expr {
	// true 可能被 shadow, 用 42 == 42 来表示 false
	x := &ast.BasicLit{Kind: token.INT, Value: "42"}
	return &ast.BinaryExpr{X: x, Op: token.EQL, Y: x}
}
