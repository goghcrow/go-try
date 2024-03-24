//go:build try

package test

import (
	"fmt"
	"go/ast"

	. "github.com/goghcrow/go-try"
)

func type_assertions() error {
	{
		expr, ok := Try(ret1Err[ast.Node]()).(ast.Expr)
		_, _ = expr, ok
	}

	{
		var (
			expr ast.Expr
			ok   bool
		)
		expr, ok = Try(ret1Err[ast.Node]()).(ast.Expr)
		_, _ = expr, ok
	}

	{
		var expr, ok = Try(ret1Err[ast.Node]()).(ast.Expr)
		_, _ = expr, ok
	}

	{
		var expr, ok = Try(ret1Err[ast.Node]()).(ast.Expr)
		_, _ = expr, ok
	}

	{
		var expr, ok interface{} = Try(ret1Err[ast.Node]()).(ast.Expr)
		_, _ = expr, ok
	}

	{
		var expr, ok any = Try(ret1Err[ast.Node]()).(ast.Expr)
		_, _ = expr, ok
	}

	{
		expr := Try(ret1Err[ast.Node]()).(ast.Expr)
		_ = expr
	}

	{
		type Int int
		_ = Int(Try(ret1Err[int]()))
	}

	{
		switch Try(ret1Err[error]()).(type) {
		}
	}
	{
		switch n := 1; Try(func1[int, fmt.Stringer](n)).(type) {
		}
	}
	return nil
}
