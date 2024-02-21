//go:build try

package test

import (
	. "github.com/goghcrow/go-try"
)

func emptyStmt() error {
	switch Try(ret1Err[int]()) {
	case 0:
		Try0(ret0Err())
	}
	return nil
}
