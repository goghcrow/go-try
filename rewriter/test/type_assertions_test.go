//go:build !try

// Code generated by github.com/goghcrow/go-try DO NOT EDIT.
package test

import (
	"fmt"
	"go/ast"
)

func type_assertions() error {
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[ast.Node]()
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		expr, ok := 𝘃𝗮𝗹𝟭.(ast.Expr)
		_, _ = expr, ok
	}
	{
		var (
			expr ast.Expr
			ok   bool
		)
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := ret1Err[ast.Node]()
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
		expr, ok = 𝘃𝗮𝗹𝟮.(ast.Expr)
		_, _ = expr, ok
	}
	{
		𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟯 := ret1Err[ast.Node]()
		if 𝗲𝗿𝗿𝟯 != nil {
			return 𝗲𝗿𝗿𝟯
		}
		var expr, ok = 𝘃𝗮𝗹𝟯.(ast.Expr)
		_, _ = expr, ok
	}
	{
		𝘃𝗮𝗹𝟰, 𝗲𝗿𝗿𝟰 := ret1Err[ast.Node]()
		if 𝗲𝗿𝗿𝟰 != nil {
			return 𝗲𝗿𝗿𝟰
		}
		var expr, ok = 𝘃𝗮𝗹𝟰.(ast.Expr)
		_, _ = expr, ok
	}
	{
		𝘃𝗮𝗹𝟱, 𝗲𝗿𝗿𝟱 := ret1Err[ast.Node]()
		if 𝗲𝗿𝗿𝟱 != nil {
			return 𝗲𝗿𝗿𝟱
		}
		var expr, ok interface{} = 𝘃𝗮𝗹𝟱.(ast.Expr)
		_, _ = expr, ok
	}
	{
		𝘃𝗮𝗹𝟲, 𝗲𝗿𝗿𝟲 := ret1Err[ast.Node]()
		if 𝗲𝗿𝗿𝟲 != nil {
			return 𝗲𝗿𝗿𝟲
		}
		var expr, ok any = 𝘃𝗮𝗹𝟲.(ast.Expr)
		_, _ = expr, ok
	}
	{
		𝘃𝗮𝗹𝟳, 𝗲𝗿𝗿𝟳 := ret1Err[ast.Node]()
		if 𝗲𝗿𝗿𝟳 != nil {
			return 𝗲𝗿𝗿𝟳
		}
		𝘃𝗮𝗹𝟴 := 𝘃𝗮𝗹𝟳.(ast.Expr)
		expr := 𝘃𝗮𝗹𝟴
		_ = expr
	}
	{
		type Int int
		𝘃𝗮𝗹𝟵, 𝗲𝗿𝗿𝟴 := ret1Err[int]()
		if 𝗲𝗿𝗿𝟴 != nil {
			return 𝗲𝗿𝗿𝟴
		}
		𝘃𝗮𝗹𝟭𝟬 := Int(𝘃𝗮𝗹𝟵)
		_ = 𝘃𝗮𝗹𝟭𝟬
	}
	{
		𝘃𝗮𝗹𝟭𝟭, 𝗲𝗿𝗿𝟵 := ret1Err[error]()
		if 𝗲𝗿𝗿𝟵 != nil {
			return 𝗲𝗿𝗿𝟵
		}
		switch 𝘃𝗮𝗹𝟭𝟭.(type) {
		}
	}
	{
		{
			n := 1
			𝘃𝗮𝗹𝟭𝟮, 𝗲𝗿𝗿𝟭𝟬 := func1[int, fmt.Stringer](n)
			if 𝗲𝗿𝗿𝟭𝟬 != nil {
				return 𝗲𝗿𝗿𝟭𝟬
			}
			switch 𝘃𝗮𝗹𝟭𝟮.(type) {
			}
		}
	}
	return nil
}
