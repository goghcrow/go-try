//go:build try

package test

import (
	. "github.com/goghcrow/go-try"
)

func goto_label() error {
L:
	var a = Try(ret1Err[int]())
	goto L
	println(a)
	return nil
}

func goto_label1() error {
	if true {

	} else {
	L:
		var a = Try(ret1Err[int]())
		goto L
		println(a)
	}
	return nil
}

func range_with_break() error {
L:
	for range Try(ret1Err[[]int]()) {
		break L
	}
	return nil
}

func range_with_break1() error {
Outer:
	for range Try(ret1Err[[]int]()) {
	Inner:
		for range Try(ret1Err[[]string]()) {
			if true {
				break Inner
			}
			break Outer
		}
	}
	return nil
}

func range_with_goto() error {
L:
	for range Try(ret1Err[[]int]()) {
		goto L
	}
	return nil
}

func range_with_goto1() error {
Outer:
	for range Try(ret1Err[[]int]()) {
	Inner:
		for range Try(ret1Err[[]string]()) {
			if true {
				goto Inner
			}
			goto Outer
		}
	}
	return nil
}

func range_with_break_goto(a int) error {
Outer:
	for range Try(ret1Err[[]int]()) {
	Inner:
		for range Try(ret1Err[[]string]()) {
			switch a {
			case 1:
				goto Inner
			case 2:
				goto Outer
			case 3:
				break Inner
			case 4:
				break Outer
			case 5:
				continue Inner
			case 6:
				continue Outer
			}
		}
	}
	return nil
}

func for_with_break() error {
L:
	for i, xs := Try(ret1Err[int]()), Try(ret1Err[[]int]()); i < 42; i++ {
		if xs[i] == 0 {
			break L
		}
	}
	return nil
}

func for_with_break1() error {
Outer:
	for i, xs := Try(ret1Err[int]()), Try(ret1Err[[]int]()); i < 42; i++ {
	Inner:
		for j, ys := Try(ret1Err[int]()), Try(ret1Err[[]int]()); j < 42; j++ {
			if xs[i] == 0 {
				break Inner
			} else if ys[j] == 0 {
				break Outer
			}
		}
	}
	return nil
}

func for_with_goto() error {
L:
	for i, xs := Try(ret1Err[int]()), Try(ret1Err[[]int]()); i < 42; i++ {
		if xs[i] == 0 {
			goto L
		}
	}
	return nil
}

func for_with_goto1() error {
Outer:
	for i, xs := Try(ret1Err[int]()), Try(ret1Err[[]int]()); i < 42; i++ {
	Inner:
		for j, ys := Try(ret1Err[int]()), Try(ret1Err[[]int]()); j < 42; j++ {
			if xs[i] == 0 {
				goto Inner
			} else if ys[j] == 0 {
				goto Outer
			}
		}
	}
	return nil
}

func for_with_break_goto(a int) error {
Outer:
	for i, _ := Try(ret1Err[int]()), Try(ret1Err[[]int]()); i < 42; i++ {
	Inner:
		for j, _ := Try(ret1Err[int]()), Try(ret1Err[[]int]()); j < 42; j++ {
			switch a {
			case 1:
				goto Inner
			case 2:
				goto Outer
			case 3:
				break Inner
			case 4:
				break Outer
			case 5:
				continue Inner
			case 6:
				continue Outer
			}
		}
	}
	return nil
}

func switch_with_break() error {
L:
	switch Try(ret1Err[int]()) {
	case 42:
		break L
	}
	return nil
}

func switch_with_break1() error {
Outer:
	switch Try(ret1Err[int]()) {
	case 42:
	Inner:
		switch Try(ret1Err[int]()) {
		case 42:
			break Outer
		case 100:
			break Inner
		}

	}
	return nil
}

func switch_with_goto() error {
L:
	switch Try(ret1Err[int]()) {
	case 42:
		goto L
	}
	return nil
}

func switch_with_goto1() error {
Outer:
	switch Try(ret1Err[int]()) {
	case 42:
	Inner:
		switch Try(ret1Err[int]()) {
		case 42:
			goto Outer
		case 100:
			goto Inner
		}
	}
	return nil
}

func switch_with_break_goto(a int) error {
Outer:
	switch Try(ret1Err[int]()) {
	case 42:
	Inner:
		switch Try(ret1Err[int]()) {
		case 1:
			goto Inner
		case 2:
			goto Outer
		case 3:
			break Inner
		case 4:
			break Outer
		}
	}
	return nil
}

func select_with_break() error {
L:
	select {
	case <-Try(ret1Err[chan int]()):
		break L
	}
	return nil
}

func select_with_break1() error {
Outer:
	select {
	case <-Try(ret1Err[chan int]()):
	Inner:
		select {
		case <-Try(ret1Err[chan int]()):
			break Outer
		case <-Try(ret1Err[chan int]()):
			break Inner
		}
	}
	return nil
}

func select_with_goto() error {
L:
	select {
	case <-Try(ret1Err[chan int]()):
		goto L
	}
	return nil
}

func select_with_goto1() error {
Outer:
	select {
	case <-Try(ret1Err[chan int]()):
	Inner:
		select {
		case <-Try(ret1Err[chan int]()):
			goto Outer
		case <-Try(ret1Err[chan int]()):
			goto Inner
		}
	}
	return nil
}

func select_with_break_goto(a int) error {
Outer:
	select {
	case <-Try(ret1Err[chan int]()):
	Inner:
		select {
		case <-Try(ret1Err[chan int]()):
			goto Outer
		case <-Try(ret1Err[chan int]()):
			goto Inner
		case <-Try(ret1Err[chan int]()):
			break Outer
		case <-Try(ret1Err[chan int]()):
			break Inner
		}
	}
	return nil
}
