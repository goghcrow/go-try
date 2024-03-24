package rewriter

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ast/astutil"
)

func isPkgSel(info *types.Info, sel *ast.SelectorExpr) bool {
	_, ok := info.Selections[sel]
	if ok {
		return false
	}
	return info.TypeOf(sel.X) == types.Typ[types.Invalid]
}

func isTypeName(info *types.Info, x ast.Expr) (ok bool) {
	switch x := x.(type) {
	case *ast.Ident:
		_, ok = info.ObjectOf(x).(*types.TypeName)
	case *ast.SelectorExpr:
		_, ok = info.ObjectOf(x.Sel).(*types.TypeName)
	case *ast.ArrayType, *ast.StructType, *ast.FuncType, *ast.InterfaceType, *ast.MapType, *ast.ChanType:
		ok = true
	case *ast.StarExpr:
		return isTypeName(info, x.X)
	}
	return
}

func isMapIndex(info *types.Info, n *ast.IndexExpr) (*types.Map, bool) {
	return typeOf[*types.Map](info.TypeOf(n.X), true)
}

func allTypeNames(info *types.Info, xs ...ast.Expr) bool {
	for _, x := range xs {
		if !isTypeName(info, x) {
			return false
		}
	}
	return true
}

func isNode[T ast.Node](n ast.Node) bool {
	_, ok := n.(T)
	return ok
}

func typeOf[T types.Type](ty types.Type, underlying bool) (z T, is bool) {
	if ty == nil {
		return z, false
	}
	if underlying {
		ty = ty.Underlying()
	}
	t, ok := ty.(T)
	return t, ok
}

func isType[T types.Type](ty types.Type, underlying bool) (is bool) {
	_, is = typeOf[T](ty, underlying)
	return
}

func typeOfStructPtr(ty types.Type, underlying bool) (*types.Struct, bool) {
	ptr, ok := typeOf[*types.Pointer](ty, underlying)
	if !ok {
		return nil, false
	}
	return typeOf[*types.Struct](ptr.Elem(), underlying)
}

func isIface(ty types.Type, underlying bool) bool {
	return isType[*types.Interface](ty, underlying)
}

func isAny(ty types.Type, underlying bool) bool {
	iface, ok := typeOf[*types.Interface](ty, underlying)
	return ok && iface.Empty()
}

func isBool(ty types.Type, underlying bool) bool {
	if b, ok := typeOf[*types.Basic](ty, underlying); ok {
		return b.Kind() == types.Bool
	}
	return false
}

func isRecvOperation(e ast.Expr) bool {
	u, _ := astutil.Unparen(e).(*ast.UnaryExpr) // rhs
	return u != nil && u.Op == token.ARROW
}

func isGoCall(ctx *rewriteCtx, n *ast.CallExpr) bool {
	if g, ok := ctx.unparenParentNode().(*ast.GoStmt); ok {
		return g.Call == n
	}
	return false
}

func isDeferCall(ctx *rewriteCtx, n *ast.CallExpr) bool {
	if d, ok := ctx.unparenParentNode().(*ast.DeferStmt); ok {
		return d.Call == n
	}
	return false
}

func isCallFun(ctx *rewriteCtx, n *ast.SelectorExpr) bool {
	if c, ok := ctx.unparenParentNode().(*ast.CallExpr); ok {
		if c.Fun == n {
			return true
		}
	}
	return false
}

func isStructField(x types.Object, st *types.Struct) bool {
	for i := 0; i < st.NumFields(); i++ {
		if st.Field(i) == x {
			return true
		}
	}
	return false
}

func isTuple2AssignRhs[T any](n ast.Expr, lhs []T, rhs []ast.Expr) bool {
	if len(lhs) == 2 && len(rhs) == 1 {
		return astutil.Unparen(rhs[0]) == n
	}
	return false
}

func isTuple2Assign(ctx *rewriteCtx, n ast.Expr) bool {
	switch p := ctx.unparenParentNode().(type) {
	case nil:
		assert(false)
		return false
	case *ast.ValueSpec:
		return isTuple2AssignRhs(n, p.Names, p.Values)
	case *ast.AssignStmt:
		return isTuple2AssignRhs(n, p.Lhs, p.Rhs)
	default:
		return false
	}
}
