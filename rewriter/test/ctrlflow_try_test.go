//go:build try

package test

import (
	. "github.com/goghcrow/go-try"
)

func swith1() error {
	type (
		A = int
		B = int
		C = int
		D = int
	)

	switch Try(ret1Err[A]()) {
	case 0:
		Try0(ret0Err())
	case 1:
		Try(ret1Err[B]())
	case 2:
		Try2(ret2Err[C, C]())
	case 3:
		Try3(ret3Err[D, D, D]())
	}
	return nil
}

func swith2() error {
	type (
		A = int
		B = int
		C = int
		D = int
		E = int
	)

	switch a, b := Try(ret1Err[A]()), Try(ret1Err[B]()); {
	case a == 1:
		Try0(ret0Err())
	case b == 2:
		Try(ret1Err[C]())
	case a == 3:
		Try2(ret2Err[D, D]())
	case b == 4:
		Try3(ret3Err[E, E, E]())
	}

	return nil
}

func switch3() error {
	switch Try(ret1Err[error]()).(type) {
	}
	return nil
}

func if1() error {
	type (
		A = int
		B = int
		C = int
	)
	if true {
		Try0(ret0Err())
		println(Try(ret1Err[A]()))
		println(Try2(ret2Err[B, string]()))
		println(Try3(ret3Err[C, string, rune]()))
	}
	return nil
}

func for1() error {
	type (
		A = int
		B = int
		C = int
	)

	for {
		Try0(ret0Err())
		println(Try(ret1Err[A]()))
		println(Try2(ret2Err[B, string]()))
		println(Try3(ret3Err[C, string, rune]()))
	}

	return nil
}
