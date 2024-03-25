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
			a = -1
			fmt.Println(err)
		}
	}()
	Try(ret1Err[bool]())
	a = 42
	return
}

func error_wrapping() (a int, err error) {
	defer handleErrorf(&err, "something wrong")
	Try(ret1Err[bool]())
	a = 42
	return
}

func handleErrorf(err *error, format string, args ...interface{}) {
	if *err != nil {
		*err = fmt.Errorf(format+": %v", append(args, *err)...)
	}
}
