//go:build try

package test

import (
	. "github.com/goghcrow/go-try"
)

func dbg() (err error) {
	//for a := 1; a > 0; a = Try(ret1Err[int]()) {
	//	a := 42
	//	_ = a
	//}

	//for a := 1; a > 0; println(Try(func1[int, int](a))) {
	//	a := 42
	//	_ = a
	//}

	//for a := 1; a > 0; a = Try(ret1Err[int]()) {
	//	a := 42
	//	_ = a
	//	if true {
	//
	//		continue
	//	}
	//	println(a)
	//}

	for i := Try(ret1Err[A]()); Try(func1[int, bool](i)); Try(func1[A, C](i)) {
		println(i)
	}

	return nil
}
