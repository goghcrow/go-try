//go:build try

package test

import (
	. "github.com/goghcrow/go-try"
)

func for_test() error {
	for ; ; id(Try(ret1Err[A]())) {
	}

	for id(Try(ret1Err[A]())) == 42 {
	}

	for i := 1; ; id(Try(func1[int, int](i))) {
	}

	for i := 1; id(Try(func1[int, int](i))) == 42; {
	}

	for i := Try(ret1Err[A]()); Try(func1[int, bool](i)); Try(func1[A, C](i)) {

	}
	return nil
}

func for_try_in_init() error {
	for i := Try(func1[int, int](0)); ; {
		println(i)
	}
}

func for_try_in_cond() error {
	for i := 0; i < Try(func1[int, int](i)); i++ {
		println(i)
	}
	return nil
}

func for_try_in_cond1() error {
	for i := 0; id[bool](false) || Try(func1[int, int](i)) > 1; i++ {
		println(i)
	}
	return nil
}

func for_try_in_post() error {
	for i := 0; i < 42; Try(func1[int, int](i)) {
		println(i)
		i++
	}
	return nil
}

func for_try_in_post1() error {
	for ; ; Try0(ret0()) {
	}
	return nil
}

func for_try_in_post2() error {
	for ; ; Try0(ret0()) {
		println(1)
		continue
	}
	return nil
}

func for_try_in_post21() error {
	for ; ; Try0(ret0()) {
		println(1)
		continue
		println(2)
	}
	return nil
}

func for_try_in_post22() error {
	for ; ; Try0(ret0()) {
		println(1)
		panic(nil)
	}
	return nil
}

func for_try_in_post23() error {
	for ; ; Try0(ret0()) {
		println(1)
		panic(nil)
		println(2)
	}
	return nil
}

func for_try_in_post24() error {
	for ; ; Try0(ret0()) {
		continue
		x := 1
		println(x)
	}
	return nil
}

func for_try_in_post3() error {
L:
	for ; ; Try0(ret0()) {
		continue L
	}
	return nil
}

func for_try_in_post4() error {
L:
	for ; ; Try(func1[int, int](0)) {
		for ; ; Try(func1[int, int](1)) {
			continue L
		}
	}
	return nil
}

func for_try_in_post5() error {
outer:
	for ; ; Try(func1[int, int](1)) {
	inner:
		for ; ; Try(func1[int, int](2)) {
			if true {
				continue inner
			} else {
				continue outer
			}
		}
		continue outer
	}
	return nil
}

func for_try_in_init_cond_post() error {
	for i := Try(func1[int, int](42)); i < Try(func1[int, int](i)); Try(func1[int, int](i)) {
		println(Try(func1[int, int](i)))
	}
	return nil
}

func for_labeled_brk() error {
L:
	for i := Try(func1[int, int](42)); i < Try(func1[int, int](i)); Try(func1[int, int](i)) {
		break L
	}
	return nil
}

func for_labeled_continue() error {
L:
	for i := Try(func1[int, int](42)); i < Try(func1[int, int](i)); Try(func1[int, int](i)) {
		println(i)
		continue L
	}
	return nil
}

func for_labeled_cont_brk_goto() error {
L:
	for i := 0; i < 42; Try(func1[int, int](i)) {
		println(i)
		if i == 42 {
			continue
		}
		if i == 42 {
			continue L
		}
		if i == 42 {
			goto L // skip post
		}
		if i == 42 {
			break
		} else {
			i++
		}
	}
	return nil
}

func for_nested_labeled() error {
	type (
		innerPost = int
		outerPost = int
	)
outer:
	for i := Try(func1[int, int](1)); i < Try(func1[int, int](i+1)); Try(func1[int, outerPost](i + 2)) {
	inner:
		for j := Try(func1[int, int](2)); j < Try(func1[int, int](j+3)); Try(func1[int, innerPost](j + 4)) {
			switch Try(func1[int, int](3)) {
			case Try(func1[int, int](4)):
				continue outer
			case Try(func1[int, int](5)):
				continue inner
			case Try(func1[int, int](6)):
				break outer
			case Try(func1[int, int](7)):
				break inner
			case Try(func1[int, int](8)):
				goto outer
			case Try(func1[int, int](9)):
				goto inner
			default:
				break
			}
		}
	}
	return nil
}
