//go:build try

package test

import (
	. "github.com/goghcrow/go-try"
)

func range_without_kv() error {
	for range Try(ret1Err[[]int]()) {
	}
	return nil
}

func range_selector_assign() error {
	type X struct{ i int }
	type Y = X

	{
		for Try(ret1Err[*X]()).i, Try(ret1Err[*Y]()).i = range []int{} {
			println()
		}
	}
	{
		for id[*X](nil).i, Try(ret1Err[*Y]()).i = range []int{} {
			println()
		}
	}
	{
		for Try(ret1Err[*X]()).i, id[*Y](nil).i = range []int{} {
			println()
		}
	}
	{
		var i int
		for i, Try(ret1Err[*X]()).i = range []int{} {
			println(i)
		}
	}
	{
		var v int
		for Try(ret1Err[*X]()).i, v = range []int{} {
			println(v)
		}
	}
	{
		for Try(ret1Err[*X]()).i = range []int{} {
		}
	}

	return nil
}

func range_index_assign() error {
	type X struct{ i int }

	{
		for Try(ret1Err[[]int]())[0], Try(ret1Err[[]int]())[0] = range []int{} {
			println()
		}
	}
	{
		for []int{}[0], Try(ret1Err[[]int]())[0] = range []int{} {
			println()
		}
	}
	{
		for Try(ret1Err[[]int]())[0], []int{}[0] = range []int{} {
			println()
		}
	}
	{
		var i int
		for i, Try(ret1Err[[]int]())[0] = range []int{} {
			println(i)
		}
	}
	{
		var v int
		for Try(ret1Err[[]int]())[0], v = range []int{} {
			println(v)
		}
	}
	{
		for Try(ret1Err[[]int]())[0] = range []int{} {
		}
	}

	return nil
}

func range_map_index_assign() error {
	{
		for Try(ret1Err[map[int]int]())[0], Try(ret1Err[map[int]int]())[0] = range []int{} {
			println()
		}
	}
	{
		for map[int]int{}[0], Try(ret1Err[map[int]int]())[0] = range []int{} {
			println()
		}
	}
	{
		for Try(ret1Err[map[int]int]())[0], map[int]int{}[0] = range []int{} {
			println()
		}
	}
	{
		var i int
		for i, Try(ret1Err[map[int]int]())[0] = range []int{} {
			println(i)
		}
	}
	{
		var v int
		for Try(ret1Err[map[int]int]())[0], v = range []int{} {
			println(v)
		}
	}
	{
		for Try(ret1Err[map[int]int]())[0] = range []int{} {
		}
	}
	return nil
}
