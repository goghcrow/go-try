//go:build try

package test

import (
	. "github.com/goghcrow/go-try"
)

func sssdfdfd() error {
	// 写入 nil map 异常
	{
		var m map[int]int
		m[Try(ret1Err[int]())] = 1
	}
	// 读 nil map 不会异常
	{
		var m map[int]int
		println(m[0], Try(ret1Err[int]()))
	}
	// map[any]T 读取异常
	{
		var m map[any]int
		println(m[0], Try(ret1Err[int]()))
	}
	return nil
}

func map_index() error {
	// 写入 nil map 异常
	{
		var m map[int]int
		m[Try(ret1Err[int]())] = 1
	}

	{
		// 读 nil map 不会异常
		{
			var m map[int]int
			println(m[0], Try(ret1Err[int]()))
		}
		// map[any]T 读取异常
		{
			var m map[any]int
			println(m[0], Try(ret1Err[int]()))
		}
	}

	{
		var m map[int]int
		println(m[func() int {
			println("called")
			return 0
		}()], Try(ret1Err[int]()))
	}
	return nil
}
