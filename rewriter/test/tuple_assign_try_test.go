//go:build try

package test

import (
	. "github.com/goghcrow/go-try"
)

func tuple_assign_map_index() error {
	{
		_, _ = map[int]int{}[Try(ret1Err[int]())]
	}
	{
		v, ok := map[int]int{}[Try(ret1Err[int]())]
		_, _ = v, ok
	}
	{
		var v, ok = map[int]int{}[Try(ret1Err[int]())]
		_, _ = v, ok
	}
	{
		var v int
		v, _ = map[int]int{}[Try(ret1Err[int]())]
		_ = v
	}
	{
		var v int
		v, ok := map[int]int{}[Try(ret1Err[int]())]
		_, _ = v, ok
	}
	{
		var v int
		var ok bool
		v, ok = map[int]int{}[Try(ret1Err[int]())]
		_, _ = v, ok
	}

	return nil
}

func tuple_assign_any_map_index() error {
	{
		_, _ = map[any]int{}[Try(ret1Err[int]())]
	}
	{
		v, ok := map[any]int{}[Try(ret1Err[int]())]
		_, _ = v, ok
	}
	{
		var v, ok = map[any]int{}[Try(ret1Err[int]())]
		_, _ = v, ok
	}
	{
		var v int
		v, _ = map[any]int{}[Try(ret1Err[int]())]
		_ = v
	}
	{
		var v int
		v, ok := map[any]int{}[Try(ret1Err[int]())]
		_, _ = v, ok
	}
	{
		var v int
		var ok bool
		v, ok = map[any]int{}[Try(ret1Err[int]())]
		_, _ = v, ok
	}

	return nil
}
