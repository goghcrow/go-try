//go:build try

package test

import (
	"fmt"

	. "github.com/goghcrow/go-try"
)

func selecotr_val() (_ int, err error) {
	type X struct{ x int }
	var x X
	return x.x + Try(ret1Err[int]()), nil
}

func panic_selecotr_ptr() (_ int, err error) {
	type X struct{ x int }
	var x *X
	return x.x + Try(ret1Err[int]()), nil
}

func panic_selecotr_iface() (_ string, err error) {
	var s fmt.Stringer
	return s.String() + Try(ret1Err[string]()), nil
}

func panic_quo(zero int) (_ int, err error) {
	return 42/zero + Try(ret1Err[int]()), nil
}

func panic_shl(negative int) (_ int, err error) {
	return 42<<negative + Try(ret1Err[int]()), nil
}

func panicRet[T any]() (z T) {
	panic("panic")
	return z
}

func panic_Panic() (_ int, err error) {
	return panicRet[int]() + Try(ret1Err[int]()), nil
}

func panic_call_nil_fun() (_ int, err error) {
	var f func() int
	return f() + Try(ret1Err[int]()), nil
}

func panic_call_nil_fun1() (_ int, err error) {
	type X struct {
		f func() int
	}
	var x *X
	return x.f() + Try(ret1Err[int]()), nil
}

func panic_type_conv(x any) (_ int, err error) {
	return x.(int) + Try(ret1Err[int]()), nil
}

func panic_type_conv1() (_ int, err error) {
	consume2((*[1]int)([]int{}), Try(ret1Err[int]()))
	return
}

func panic_type_conv2() (_ int, err error) {
	consume2((*[1]int)([]int{}), Try(ret1Err[int]()))
	return
}

func panic_make(negitive int) (_ int, err error) {
	consume2(make([]int, negitive), Try(ret1Err[int]()))
	return
}
