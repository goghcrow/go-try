//go:build try

package test

import (
	. "github.com/goghcrow/go-try"
)

type (
	A = int
	B = int
	C = int
	D = int
	E = int
)

func switch_underlying_nil() error {
	type (
		Nil any
	)
	switch {
	case Try(ret1Err[any]()):
	case Try(ret1Err[Nil]()):
	}
	return nil
}

func switch_fallthrough_copy(i int) (err error) {
	type (
		A = int
		B = int
		C = int
		D = int
		E = int
	)
	switch i {
	case Try(ret1Err[A]()):
		fallthrough
	case Try(ret1Err[C]()):
		fallthrough
	case Try(ret1Err[D]()):
		Try(ret1Err[E]())
	}
	return nil
}

func switch_fallthrough_copy1(i int) (err error) {
	type (
		A = int
		B = int
		C = int
		D = int
		E = int
	)
	switch i {
	case Try(ret1Err[A]()):
		fallthrough
	case Try(ret1Err[C]()):
		// fallthrough
	case Try(ret1Err[D]()):
		Try(ret1Err[E]())
	}
	return nil
}

func switch_fallthrough_copy2(i int) (err error) {
	type (
		A = int
		B = int
		C = int
		D = int
		E = int
		F = int
	)
	switch i {
	case Try(ret1Err[A]()):
		Try(ret1Err[B]())
		fallthrough
	case Try(ret1Err[C]()):
		Try(ret1Err[D]())
		fallthrough
	case Try(ret1Err[E]()):
		Try(ret1Err[F]())
	}

	return nil
}

func switch_fallthrough() (err error) {
	a := 1
	switch a {
	case Try(ret1Err[A]()):
		goto L
		println("1")
	L:
		fallthrough
	default:
		println("default")
		fallthrough
	case Try(ret1Err[B]()):
		println("2")
	}
	return
}

func switch_scope_shadow() error {
	var x int
	switch {
	case Try(ret1Err[int]()) == 1:
		x := 1
		println(x + Try(ret1Err[int]()))

		fallthrough
	default:
		x = 2
	}
	println(x)
	return nil
}

func switch_scope_conflict() error {
	var x int
	switch {
	case Try(ret1Err[int]()) == 1:
		x := 1
		println(x + Try(ret1Err[int]()))

		fallthrough
	default:
		x := 2
		println(x)
	}
	println(x)
	return nil
}

func switch_try_in_init() error {
	switch Try(func1[int, A](0)); {
	default:
	}

	return nil
}

func switch_try_in_tag() error {
	switch Try(func1[int, A](0)) {
	default:
	}

	return nil
}

func switch_case_use_init() error {
	switch i := Try(func1[int, A](0)); {
	case Try(func1[int, int](i)) == 42:
		println("hello")
	default:
	}

	return nil
}

func switch_case_use_init1() error {
	switch i := Try(func1[int, A](0)); {
	case Try(func1[int, int](i)) == 42 || Try(func1[int, int](i+1)) == 100:
		println("hello")
	default:
	}

	return nil
}

func switch_cond_use_init() error {
	switch i := Try(func1[int, A](0)); Try(func1[int, B](i)) {
	case Try(func1[int, C](i)):
		println("hello")
	default:
	}

	return nil
}

func switch_cond_use_init1() error {
	switch i := 42; Try(func1[int, B](i)) {
	default:
	}
	return nil
}

func swith_try_in_case_no_tag() error {
	switch i := Try(func1[int, A](0)); {
	case Try(func1[int, B](i)) == 42:
		println("B")
	case Try(func1[int, C](i)) == 42:
		println("C")
	default:
		println("D")
	}

	return nil
}

func swith_mixed_cases_no_tag() error {
	switch i := Try(func1[int, A](0)); {
	case Try(func1[int, B](i)) == 42:
		println("B")
	case id[int](i) == 42:
		println("C")
	case Try(func1[int, D](i)) == 42:
		println("D1")
	case id[int](i) == 42:
		println("E")
	case Try(func1[int, D](i)) == 42:
		println("D2")
	default:
		println("default")
	}

	return nil
}

func swith_mixed_cases() error {
	switch i := Try(func1[int, A](0)); Try(func1[int, A](i)) {
	case Try(func1[int, B](i)):
		println("B")
	case id[int](i):
		println("C")
	case Try(func1[int, D](i)):
		println("D1")
	case id[int](i):
		println("E")
	case Try(func1[int, D](i)):
		println("D2")
	default:
		println("default")
	}

	return nil
}

