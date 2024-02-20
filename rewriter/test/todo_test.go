//go:build try

package test

import (
	. "github.com/goghcrow/go-try"
)

// var n int = try.Try(ret1Err[int]())

func todo1() (_ error) {
	type (
		fst = int
		snd = int
	)
	println(Try(ret1Err[fst]()) + Try(ret1Err[snd]()))

	// todo 应该先执行 return1 再执行 ret1Err
	consume2(return1[fst](), Try(ret1Err[snd]()))

	// todo 应该先执行 return1 再执行 ret1Err
	println(return1[int]() + Try(ret1Err[int]()))

	return
}

func todo2() error {
	type (
		A = int
		B = int
		C = int
		D = int
		E = int
		F = int
		G = int
		H = int
	)

	// if id(Try(ret1Err[A]())) == 42 {
	// } else if Try(ret1Err[B]()) == 42 {
	// }

	// if true {
	// } else if Try(ret1Err[B]()) == 42 {
	// }

	// switch {
	// case Try(ret1Err[A]()) == 42:
	// case id(Try(ret1Err[B]())) == 42:
	// }

	// switch {
	// case true:
	// case id(Try(ret1Err[B]())) == 42:
	// }

	// switch i := 0; Try(func1[int, A](i)) {
	// }

	switch Try(ret1Err[error]()).(type) {
	}

	// switch n := 1; Try(func1[int, fmt.Stringer](n)).(type) {
	//
	// }

	// for ; ; id(Try(ret1Err[A]())) {
	// }

	// for id(Try(ret1Err[A]())) == 42 {
	// }

	// for i := 1; ; id(Try(func1[int, int](i))) {
	// }

	// for i := 1; id(Try(func1[int, int](i))) == 42; {
	// }

	// for i := Try(ret1Err[A]()); Try(ret1Err[bool]()); Try(func1[A, C](i)) {
	//
	// }

	// go try.Try0(ret0())
	// defer try.Try0(ret0())

	// go consume1(try.Try(ret1Err[int]()))
	// defer consume1(try.Try(ret1Err[int]()))

	return nil
}
