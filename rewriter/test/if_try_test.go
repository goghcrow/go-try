//go:build try

package test

import (
	. "github.com/goghcrow/go-try"
)

func try_in_if_init_or_if_cond() error {
	if Try(func1[int, bool](1)) {

	} else if false {

	} else if a := Try(func1[int, bool](2)); a {

	} else if Try(func1[int, bool](3)) {

	} else if true {

	}

	return nil
}