func switch_try_in_case() error {
	switch i := Try(func1[int, A](0)); Try(func1[int, A](i)) {
	case Try(func1[int, B](i)):
		println("B")
	case Try(func1[int, C](i)):
		println("C")
	default:
		println("D")
	}

	return nil
}

func switch_multi_case() error {
	switch i := Try(func1[int, A](0)); Try(func1[int, B](i)) {
	case Try(func1[int, C](i)), Try(func1[int, D](i)):
		println("hello")
	case Try(func1[int, E](i)):
		println("hello")
	default:
	}

	return nil
}

func switch_multi_case2() error {
	switch {
	case Try(func1[int, A](1)) == 1, Try(func1[int, A](1)) == 2:
		println("1,2")
	case Try(func1[int, A](1)) == 3:
		println("3")
	}

	return nil
}

func switch_break() error {
	switch {
	case Try(func1[int, A](1)) == 42:
		println("hello")
		break
	}

	return nil
}

func switch_nested_break() error {
	switch {
	case Try(func1[int, A](1)) == 42:
		println("hello")
		switch {
		case Try(func1[int, A](1)) == 42:
			println("hello")
			break
		}
		break
	}

	return nil
}

func switch_labeled_break() error {
L:
	switch {
	case Try(func1[int, A](1)) == 42:
		println("hello")
		break L
	}

	return nil
}

func switch_goto() error {
L:
	switch {
	case Try(func1[int, A](1)) == 42:
		goto L
	}

	return nil
}

func switch_labeled_break_and_goto() error {
L:
	switch {
	case Try(func1[int, A](1)) == 42:
		break L
	case Try(func1[int, B](1)) == 42:
		goto L
	}

	return nil
}

func switch_nested_goto() error {
L:
	switch {
	case Try(func1[int, A](1)) == 42:
		switch {
		case Try(func1[int, B](1)) == 42:
			for {
				goto L
			}
		}
	}

	return nil
}

func switch_nested_labeled_break() error {
L:
	switch {
	case Try(func1[int, A](1)) == 42:
		println("outer")
		switch {
		case Try(func1[int, B](1)) == 42:
			println("inner")
			for {
				break L
			}
		}
	}

	return nil
}

func switch_nested_labeled_break_and_goto() error {
L:
	switch {
	case Try(func1[int, A](1)) == 42:
		println("outer")
		switch {
		case Try(func1[int, B](1)) == 42:
			println("inner")
			for {
				if true {
					break L
				} else {
					goto L
				}
			}
		}
	}

	return nil
}

func switch_nested() error {
outer:
	switch {
	case Try(func1[int, A](1)) == 42:
		println("outer")
	inner:
		switch {
		case Try(func1[int, B](1)) == 42:
			break inner
		case Try(func1[int, C](1)) == 42:
			goto inner
		case Try(func1[int, D](1)) == 42:
			println("inner")
			break outer
		case Try(func1[int, E](1)) == 42:
			println("inner")
			goto outer
		}
	default:
		println("default")
	}

	return nil
}

func switch_fallthrough_default() error {
	switch {
	default:
		println("fallthrough")
		fallthrough
	case Try(func1[int, int](1)) == 42:
		println("hello")
	}

	return nil
}

func switch_labeled_fallthrough() error {
	switch {
	case Try(func1[int, int](1)) == 42:
		println("hello")
		if false {
			goto L
		}
	L:
		fallthrough
	default:
		println("default")
	}

	return nil
}

func switch_mixed() error {
L:
	switch i := Try(func1[int, A](0)); {
	case Try(func1[int, B](i)) == 42:
		switch {
		default:
			for i := 0; i < 1; i++ {
				goto labeldFall
			}
		}
		println("B")
	labeldFall:
		fallthrough
	default:
		println("default")
		fallthrough
	case id[int](i) == 42:
		println("42-a")
	case Try(func1[int, C](i)) == 42:
		println("C")
		break L
	case Try(func1[int, C](i)) == 42:
		for i := 0; i < 10; i++ {
			switch Try(func1[int, error](0)).(type) {
			case error:
				println("C2")
				break L
			case nil:
				for {
					goto L
				}
			}
		}
	case id[int](i) == 42:
		println("42-b")
		goto L
	}
	goto L
}
