//go:build try

package test

import (
	. "github.com/goghcrow/go-try"
)

func blank_in_range() error {
	for _, _ = range Try(ret1Err[[]int]()) {
	}
	for _ = range Try(ret1Err[[]int]()) {
	}
	return nil
}

func blank_in_map_index() error {
	{
		_ = map[int]int{}[Try(ret1Err[int]())]
	}
	// tuple assign
	{
		_, _ = map[int]int{}[Try(ret1Err[int]())]
	}
	// tuple assign + nil key map index
	{
		_, _ = map[any]int{}[Try(ret1Err[int]())]
	}
	// underlying
	{
		type Nil any
		_, _ = map[Nil]int{}[Try(ret1Err[int]())]
	}
	return nil
}

func blank_in_select() error {
	ch := (<-chan int)(nil)
	select {
	case _, _ = <-ch:
	case _, _ = <-Try(ret1Err[chan int]()):
	default:
	}
	return nil
}
