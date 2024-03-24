//go:build try

package test

import (
	. "github.com/goghcrow/go-try"
)

func try_underlying_fun() error {
	{
		type F func() int
		var f F
		_ = f() + Try(ret1Err[int]())
	}

	{
		type ErrF func() (int, error)
		var errF ErrF
		_ = Try(errF())
	}

	return nil
}
