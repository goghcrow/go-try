//go:build try

package test

import (
	"go/ast"

	. "github.com/goghcrow/go-try"
)

func rewrite_selector_expr() error {
	type X struct{ x string }
	{
		var x *X
		_ = x.x + Try(func1[int, string](2))
	}

	{
		println(Try(ret1Err[*X]()).x + Try(ret1Err[string]()))
	}

	{
		println(Try(ret1Err[X]()).x + Try(ret1Err[string]()))
	}
	return nil
}

func rewrite_ptr_selector_expr() error {
	var x *ast.CallExpr
	{
		// 可能 panic
		consume2(x.Args, Try(ret1Err[string]()))
	}
	{
		// 不会 panic
		consume2(x.Pos, Try(ret1Err[string]()))
	}
	{
		// 可能 panic
		consume2(x.Pos(), Try(ret1Err[string]()))
	}
	return nil
}

func rewrite_iface_selector_expr() error {
	var x ast.Node
	{
		// 可能 panic
		consume2(x.Pos, Try(ret1Err[string]()))
	}
	{
		// 可能 panic
		consume2(x.Pos(), Try(ret1Err[string]()))
	}
	return nil
}
