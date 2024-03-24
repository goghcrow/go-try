//go:build try

package test

import (
	. "github.com/goghcrow/go-try"
)

func logical_or() error {
	{
		_ = id(true) || id(false)
	}

	{
		_ = Try(func1[int, bool](1)) || Try(func1[int, bool](2))
	}

	{
		_ = id(true) || Try(func1[int, bool](2))
	}

	{
		_ = id(true) || id(false) || Try(func1[int, bool](2))
	}

	{
		_ = id(true) || Try(func1[int, bool](1)) || Try(func1[int, bool](2))
	}

	{
		_ = id(true) || Try(func1[int, bool](2)) || id(false)
	}

	return nil
}

func logical_and() error {
	{
		_ = id(true) && id(false)
	}

	{
		_ = Try(func1[int, bool](1)) && Try(func1[int, bool](2))
	}

	{
		_ = id(true) && Try(func1[int, bool](2))
	}

	{
		_ = id(true) && id(false) && Try(func1[int, bool](2))
	}

	{
		_ = id(true) && Try(func1[int, bool](1)) && Try(func1[int, bool](2))
	}

	{
		_ = id(true) && Try(func1[int, bool](2)) && id(false)
	}

	return nil
}

func logical_and_or() error {
	{
		_ = id(true) && id(false) || Try(func1[int, bool](2))
	}

	{
		_ = id(true) && (id(false) || Try(func1[int, bool](2)))
	}

	{
		_ = id(true) || Try(func1[int, bool](1)) && Try(func1[int, bool](2))
	}

	{
		_ = (id(true) || Try(func1[int, bool](1))) && Try(func1[int, bool](2))
	}

	return nil
}

func logical_special_case0() error {
	{
		_ = []int{}[0] > 0 || Try(func1[int, bool](1))
	}
	{
		_ = (&(struct{ a int }{})).a > 0 || Try(func1[int, bool](1))
	}
	{
		_ = struct{ a int }{}.a > 0 || Try(func1[int, bool](1))
	}
	{
		var a any
		_ = a.(bool) || Try(func1[int, bool](1))
	}
	{
		var a *bool
		_ = *a || Try(func1[int, bool](1))
	}

	return nil
}

func logical_special_case() error {
	{
		_ = 1 == 1 || Try(func1[int, bool](1))
	}
	{
		_ = Try(func1[int, bool](1)) || 1 == 1
	}
	{
		_ = 1 == 1 && Try(func1[int, bool](1))
	}
	{
		_ = Try(func1[int, bool](1)) && 1 == 1
	}
	{
		_ = 1 == 1 || 1 == 1 && Try(func1[int, bool](2))
	}
	{
		_ = Try(func1[int, bool](1)) || 1 == 1 || Try(func1[int, bool](2))
	}
	return nil
}
