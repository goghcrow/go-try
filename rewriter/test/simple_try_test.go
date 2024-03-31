//go:build try

package test

import (
	"errors"

	. "github.com/goghcrow/go-try"
)

type (
	Int int
	Str string
)

var helloErr = errors.New("hello error!")

func answer() (int, error) {
	return Try(ret1Err[int]()), nil
}

func assign() (_ error) {
	x, y := Try(ret1Err[int]()), 42
	consume2(x, y)

	v1, v2 := Try(42, helloErr), 42
	consume2(v1, v2)
	return nil
}

func binary() (_ error) {
	consume1(Try(ret1Err[int]()) + 1)
	return nil
}

func ret0() (_ error) {
	Try0(ret0Err())

	Try0(helloErr)
	return
}

func ret1() (_ Int, _ error) {
	Try(ret1Err[int]())
	consume1(Try(ret1Err[int]()) + 1)

	Try(42, helloErr)
	consume1(Try(42, helloErr) + 1)

	return
}

func ret2() (_ Int, _ Str, _ error) {
	Try2(ret2Err[int, string]())
	iV, bV := Try2(ret2Err[int, string]())
	consume2(iV, bV)
	consume2(Try2(ret2Err[int, string]()))

	Try2(42, "answer", helloErr)
	iV, bV = Try2(42, "answer", helloErr)
	consume2(iV, bV)
	consume2(Try2(42, "answer", helloErr))
	return
}

func ret2_grouped_ret() (_, _ Int, _ error) {
	Try2(ret2Err[int, byte]())
	return
}

func ret3() (_ *Int, _ error) {
	Try3(ret3Err[int, rune, string]())
	iV, bV, sV := Try3(ret3Err[int, rune, string]())
	consume3(iV, bV, sV)
	consume3(Try3(ret3Err[int, rune, string]()))

	Try3(42, 'a', "hello", helloErr)
	iV, bV, sV = Try3(42, 'a', "hello", helloErr)
	consume3(iV, bV, sV)
	consume3(Try3(42, 'a', "hello", helloErr))

	func(int, rune, string) {}(Try3(ret3Err[int, rune, string]()))
	return
}

func funcLit() {
	go func() {
		_ = func() error {
			Try0(ret0())
			consume1(Try(ret1Err[int]()))
			return nil
		}()
	}()

	defer func() {
		_ = func() error {
			Try0(ret0())
			consume1(Try(ret1Err[int]()))
			return nil
		}()
	}()

	if func() int {
		func() (int, error) {
			return id(Try(ret1Err[int]())), nil
		}()
		return 42
	}() == 42 {
	}
}

func fnlit() {
	go func() {
		_ = func() error {
			Try0(ret0())
			consume1(Try(ret1Err[int]()))
			return nil
		}()
	}()

	func() {
		_ = func() error {
			Try0(ret0())
			consume1(Try(ret1Err[int]()))
			return nil
		}()
	}()
}

func unparen() error {
	{
		_ = Try((ret1Err[int]()))
	}
	return nil
}
