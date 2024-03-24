//go:build try

package test

import (
	"go/ast"

	. "github.com/goghcrow/go-try"
)

func lhs_func_call_assign() (i int, err error) {
	*id(&i) = Try(ret1Err[int]())

	{
		(*id(&i)) = (Try(ret1Err[int]()))
	}
	return
}

func selector_assign_expr() error {
	type X struct{ x int }
	{
		var x X
		x.x = Try(ret1Err[int]())
	}

	{
		var x *X
		x.x = Try(ret1Err[int]())
	}

	// ========================================

	{

		Try(ret1Err[*X]()).x = 42
	}

	{

		(Try(ret1Err[*X]()).x) = 42
	}

	{
		id[*X](nil).x = Try(ret1Err[int]())
	}

	{
		Try(ret1Err[*X]()).x = Try(ret1Err[int]())
	}
	return nil
}

func index_assign_expr() error {
	{
		Try(ret1Err[[]int]())[0] = 42
	}

	{
		id[[]int](nil)[0] = Try(ret1Err[int]())
	}

	{
		Try(ret1Err[[]int]())[0] = Try(ret1Err[int]())
	}

	{
		Try(ret1Err[map[int]string]())[0] = "Hello"
	}

	{
		map[int]string{}[0] = Try(ret1Err[string]())
	}

	{
		Try(ret1Err[map[int]string]())[0] = Try(ret1Err[string]())
	}

	return nil
}

func if_init_assign_expr() (err error) {
	type X struct{ x int }

	{
		if id[*X](nil).x = 42; Try(ret1Err[bool]()) {
		}
	}
	{
		if Try(ret1Err[*X]()).x = 42; id[bool](true) {
		}
	}
	{
		if Try(ret1Err[*X]()).x = 42; Try(ret1Err[bool]()) {
		}
	}

	{
		if id[[]int](nil)[0] = 42; Try(ret1Err[bool]()) {
		}
	}
	{
		if Try(ret1Err[[]int]())[0] = 42; id[bool](true) {
		}
	}
	{
		if Try(ret1Err[[]int]())[0] = 42; Try(ret1Err[bool]()) {
		}
	}

	{
		if map[int]string{}[0] = "hello"; Try(ret1Err[bool]()) {
		}
	}
	{
		if Try(ret1Err[map[int]string]())[0] = "hello"; id[bool](true) {
		}
	}
	{
		if Try(ret1Err[map[int]string]())[0] = "hello"; Try(ret1Err[bool]()) {
		}
	}
	return nil
}

func switch_init_assign_expr() (err error) {
	type X struct{ x int }

	{
		switch id[*X](nil).x = 42; Try(ret1Err[bool]()) {
		}
	}
	{
		switch Try(ret1Err[*X]()).x = 42; id[int](42) {
		}
	}
	{
		switch Try(ret1Err[*X]()).x = 42; Try(ret1Err[bool]()) {
		}
	}

	{
		switch id[[]int](nil)[0] = 42; Try(ret1Err[bool]()) {
		}
	}
	{
		switch Try(ret1Err[[]int]())[0] = 42; id[int](42) {
		}
	}
	{
		switch Try(ret1Err[[]int]())[0] = 42; Try(ret1Err[bool]()) {
		}
	}

	{
		switch map[int]string{}[0] = "hello"; Try(ret1Err[bool]()) {
		}
	}
	{
		switch Try(ret1Err[map[int]string]())[0] = "hello"; id[int](42) {
		}
	}
	{
		switch Try(ret1Err[map[int]string]())[0] = "hello"; Try(ret1Err[bool]()) {
		}
	}

	return nil
}

func type_switch_init_assign_expr() (err error) {
	type X struct{ x int }
	{
		switch id[*X](nil).x = 42; n := Try(ret1Err[ast.Node]()).(type) {
		default:
			_ = n
		}
	}
	{
		switch Try(ret1Err[*X]()).x = 42; n := ast.Node(nil).(type) {
		default:
			_ = n
		}
	}
	{
		switch Try(ret1Err[*X]()).x = 42; n := Try(ret1Err[ast.Node]()).(type) {
		default:
			_ = n
		}
	}

	{
		switch id[[]int](nil)[0] = 42; n := Try(ret1Err[ast.Node]()).(type) {
		default:
			_ = n
		}
	}
	{
		switch Try(ret1Err[[]int]())[0] = 42; n := ast.Node(nil).(type) {
		default:
			_ = n
		}
	}
	{
		switch Try(ret1Err[[]int]())[0] = 42; n := Try(ret1Err[ast.Node]()).(type) {
		default:
			_ = n
		}
	}

	{
		switch map[int]string{}[0] = "hello"; n := Try(ret1Err[ast.Node]()).(type) {
		default:
			_ = n
		}
	}
	{
		switch Try(ret1Err[map[int]string]())[0] = "hello"; n := ast.Node(nil).(type) {
		default:
			_ = n
		}
	}
	{
		switch Try(ret1Err[map[int]string]())[0] = "hello"; n := Try(ret1Err[ast.Node]()).(type) {
		default:
			_ = n
		}
	}
	return nil
}

func for_switch_init_assign_expr() (err error) {
	type X struct{ x int }

	{
		for id[*X](nil).x = 42; Try(ret1Err[bool]()); {
		}
	}
	{
		for Try(ret1Err[*X]()).x = 42; id[bool](true); {
		}
	}
	{
		for Try(ret1Err[*X]()).x = 42; Try(ret1Err[bool]()); {
		}
	}

	{
		for id[[]int](nil)[0] = 42; Try(ret1Err[bool]()); {
		}
	}
	{
		for Try(ret1Err[[]int]())[0] = 42; id[bool](true); {
		}
	}
	{
		for Try(ret1Err[[]int]())[0] = 42; Try(ret1Err[bool]()); {
		}
	}

	{
		for map[int]string{}[0] = "hello"; Try(ret1Err[bool]()); {
		}
	}
	{
		for Try(ret1Err[map[int]string]())[0] = "hello"; id[bool](true); {
		}
	}
	{
		for Try(ret1Err[map[int]string]())[0] = "hello"; Try(ret1Err[bool]()); {
		}
	}
	return nil
}
