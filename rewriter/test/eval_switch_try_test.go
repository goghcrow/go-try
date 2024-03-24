//go:build try

package test

import (
	. "github.com/goghcrow/go-try"
)

func assert(cond bool, msg string) {
	if !cond {
		print("assertion fail: ", msg, "\n")
		panic(1)
	}
}

func testSwitchStmt() (err error) {
	i5 := 5
	i7 := 7
	hello := "hello"

	switch Try(lit(true)) {
	case Try(lit(i5)) < 5:
		assert(false, "<")
	case Try(lit(i5)) == 5:
		assert(true, "!")
	case Try(lit(i5)) > 5:
		assert(false, ">")
	}

	switch {
	case i5 < Try(lit(5)):
		assert(false, "<")
	case i5 == Try(lit(5)):
		assert(true, "!")
	case i5 > Try(lit(5)):
		assert(false, ">")
	}

	switch x := Try(lit(5)); Try(lit(true)) {
	case Try(lit(i5)) < x:
		assert(false, "<")
	case Try(lit(i5)) == x:
		assert(true, "!")
	case Try(lit(i5)) > x:
		assert(false, ">")
	}

	switch x := Try(lit(5)); Try(lit(true)) {
	case i5 < Try(lit(x)):
		assert(false, "<")
	case i5 == Try(lit(x)):
		assert(true, "!")
	case i5 > Try(lit(x)):
		assert(false, ">")
	}

	switch Try(lit(i5)) {
	case Try(lit(0)):
		assert(false, "0")
	case Try(lit(1)):
		assert(false, "1")
	case Try(lit(2)):
		assert(false, "2")
	case Try(lit(3)):
		assert(false, "3")
	case Try(lit(4)):
		assert(false, "4")
	case Try(lit(5)):
		assert(true, "5")
	case Try(lit(6)):
		assert(false, "6")
	case Try(lit(7)):
		assert(false, "7")
	case Try(lit(8)):
		assert(false, "8")
	case Try(lit(9)):
		assert(false, "9")
	default:
		assert(false, "default")
	}

	switch Try(lit(i5)) {
	case Try(lit(0)), 1, Try(lit(2)), Try(lit(3)), 4:
		assert(false, "4")
	case Try(lit(5)):
		assert(true, "5")
	case 6, Try(lit(7)), 8, Try(lit(9)):
		assert(false, "9")
	default:
		assert(false, "default")
	}

	switch Try(lit(i5)) {
	case Try(lit(0)):
	case 1:
	case Try(lit(2)):
	case 3:
	case Try(lit(4)):
		assert(false, "4")
	case 5:
		assert(true, "5")
	case Try(lit(6)):
	case 7:
	case Try(lit(8)):
	case 9:
	default:
		assert(i5 == 5, "good")
	}

	switch Try(lit(i5)) {
	case Try(lit(0)):
		dummy := 0
		_ = dummy
		fallthrough
	case Try(lit(1)):
		dummy := 0
		_ = dummy
		fallthrough
	case Try(lit(2)):
		dummy := 0
		_ = dummy
		fallthrough
	case Try(lit(3)):
		dummy := 0
		_ = dummy
		fallthrough
	case Try(lit(4)):
		dummy := 0
		_ = dummy
		assert(false, "4")
	case Try(lit(5)):
		dummy := 0
		_ = dummy
		fallthrough
	case Try(lit(6)):
		dummy := 0
		_ = dummy
		fallthrough
	case Try(lit(7)):
		dummy := 0
		_ = dummy
		fallthrough
	case Try(lit(8)):
		dummy := 0
		_ = dummy
		fallthrough
	case Try(lit(9)):
		dummy := 0
		_ = dummy
		fallthrough
	default:
		dummy := 0
		_ = dummy
		assert(i5 == 5, "good")
	}

	fired := false
	switch Try(lit(i5)) {
	case Try(lit(0)):
		dummy := 0
		_ = dummy
		fallthrough // tests scoping of cases
	case Try(lit(1)):
		dummy := 0
		_ = dummy
		fallthrough
	case Try(lit(2)):
		dummy := 0
		_ = dummy
		fallthrough
	case Try(lit(3)):
		dummy := 0
		_ = dummy
		fallthrough
	case Try(lit(4)):
		dummy := 0
		_ = dummy
		assert(false, "4")
	case Try(lit(5)):
		dummy := 0
		_ = dummy
		fallthrough
	case Try(lit(6)):
		dummy := 0
		_ = dummy
		fallthrough
	case Try(lit(7)):
		dummy := 0
		_ = dummy
		fallthrough
	case Try(lit(8)):
		dummy := 0
		_ = dummy
		fallthrough
	case Try(lit(9)):
		dummy := 0
		_ = dummy
		fallthrough
	default:
		dummy := 0
		_ = dummy
		fired = !fired
		assert(i5 == 5, "good")
	}
	assert(fired, "fired")

	count := 0
	switch Try(lit(i5)) {
	case Try(lit(0)):
		count = count + 1
		fallthrough
	case Try(lit(1)):
		count = count + 1
		fallthrough
	case Try(lit(2)):
		count = count + 1
		fallthrough
	case Try(lit(3)):
		count = count + 1
		fallthrough
	case Try(lit(4)):
		count = count + 1
		assert(false, "4")
	case Try(lit(5)):
		count = count + 1
		fallthrough
	case Try(lit(6)):
		count = count + 1
		fallthrough
	case Try(lit(7)):
		count = count + 1
		fallthrough
	case Try(lit(8)):
		count = count + 1
		fallthrough
	case Try(lit(9)):
		count = count + 1
		fallthrough
	default:
		assert(i5 == count, "good")
	}
	assert(fired, "fired")

	switch Try(lit(hello)) {
	case "wowie":
		assert(false, "wowie")
	case Try(lit("hello")):
		assert(true, "hello")
	case Try(lit("jumpn")):
		assert(false, "jumpn")
	default:
		assert(false, "default")
	}

	fired = false
	switch i := i5 + Try(lit(2)); Try(lit(i)) {
	case Try(lit(i7)):
		fired = true
	default:
		assert(false, "fail")
	}
	assert(fired, "var")

	// switch on nil-only comparison types
	switch f := Try(lit(func() {})); f {
	case nil:
		assert(false, "f should not be nil")
	default:
	}

	switch m := make(map[int]int); Try(lit(m)) {
	case nil:
		assert(false, "m should not be nil")
	default:
	}

	switch a := make([]int, Try(lit(1))); a {
	case nil:
		assert(false, "m should not be nil")
	default:
	}

	// switch on interface.
	switch i := interface{}("hello"); i {
	case Try(lit(42)):
		assert(false, `i should be "hello"`)
	case Try(lit("hello")):
		assert(true, "hello")
	default:
		assert(false, `i should be "hello"`)
	}

	// switch on implicit bool converted to interface
	// was broken: see issue 3980
	switch i := interface{}(true); {
	case Try(lit(i)):
		assert(true, "true")
	case Try(lit(false)):
		assert(false, "i should be true")
	default:
		assert(false, "i should be true")
	}

	// switch on interface with constant cases differing by type.
	// was rejected by compiler: see issue 4781
	type T int
	type B bool
	type F float64
	type S string
	switch i := interface{}(float64(1.0)); i {
	case nil:
		assert(false, "i should be float64(1.0)")
	case Try(lit((*int)(nil))):
		assert(false, "i should be float64(1.0)")
	case Try(lit(1)):
		assert(false, "i should be float64(1.0)")
	case Try(lit(T(1))):
		assert(false, "i should be float64(1.0)")
	case Try(lit(F(1.0))):
		assert(false, "i should be float64(1.0)")
	case Try(lit(1.0)):
		assert(true, "true")
	case Try(lit("hello")):
		assert(false, "i should be float64(1.0)")
	case Try(lit(S("hello"))):
		assert(false, "i should be float64(1.0)")
	case true, Try(lit(B(false))):
		assert(false, "i should be float64(1.0)")
	case false, Try(lit(B(true))):
		assert(false, "i should be float64(1.0)")
	}

	// switch on array.
	switch ar := [3]int{1, 2, Try(lit(3))}; ar {
	case Try(lit([3]int{1, 2, Try(lit(3))})):
		assert(true, "[1 2 3]")
	case Try(lit([3]int{4, 5, Try(lit(6))})):
		assert(false, "ar should be [1 2 3]")
	default:
		assert(false, "ar should be [1 2 3]")
	}

	// switch on channel
	switch c1, c2 := make(chan int), Try(lit(make(chan int))); c1 {
	case nil:
		assert(false, "c1 did not match itself")
	case Try(lit(c2)):
		assert(false, "c1 did not match itself")
	case Try(lit(c1)):
		assert(true, "chan")
	default:
		assert(false, "c1 did not match itself")
	}

	// empty switch
	switch {
	}

	// empty switch with default case.
	fired = false
	switch {
	default:
		fired = Try(lit(true))
	}
	assert(fired, "fail")

	// Default and fallthrough.
	count = 0
	switch {
	default:
		count++
		fallthrough
	case Try(lit(false)):
		count++
	}
	assert(count == 2, "fail")

	// fallthrough to default, which is not at end.
	count = 0
	switch Try(lit(i5)) {
	case Try(lit(5)):
		count++
		fallthrough
	default:
		count++
	case Try(lit(6)):
		count++
	}
	assert(count == 2, "fail")

	i := 0
	fired = false
	switch x := 5; {
	case i < x:
		fired = true
	case i == Try(lit(x)):
		assert(false, "fail")
	case Try(lit(i > x)):
		assert(false, "fail")
	}
	assert(fired, "fail")

	// Unified IR converts the tag and all case values to empty
	// interface, when any of the case values aren't assignable to the
	// tag value's type. Make sure that `case nil:` compares against the
	// tag type's nil value (i.e., `(*int)(nil)`), not nil interface
	// (i.e., `any(nil)`).
	switch (*int)(nil) {
	case nil:
		// ok
	case any(nil):
		assert(false, "case any(nil) matched")
	default:
		assert(false, "default matched")
	}

	return
}
