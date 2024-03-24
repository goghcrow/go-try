//go:build try

package test

import (
	"fmt"

	. "github.com/goghcrow/go-try"
)

func named_error() (a int, err error) {
	Try(ret1Err[bool]())
	a = 42
	return
}

func catch_and_log() (a int, err error) {
	defer func() {
		if err != nil {
			fmt.Println(err)
		}
	}()
	Try(ret1Err[bool]())
	a = 42
	return
}
