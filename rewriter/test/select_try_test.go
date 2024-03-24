//go:build try

package test

import (
	. "github.com/goghcrow/go-try"
)

func select_eval_order0() {
	getRecvCh := func(s string, buf int) <-chan int {
		println("getRecvCh called: " + s)
		ch := make(chan int, buf)
		if buf > 0 {
			ch <- 42
		}
		return ch
	}
	mkValToSend := func() int {
		println("mkValToSend called")
		return 0
	}
	mkSndCh := func() chan<- int {
		println("mkSndCh called")
		return make(chan int)
	}
	f := func() int {
		println("f called")
		return 0
	}
	getvar := func(p *int) *int {
		println("getvar called")
		return p
	}

	var a = []int{1, 2, 3}
	var i1, i2 int

	select {
	case *getvar(&i1) = <-getRecvCh("c1", 1):
		print("received ", i1, " from c1\n")
	case mkSndCh() <- mkValToSend():
		print("sent ", i2, " to c2\n")
	case a[f()] = <-getRecvCh("c2", 0):
		// same as:
		// case t := <-c4
		//	a[f()] = t
	default:
		print("no communication\n")
	}

	// output
	// getRecvCh called: c1
	// mkSndCh called
	// mkValToSend called
	// getRecvCh called: c2
	// getvar called
	// received 42 from c1
}

func select_eval_order() error {
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
	select {
	case <-Try(ret1Err[chan A]()):
	case *Try(ret1Err[*B]()), *Try(ret1Err[*bool]()) = <-Try(ret1Err[chan C]()):
	case Try(ret1Err[chan D]()) <- Try(ret1Err[E]()):
	case Try(ret1Err[[]F]())[Try(ret1Err[G]())] = <-Try(ret1Err[chan H]()):
	default:
	}
	return nil
}

func select_eval_order1() error {
	type (
		A = int
		B = int
		C = int
	)
	var ch <-chan int
	select {
	case <-ch:
	case *Try(ret1Err[*A]()), _ = <-ch:
	case *Try(ret1Err[*B]()) = <-ch:
	case *Try(ret1Err[*C]()), *Try(ret1Err[*bool]()) = <-ch:
	case _, *Try(ret1Err[*bool]()) = <-ch:
	default:
	}
	return nil
}

func rewrite_select() error {
	{
		select {
		case x, ok := (<-Try(ret1Err[chan int]())):
			_, _ = x, ok
		}
	}

	{
		select {
		case *Try(ret1Err[*int]()) = <-(<-chan int)(nil):
		default:
		}
	}

	return nil
}
