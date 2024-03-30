//go:build !try

// Code generated by github.com/goghcrow/go-try DO NOT EDIT.
package test

type (
	A = int
	B = int
	C = int
	D = int
	E = int
)

func switch_underlying_nil() error {
	type Nil any
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[any]()
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		if (42 == 42) == 𝘃𝗮𝗹𝟭 {
		} else {
			𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := ret1Err[Nil]()
			if 𝗲𝗿𝗿𝟮 != nil {
				return 𝗲𝗿𝗿𝟮
			}
			if (42 == 42) == 𝘃𝗮𝗹𝟮 {
			}
		}
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
	{
		𝘃𝗮𝗹𝟭 := i
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟭 := ret1Err[A]()
		if 𝗲𝗿𝗿𝟭 != nil {
			err = 𝗲𝗿𝗿𝟭
			return
		}
		if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟮 {
			_, 𝗲𝗿𝗿𝟮 := ret1Err[E]()
			if 𝗲𝗿𝗿𝟮 != nil {
				err = 𝗲𝗿𝗿𝟮
				return
			}
		} else {
			𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟯 := ret1Err[C]()
			if 𝗲𝗿𝗿𝟯 != nil {
				err = 𝗲𝗿𝗿𝟯
				return
			}
			if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟯 {
				_, 𝗲𝗿𝗿𝟰 := ret1Err[E]()
				if 𝗲𝗿𝗿𝟰 != nil {
					err = 𝗲𝗿𝗿𝟰
					return
				}
			} else {
				𝘃𝗮𝗹𝟰, 𝗲𝗿𝗿𝟱 := ret1Err[D]()
				if 𝗲𝗿𝗿𝟱 != nil {
					err = 𝗲𝗿𝗿𝟱
					return
				}
				if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟰 {
					_, 𝗲𝗿𝗿𝟲 := ret1Err[E]()
					if 𝗲𝗿𝗿𝟲 != nil {
						err = 𝗲𝗿𝗿𝟲
						return
					}
				}
			}
		}
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
	{
		𝘃𝗮𝗹𝟭 := i
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟭 := ret1Err[A]()
		if 𝗲𝗿𝗿𝟭 != nil {
			err = 𝗲𝗿𝗿𝟭
			return
		}
		if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟮 {
		} else {
			𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟮 := ret1Err[C]()
			if 𝗲𝗿𝗿𝟮 != nil {
				err = 𝗲𝗿𝗿𝟮
				return
			}
			if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟯 {
			} else {
				𝘃𝗮𝗹𝟰, 𝗲𝗿𝗿𝟯 := ret1Err[D]()
				if 𝗲𝗿𝗿𝟯 != nil {
					err = 𝗲𝗿𝗿𝟯
					return
				}
				if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟰 {
					_, 𝗲𝗿𝗿𝟰 := ret1Err[E]()
					if 𝗲𝗿𝗿𝟰 != nil {
						err = 𝗲𝗿𝗿𝟰
						return
					}
				}
			}
		}
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
	{
		𝘃𝗮𝗹𝟭 := i
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟭 := ret1Err[A]()
		if 𝗲𝗿𝗿𝟭 != nil {
			err = 𝗲𝗿𝗿𝟭
			return
		}
		if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟮 {
			{
				_, 𝗲𝗿𝗿𝟮 := ret1Err[B]()
				if 𝗲𝗿𝗿𝟮 != nil {
					err = 𝗲𝗿𝗿𝟮
					return
				}
			}
			{
				_, 𝗲𝗿𝗿𝟯 := ret1Err[D]()
				if 𝗲𝗿𝗿𝟯 != nil {
					err = 𝗲𝗿𝗿𝟯
					return
				}
			}
			{
				_, 𝗲𝗿𝗿𝟰 := ret1Err[F]()
				if 𝗲𝗿𝗿𝟰 != nil {
					err = 𝗲𝗿𝗿𝟰
					return
				}
			}
		} else {
			𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟱 := ret1Err[C]()
			if 𝗲𝗿𝗿𝟱 != nil {
				err = 𝗲𝗿𝗿𝟱
				return
			}
			if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟯 {
				{
					_, 𝗲𝗿𝗿𝟲 := ret1Err[D]()
					if 𝗲𝗿𝗿𝟲 != nil {
						err = 𝗲𝗿𝗿𝟲
						return
					}
				}
				{
					_, 𝗲𝗿𝗿𝟳 := ret1Err[F]()
					if 𝗲𝗿𝗿𝟳 != nil {
						err = 𝗲𝗿𝗿𝟳
						return
					}
				}
			} else {
				𝘃𝗮𝗹𝟰, 𝗲𝗿𝗿𝟴 := ret1Err[E]()
				if 𝗲𝗿𝗿𝟴 != nil {
					err = 𝗲𝗿𝗿𝟴
					return
				}
				if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟰 {
					_, 𝗲𝗿𝗿𝟵 := ret1Err[F]()
					if 𝗲𝗿𝗿𝟵 != nil {
						err = 𝗲𝗿𝗿𝟵
						return
					}
				}
			}
		}
	}
	return nil
}
func switch_fallthrough() (err error) {
	a := 1
	{
		𝘃𝗮𝗹𝟭 := a
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟭 := ret1Err[A]()
		if 𝗲𝗿𝗿𝟭 != nil {
			err = 𝗲𝗿𝗿𝟭
			return
		}
		if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟮 {
			{
				goto L
				println("1")
			L:
			}
			{
				println("default")
			}
			{
				println("2")
			}
		} else {
			𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟮 := ret1Err[B]()
			if 𝗲𝗿𝗿𝟮 != nil {
				err = 𝗲𝗿𝗿𝟮
				return
			}
			if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟯 {
				println("2")
			} else {
				{
					println("default")
				}
				{
					println("2")
				}
			}
		}
	}
	return
}
func switch_scope_shadow() error {
	var x int
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[int]()
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		if 𝘃𝗮𝗹𝟭 == 1 {
			{
				x := 1
				𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := ret1Err[int]()
				if 𝗲𝗿𝗿𝟮 != nil {
					return 𝗲𝗿𝗿𝟮
				}
				println(x + 𝘃𝗮𝗹𝟮)
			}
			{
				x = 2
			}
		} else {
			x = 2
		}
	}
	println(x)
	return nil
}
func switch_scope_conflict() error {
	var x int
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[int]()
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		if 𝘃𝗮𝗹𝟭 == 1 {
			{
				x := 1
				𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := ret1Err[int]()
				if 𝗲𝗿𝗿𝟮 != nil {
					return 𝗲𝗿𝗿𝟮
				}
				println(x + 𝘃𝗮𝗹𝟮)
			}
			{
				x := 2
				println(x)
			}
		} else {
			x := 2
			println(x)
		}
	}
	println(x)
	return nil
}
func switch_try_in_init() error {
	_, 𝗲𝗿𝗿𝟭 := func1[int, A](0)
	if 𝗲𝗿𝗿𝟭 != nil {
		return 𝗲𝗿𝗿𝟭
	}
	switch {
	default:
	}
	return nil
}
func switch_try_in_tag() error {
	𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, A](0)
	if 𝗲𝗿𝗿𝟭 != nil {
		return 𝗲𝗿𝗿𝟭
	}
	switch 𝘃𝗮𝗹𝟭 {
	default:
	}
	return nil
}
func switch_case_use_init() error {
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, A](0)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		i := 𝘃𝗮𝗹𝟭
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := func1[int, int](i)
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
		if 𝘃𝗮𝗹𝟮 == 42 {
			println("hello")
		} else {
		}
	}
	return nil
}
func switch_case_use_init1() error {
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, A](0)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		i := 𝘃𝗮𝗹𝟭
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := func1[int, int](i)
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
		𝘃𝗮𝗹𝟰 := 𝘃𝗮𝗹𝟮 == 42
		if !𝘃𝗮𝗹𝟰 {
			𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟯 := func1[int, int](i + 1)
			if 𝗲𝗿𝗿𝟯 != nil {
				return 𝗲𝗿𝗿𝟯
			}
			𝘃𝗮𝗹𝟰 = 𝘃𝗮𝗹𝟯 == 100
		}
		if 𝘃𝗮𝗹𝟰 {
			println("hello")
		} else {
		}
	}
	return nil
}
func switch_cond_use_init() error {
	{
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟭 := func1[int, A](0)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		i := 𝘃𝗮𝗹𝟮
		𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟮 := func1[int, B](i)
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
		𝘃𝗮𝗹𝟭 := 𝘃𝗮𝗹𝟯
		𝘃𝗮𝗹𝟰, 𝗲𝗿𝗿𝟯 := func1[int, C](i)
		if 𝗲𝗿𝗿𝟯 != nil {
			return 𝗲𝗿𝗿𝟯
		}
		if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟰 {
			println("hello")
		} else {
		}
	}
	return nil
}
func switch_cond_use_init1() error {
	{
		i := 42
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, B](i)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		switch 𝘃𝗮𝗹𝟭 {
		default:
		}
	}
	return nil
}
func swith_try_in_case_no_tag() error {
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, A](0)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		i := 𝘃𝗮𝗹𝟭
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := func1[int, B](i)
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
		if 𝘃𝗮𝗹𝟮 == 42 {
			println("B")
		} else {
			𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟯 := func1[int, C](i)
			if 𝗲𝗿𝗿𝟯 != nil {
				return 𝗲𝗿𝗿𝟯
			}
			if 𝘃𝗮𝗹𝟯 == 42 {
				println("C")
			} else {
				println("D")
			}
		}
	}
	return nil
}
func swith_mixed_cases_no_tag() error {
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, A](0)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		i := 𝘃𝗮𝗹𝟭
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := func1[int, B](i)
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
		if 𝘃𝗮𝗹𝟮 == 42 {
			println("B")
		} else if id[int](i) == 42 {
			println("C")
		} else {
			𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟯 := func1[int, D](i)
			if 𝗲𝗿𝗿𝟯 != nil {
				return 𝗲𝗿𝗿𝟯
			}
			if 𝘃𝗮𝗹𝟯 == 42 {
				println("D1")
			} else if id[int](i) == 42 {
				println("E")
			} else {
				𝘃𝗮𝗹𝟰, 𝗲𝗿𝗿𝟰 := func1[int, D](i)
				if 𝗲𝗿𝗿𝟰 != nil {
					return 𝗲𝗿𝗿𝟰
				}
				if 𝘃𝗮𝗹𝟰 == 42 {
					println("D2")
				} else {
					println("default")
				}
			}
		}
	}
	return nil
}
func swith_mixed_cases() error {
	{
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟭 := func1[int, A](0)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		i := 𝘃𝗮𝗹𝟮
		𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟮 := func1[int, A](i)
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
		𝘃𝗮𝗹𝟭 := 𝘃𝗮𝗹𝟯
		𝘃𝗮𝗹𝟰, 𝗲𝗿𝗿𝟯 := func1[int, B](i)
		if 𝗲𝗿𝗿𝟯 != nil {
			return 𝗲𝗿𝗿𝟯
		}
		if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟰 {
			println("B")
		} else if 𝘃𝗮𝗹𝟭 == id[int](i) {
			println("C")
		} else {
			𝘃𝗮𝗹𝟱, 𝗲𝗿𝗿𝟰 := func1[int, D](i)
			if 𝗲𝗿𝗿𝟰 != nil {
				return 𝗲𝗿𝗿𝟰
			}
			if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟱 {
				println("D1")
			} else if 𝘃𝗮𝗹𝟭 == id[int](i) {
				println("E")
			} else {
				𝘃𝗮𝗹𝟲, 𝗲𝗿𝗿𝟱 := func1[int, D](i)
				if 𝗲𝗿𝗿𝟱 != nil {
					return 𝗲𝗿𝗿𝟱
				}
				if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟲 {
					println("D2")
				} else {
					println("default")
				}
			}
		}
	}
	return nil
}
func switch_try_in_case() error {
	{
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟭 := func1[int, A](0)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		i := 𝘃𝗮𝗹𝟮
		𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟮 := func1[int, A](i)
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
		𝘃𝗮𝗹𝟭 := 𝘃𝗮𝗹𝟯
		𝘃𝗮𝗹𝟰, 𝗲𝗿𝗿𝟯 := func1[int, B](i)
		if 𝗲𝗿𝗿𝟯 != nil {
			return 𝗲𝗿𝗿𝟯
		}
		if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟰 {
			println("B")
		} else {
			𝘃𝗮𝗹𝟱, 𝗲𝗿𝗿𝟰 := func1[int, C](i)
			if 𝗲𝗿𝗿𝟰 != nil {
				return 𝗲𝗿𝗿𝟰
			}
			if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟱 {
				println("C")
			} else {
				println("D")
			}
		}
	}
	return nil
}
func switch_multi_case() error {
	{
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟭 := func1[int, A](0)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		i := 𝘃𝗮𝗹𝟮
		𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟮 := func1[int, B](i)
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
		𝘃𝗮𝗹𝟭 := 𝘃𝗮𝗹𝟯
		𝘃𝗮𝗹𝟰, 𝗲𝗿𝗿𝟯 := func1[int, C](i)
		if 𝗲𝗿𝗿𝟯 != nil {
			return 𝗲𝗿𝗿𝟯
		}
		𝘃𝗮𝗹𝟲 := 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟰
		if !𝘃𝗮𝗹𝟲 {
			𝘃𝗮𝗹𝟱, 𝗲𝗿𝗿𝟰 := func1[int, D](i)
			if 𝗲𝗿𝗿𝟰 != nil {
				return 𝗲𝗿𝗿𝟰
			}
			𝘃𝗮𝗹𝟲 = 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟱
		}
		if 𝘃𝗮𝗹𝟲 {
			println("hello")
		} else {
			𝘃𝗮𝗹𝟳, 𝗲𝗿𝗿𝟱 := func1[int, E](i)
			if 𝗲𝗿𝗿𝟱 != nil {
				return 𝗲𝗿𝗿𝟱
			}
			if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟳 {
				println("hello")
			} else {
			}
		}
	}
	return nil
}
func switch_multi_case2() error {
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, A](1)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		𝘃𝗮𝗹𝟯 := 𝘃𝗮𝗹𝟭 == 1
		if !𝘃𝗮𝗹𝟯 {
			𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := func1[int, A](1)
			if 𝗲𝗿𝗿𝟮 != nil {
				return 𝗲𝗿𝗿𝟮
			}
			𝘃𝗮𝗹𝟯 = 𝘃𝗮𝗹𝟮 == 2
		}
		if 𝘃𝗮𝗹𝟯 {
			println("1,2")
		} else {
			𝘃𝗮𝗹𝟰, 𝗲𝗿𝗿𝟯 := func1[int, A](1)
			if 𝗲𝗿𝗿𝟯 != nil {
				return 𝗲𝗿𝗿𝟯
			}
			if 𝘃𝗮𝗹𝟰 == 3 {
				println("3")
			}
		}
	}
	return nil
}
func switch_break() error {
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, A](1)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		if 𝘃𝗮𝗹𝟭 == 42 {
			println("hello")
			goto 𝗟_𝗕𝗿𝗸𝗧𝗼𝟭
		}
	𝗟_𝗕𝗿𝗸𝗧𝗼𝟭:
	}
	return nil
}
func switch_nested_break() error {
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, A](1)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		if 𝘃𝗮𝗹𝟭 == 42 {
			println("hello")
			{
				𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := func1[int, A](1)
				if 𝗲𝗿𝗿𝟮 != nil {
					return 𝗲𝗿𝗿𝟮
				}
				if 𝘃𝗮𝗹𝟮 == 42 {
					println("hello")
					goto 𝗟_𝗕𝗿𝗸𝗧𝗼𝟭
				}
			𝗟_𝗕𝗿𝗸𝗧𝗼𝟭:
			}
			goto 𝗟_𝗕𝗿𝗸𝗧𝗼𝟮
		}
	𝗟_𝗕𝗿𝗸𝗧𝗼𝟮:
	}
	return nil
}
func switch_labeled_break() error {
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, A](1)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		if 𝘃𝗮𝗹𝟭 == 42 {
			println("hello")
			goto 𝗟_𝗕𝗿𝗸𝗧𝗼_𝗟𝟭
		}
	𝗟_𝗕𝗿𝗸𝗧𝗼_𝗟𝟭:
	}
	return nil
}
func switch_goto() error {
L:
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, A](1)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		if 𝘃𝗮𝗹𝟭 == 42 {
			goto L
		}
	}
	return nil
}
func switch_labeled_break_and_goto() error {
L:
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, A](1)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		if 𝘃𝗮𝗹𝟭 == 42 {
			goto 𝗟_𝗕𝗿𝗸𝗧𝗼_𝗟𝟭
		} else {
			𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := func1[int, B](1)
			if 𝗲𝗿𝗿𝟮 != nil {
				return 𝗲𝗿𝗿𝟮
			}
			if 𝘃𝗮𝗹𝟮 == 42 {
				goto L
			}
		}
	𝗟_𝗕𝗿𝗸𝗧𝗼_𝗟𝟭:
	}
	return nil
}
func switch_nested_goto() error {
L:
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, A](1)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		if 𝘃𝗮𝗹𝟭 == 42 {
			𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := func1[int, B](1)
			if 𝗲𝗿𝗿𝟮 != nil {
				return 𝗲𝗿𝗿𝟮
			}
			if 𝘃𝗮𝗹𝟮 == 42 {
				for {
					goto L
				}
			}
		}
	}
	return nil
}
func switch_nested_labeled_break() error {
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, A](1)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		if 𝘃𝗮𝗹𝟭 == 42 {
			println("outer")
			{
				𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := func1[int, B](1)
				if 𝗲𝗿𝗿𝟮 != nil {
					return 𝗲𝗿𝗿𝟮
				}
				if 𝘃𝗮𝗹𝟮 == 42 {
					println("inner")
					for {
						goto 𝗟_𝗕𝗿𝗸𝗧𝗼_𝗟𝟭
					}
				}
			}
		}
	𝗟_𝗕𝗿𝗸𝗧𝗼_𝗟𝟭:
	}
	return nil
}
func switch_nested_labeled_break_and_goto() error {
L:
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, A](1)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		if 𝘃𝗮𝗹𝟭 == 42 {
			println("outer")
			{
				𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := func1[int, B](1)
				if 𝗲𝗿𝗿𝟮 != nil {
					return 𝗲𝗿𝗿𝟮
				}
				if 𝘃𝗮𝗹𝟮 == 42 {
					println("inner")
					for {
						if true {
							goto 𝗟_𝗕𝗿𝗸𝗧𝗼_𝗟𝟭
						} else {
							goto L
						}
					}
				}
			}
		}
	𝗟_𝗕𝗿𝗸𝗧𝗼_𝗟𝟭:
	}
	return nil
}
func switch_nested() error {
outer:
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, A](1)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		if 𝘃𝗮𝗹𝟭 == 42 {
			println("outer")
		inner:
			{
				𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := func1[int, B](1)
				if 𝗲𝗿𝗿𝟮 != nil {
					return 𝗲𝗿𝗿𝟮
				}
				if 𝘃𝗮𝗹𝟮 == 42 {
					goto 𝗟_𝗕𝗿𝗸𝗧𝗼_𝗶𝗻𝗻𝗲𝗿𝟭
				} else {
					𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟯 := func1[int, C](1)
					if 𝗲𝗿𝗿𝟯 != nil {
						return 𝗲𝗿𝗿𝟯
					}
					if 𝘃𝗮𝗹𝟯 == 42 {
						goto inner
					} else {
						𝘃𝗮𝗹𝟰, 𝗲𝗿𝗿𝟰 := func1[int, D](1)
						if 𝗲𝗿𝗿𝟰 != nil {
							return 𝗲𝗿𝗿𝟰
						}
						if 𝘃𝗮𝗹𝟰 == 42 {
							println("inner")
							goto 𝗟_𝗕𝗿𝗸𝗧𝗼_𝗼𝘂𝘁𝗲𝗿𝟭
						} else {
							𝘃𝗮𝗹𝟱, 𝗲𝗿𝗿𝟱 := func1[int, E](1)
							if 𝗲𝗿𝗿𝟱 != nil {
								return 𝗲𝗿𝗿𝟱
							}
							if 𝘃𝗮𝗹𝟱 == 42 {
								println("inner")
								goto outer
							}
						}
					}
				}
			𝗟_𝗕𝗿𝗸𝗧𝗼_𝗶𝗻𝗻𝗲𝗿𝟭:
			}
		} else {
			println("default")
		}
	𝗟_𝗕𝗿𝗸𝗧𝗼_𝗼𝘂𝘁𝗲𝗿𝟭:
	}
	return nil
}
func switch_fallthrough_default() error {
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, int](1)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		if 𝘃𝗮𝗹𝟭 == 42 {
			println("hello")
		} else {
			{
				println("fallthrough")
			}
			{
				println("hello")
			}
		}
	}
	return nil
}
func switch_labeled_fallthrough() error {
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, int](1)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		if 𝘃𝗮𝗹𝟭 == 42 {
			{
				println("hello")
				if false {
					goto L
				}
			L:
			}
			{
				println("default")
			}
		} else {
			println("default")
		}
	}
	return nil
}
func switch_mixed() error {
𝗟_𝗚𝗼𝘁𝗼_𝗟𝟭:
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, A](0)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		i := 𝘃𝗮𝗹𝟭
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := func1[int, B](i)
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
		if 𝘃𝗮𝗹𝟮 == 42 {
			{
				switch {
				default:
					for i := 0; i < 1; i++ {
						goto labeldFall
					}
				}
				println("B")
			labeldFall:
			}
			{
				println("default")
			}
			{
				println("42-a")
			}
		} else if id[int](i) == 42 {
			println("42-a")
		} else {
			𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟯 := func1[int, C](i)
			if 𝗲𝗿𝗿𝟯 != nil {
				return 𝗲𝗿𝗿𝟯
			}
			if 𝘃𝗮𝗹𝟯 == 42 {
				println("C")
				goto 𝗟_𝗕𝗿𝗸𝗧𝗼_𝗟𝟭
			} else {
				𝘃𝗮𝗹𝟰, 𝗲𝗿𝗿𝟰 := func1[int, C](i)
				if 𝗲𝗿𝗿𝟰 != nil {
					return 𝗲𝗿𝗿𝟰
				}
				if 𝘃𝗮𝗹𝟰 == 42 {
					for i := 0; i < 10; i++ {
						𝘃𝗮𝗹𝟱, 𝗲𝗿𝗿𝟱 := func1[int, error](0)
						if 𝗲𝗿𝗿𝟱 != nil {
							return 𝗲𝗿𝗿𝟱
						}
						switch 𝘃𝗮𝗹𝟱.(type) {
						case error:
							println("C2")
							goto 𝗟_𝗕𝗿𝗸𝗧𝗼_𝗟𝟭
						case nil:
							for {
								goto 𝗟_𝗚𝗼𝘁𝗼_𝗟𝟭
							}
						}
					}
				} else if id[int](i) == 42 {
					println("42-b")
					goto 𝗟_𝗚𝗼𝘁𝗼_𝗟𝟭
				} else {
					{
						println("default")
					}
					{
						println("42-a")
					}
				}
			}
		}
	𝗟_𝗕𝗿𝗸𝗧𝗼_𝗟𝟭:
	}
	goto 𝗟_𝗚𝗼𝘁𝗼_𝗟𝟭
}
