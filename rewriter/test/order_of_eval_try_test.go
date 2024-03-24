//go:build try

package test

import (
	"fmt"
	"go/ast"

	. "github.com/goghcrow/go-try"
)

func rewrite_binary_logical_expr() error {
	type X struct{ a int }
	yyy := func(a, b int) bool {
		return true
	}
	n := 42

	if yyy(X{id(42)}.a, Try(ret1Err[int]())) &&
		n > 100 &&
		Try(ret1Err[bool]()) &&
		yyy(1, 2) {
		println(1)
	}
	return nil
}

func rewrite_assign_and_value_spec() error {
	type (
		A = int
		B = int
	)
	{
		var a, b int
		a, b = Try(ret1Err[A]()), Try(ret1Err[B]())
		_, _ = a, b
	}
	{
		a, b := Try(ret1Err[A]()), Try(ret1Err[B]())
		_, _ = a, b
	}

	{
		var a, b int
		a, b = id[A](42), Try(ret1Err[B]())
		_, _ = a, b
	}
	{
		a, b := id[A](42), Try(ret1Err[B]())
		_, _ = a, b
	}
	return nil
}

func rewrite_compositelit_expr() error {
	type S struct {
		name string
		age  int
	}
	_ = S{
		name: id[string]("hello"),
		age:  id[int](0) + Try(ret1Err[int]()),
	}

	_ = S{
		id[string]("hello"),
		id[int](0) + Try(ret1Err[int]()),
	}
	return nil
}

func rewrite_return() (a, b int, err error) {
	type (
		A = int
		B = int
	)
	return Try(ret1Err[A]()), Try(ret1Err[B]()), nil
}

func rewrite_return1() (a, b int, err error) {
	type (
		A = int
		B = int
	)
	return id[A](42), Try(ret1Err[B]()), nil
}

func rewrite_star_expr() error {
	var ptr *int
	println(*ptr + Try(ret1Err[int]()))
	return nil
}

func rewrite_slice_expr() error {
	type (
		A = int
		B = int
		C = int
		D = int
	)
	{
		println([]int{}[100] + Try(ret1Err[int]()))
	}
	{
		println([]int{}[1:2:3], Try(ret1Err[int]()))
	}
	{
		println([]int{}[1:Try(ret1Err[A]()):3], Try(ret1Err[B]()))
	}
	{
		println([]int{}[Try(ret1Err[A]()):Try(ret1Err[B]()):Try(ret1Err[C]())], Try(ret1Err[D]()))
	}
	return nil
}

func rewrite_binary_expr() error {
	type (
		fst = int
		snd = int
	)
	{
		println(Try(ret1Err[fst]()) + Try(ret1Err[snd]()))
	}
	{
		println(return1[int]() + Try(ret1Err[int]()))
	}

	{
		println(fmt.Sprintf("") + Try(ret1Err[string]()))
	}

	{
		var a any
		println(a.(int) + Try(ret1Err[int]()))
	}

	return nil
}

func rewrite_call_args() error {
	{
		var a any
		println(a.(int), Try(ret1Err[int]()))
	}
	return nil
}

func rewrite_index_expr() error {
	{
		var x []int
		_ = x[1] + Try(ret1Err[int]())
	}
	{
		_ = []int{}[Try(ret1Err[int]())]
	}

	{
		_ = Try(ret1Err[[]func(int) int]())[Try(ret1Err[int]())]
	}

	{
		_ = Try(ret1Err[[]func(int) int]())[Try(ret1Err[int]())](42) + Try(ret1Err[int]())
	}

	{
		_ = id[int](42) + Try(ret1Err[int]())
	}
	{
		_ = &[]func(){}[0] == &[]func(){}[Try(ret1Err[int]())]
	}
	return nil
}

func rewrite_index1_expr[T int]() (a int, err error) {
	{
		_ = id[T](42) + Try(ret1Err[T]())
	}
	{
		_ = &[]func(){}[0] == &[]func(){}[Try(ret1Err[T]())]
	}
	return
}

func rewrite_type_assertion_expr() (err error) {
	{
		var n ast.Node
		_ = n.(*ast.Ident).Name + " " + Try(ret1Err[string]())
	}
	return nil
}

func rewrite_mixed() error {
	type (
		A = int
		B = int
		C = int
	)

	{
		// 必须先求值 id[int](0)+Try(ret1Err[A]())
		// 再求值 []func(int) int{}[id[int](0)+Try(ret1Err[A]())]
		// 再求值 id[int](1)+Try(ret1Err[B]())
		// 再求值 []func(int) int{}[id[int](0)+Try(ret1Err[A]())](
		//				id[int](1)+Try(ret1Err[B]()),
		//			)
		// 再求值 Try(ret1Err[C]()) + id[int](2),
		// ...
		println(
			[]func(int) int{}[id[int](0)+Try(ret1Err[A]())](
				id[int](1)+Try(ret1Err[B]()),
			) +
				Try(ret1Err[C]()) + id[int](2),
		)
	}

	{
		Try0(Try(ret1Err[error]()))
	}
	return nil
}

func rewrite_if_init_cond() error {
	{
		if n := 1; Try(func1[int, bool](n)) {
		}
	}
	{
		n := 0
		if n := 1; Try(func1[int, bool](n)) {
			n++
		}
		println(n)
	}
	{
		n := 0
		if n := Try(func1[int, int](42)); Try(func1[int, bool](n)) {
			n++
		}
		println(n)
	}
	return nil
}

func rewrite_for_init_cond() error {
	{
		for i := 0; Try(func1[int, bool](i)); Try(func1[int, int](i)) {
			i++
		}
	}
	{
		i := 0
		for i := 0; Try(func1[int, bool](i)); Try(func1[int, int](i)) {
			i++
		}
		println(i)
	}
	{
		i := 0
		for i := Try(func1[int, int](42)); Try(func1[int, bool](i)); Try(func1[int, int](i)) {
			i++
		}
		println(i)
	}
	return nil
}

func rewrite_typeswitch_init_assign() error {
	{
		switch n := 1; Try(func1[int, fmt.Stringer](n)).(type) {
		}
	}
	{
		n := 0
		switch n := 1; Try(func1[int, fmt.Stringer](n)).(type) {
		default:
			n++
		}
		println(n)
	}
	{
		n := 0
		switch n := Try(ret1Err[int]()); Try(func1[int, fmt.Stringer](n)).(type) {
		default:
			n++
		}
		println(n)
	}
	return nil
}

func rewrite_switch_init_tag() error {
	{
		switch n := 1; Try(func1[int, fmt.Stringer](n)) {
		}
	}
	{
		n := 0
		switch n := 1; Try(func1[int, fmt.Stringer](n)) {
		default:
			n++
		}
		println(n)
	}
	{
		n := 0
		switch n := Try(ret1Err[int]()); Try(func1[int, fmt.Stringer](n)) {
		default:
			n++
		}
		println(n)
	}
	return nil
}
