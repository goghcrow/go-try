//go:build try

package test

import (
	. "github.com/goghcrow/go-try"
)

type myErr struct {
}

func (e myErr) Error() string { return "" }

func try_ret_myerr() (err error) {
	_ = 1 + Try(func() (int, error) {
		return 42, myErr{}
	}())
	return nil
}
