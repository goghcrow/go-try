//go:build !try

// Code generated by github.com/goghcrow/go-try DO NOT EDIT.
package test

import . "github.com/goghcrow/go-try/rt"

func for_test() error {
	for {
		𝗽𝗼𝘀𝘁𝟭 := func() (_ E𝗿𝗿𝗼𝗿) {
			𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[A]()
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			id(𝘃𝗮𝗹𝟭)
			return
		}
		𝗲𝗿𝗿𝟭 := 𝗽𝗼𝘀𝘁𝟭()
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
	}
	for {
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟮 := ret1Err[A]()
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
		𝘃𝗮𝗹𝟮 := id(𝘃𝗮𝗹𝟭)
		if 𝘃𝗮𝗹𝟮 != 42 {
			break
		}
	}
	for i := 1; ; {
		𝗽𝗼𝘀𝘁𝟮 := func() (_ E𝗿𝗿𝗼𝗿) {
			𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, int](i)
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			id(𝘃𝗮𝗹𝟭)
			return
		}
		𝗲𝗿𝗿𝟯 := 𝗽𝗼𝘀𝘁𝟮()
		if 𝗲𝗿𝗿𝟯 != nil {
			return 𝗲𝗿𝗿𝟯
		}
	}
	for i := 1; ; {
		𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟰 := func1[int, int](i)
		if 𝗲𝗿𝗿𝟰 != nil {
			return 𝗲𝗿𝗿𝟰
		}
		𝘃𝗮𝗹𝟰 := id(𝘃𝗮𝗹𝟯)
		if 𝘃𝗮𝗹𝟰 != 42 {
			break
		}
	}
	𝘃𝗮𝗹𝟱, 𝗲𝗿𝗿𝟱 := ret1Err[A]()
	if 𝗲𝗿𝗿𝟱 != nil {
		return 𝗲𝗿𝗿𝟱
	}
	for i := 𝘃𝗮𝗹𝟱; ; {
		𝗽𝗼𝘀𝘁𝟯 := func() (_ E𝗿𝗿𝗼𝗿) {
			_, 𝗲𝗿𝗿𝟭 := func1[A, C](i)
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			return
		}
		𝘃𝗮𝗹𝟲, 𝗲𝗿𝗿𝟲 := func1[int, bool](i)
		if 𝗲𝗿𝗿𝟲 != nil {
			return 𝗲𝗿𝗿𝟲
		}
		if !𝘃𝗮𝗹𝟲 {
			break
		}
		𝗲𝗿𝗿𝟳 := 𝗽𝗼𝘀𝘁𝟯()
		if 𝗲𝗿𝗿𝟳 != nil {
			return 𝗲𝗿𝗿𝟳
		}
	}
	return nil
}
func for_try_in_init() error {
	𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, int](0)
	if 𝗲𝗿𝗿𝟭 != nil {
		return 𝗲𝗿𝗿𝟭
	}
	for i := 𝘃𝗮𝗹𝟭; ; {
		println(i)
	}
}
func for_try_in_cond() error {
	for i := 0; ; i++ {
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, int](i)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		if i >= 𝘃𝗮𝗹𝟭 {
			break
		}
		println(i)
	}
	return nil
}
func for_try_in_cond1() error {
	for i := 0; ; i++ {
		𝘃𝗮𝗹𝟭 := id[bool](false)
		𝘃𝗮𝗹𝟯 := !𝘃𝗮𝗹𝟭
		if 𝘃𝗮𝗹𝟯 {
			𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟭 := func1[int, int](i)
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			𝘃𝗮𝗹𝟯 = 𝘃𝗮𝗹𝟮 <= 1
		}
		if 𝘃𝗮𝗹𝟯 {
			break
		}
		println(i)
	}
	return nil
}
func for_try_in_post() error {
	for i := 0; i < 42; {
		𝗽𝗼𝘀𝘁𝟭 := func() (_ E𝗿𝗿𝗼𝗿) {
			_, 𝗲𝗿𝗿𝟭 := func1[int, int](i)
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			return
		}
		println(i)
		i++
		𝗲𝗿𝗿𝟭 := 𝗽𝗼𝘀𝘁𝟭()
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
	}
	return nil
}
func for_try_in_post1() error {
	for {
		𝗽𝗼𝘀𝘁𝟭 := func() (_ E𝗿𝗿𝗼𝗿) {
			𝗲𝗿𝗿𝟭 := ret0()
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			return
		}
		𝗲𝗿𝗿𝟭 := 𝗽𝗼𝘀𝘁𝟭()
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
	}
	return nil
}
func for_try_in_post2() error {
	for {
		𝗽𝗼𝘀𝘁𝟭 := func() (_ E𝗿𝗿𝗼𝗿) {
			𝗲𝗿𝗿𝟭 := ret0()
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			return
		}
		println(1)
		𝗲𝗿𝗿𝟭 := 𝗽𝗼𝘀𝘁𝟭()
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		continue
	}
	return nil
}
func for_try_in_post21() error {
	for {
		𝗽𝗼𝘀𝘁𝟭 := func() (_ E𝗿𝗿𝗼𝗿) {
			𝗲𝗿𝗿𝟭 := ret0()
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			return
		}
		println(1)
		𝗲𝗿𝗿𝟭 := 𝗽𝗼𝘀𝘁𝟭()
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		continue
		println(2)
		𝗲𝗿𝗿𝟮 := 𝗽𝗼𝘀𝘁𝟭()
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
	}
	return nil
}
func for_try_in_post22() error {
	for {
		𝗽𝗼𝘀𝘁𝟭 := func() (_ E𝗿𝗿𝗼𝗿) {
			𝗲𝗿𝗿𝟭 := ret0()
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			return
		}
		println(1)
		panic(nil)
		if 42 != 42 {
			𝗲𝗿𝗿𝟭 := 𝗽𝗼𝘀𝘁𝟭()
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
		}
	}
	return nil
}
func for_try_in_post23() error {
	for {
		𝗽𝗼𝘀𝘁𝟭 := func() (_ E𝗿𝗿𝗼𝗿) {
			𝗲𝗿𝗿𝟭 := ret0()
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			return
		}
		println(1)
		panic(nil)
		println(2)
		𝗲𝗿𝗿𝟭 := 𝗽𝗼𝘀𝘁𝟭()
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
	}
	return nil
}
func for_try_in_post24() error {
	for {
		𝗽𝗼𝘀𝘁𝟭 := func() (_ E𝗿𝗿𝗼𝗿) {
			𝗲𝗿𝗿𝟭 := ret0()
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			return
		}
		𝗲𝗿𝗿𝟭 := 𝗽𝗼𝘀𝘁𝟭()
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		continue
		x := 1
		println(x)
		𝗲𝗿𝗿𝟮 := 𝗽𝗼𝘀𝘁𝟭()
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
	}
	return nil
}
func for_try_in_post3() error {
L:
	for {
		𝗽𝗼𝘀𝘁𝟭 := func() (_ E𝗿𝗿𝗼𝗿) {
			𝗲𝗿𝗿𝟭 := ret0()
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			return
		}
		𝗲𝗿𝗿𝟭 := 𝗽𝗼𝘀𝘁𝟭()
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		continue L
	}
	return nil
}
func for_try_in_post4() error {
L:
	for {
		𝗽𝗼𝘀𝘁𝟮 := func() (_ E𝗿𝗿𝗼𝗿) {
			_, 𝗲𝗿𝗿𝟭 := func1[int, int](0)
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			return
		}
		for {
			𝗽𝗼𝘀𝘁𝟭 := func() (_ E𝗿𝗿𝗼𝗿) {
				_, 𝗲𝗿𝗿𝟭 := func1[int, int](1)
				if 𝗲𝗿𝗿𝟭 != nil {
					return 𝗲𝗿𝗿𝟭
				}
				return
			}
			𝗲𝗿𝗿𝟭 := 𝗽𝗼𝘀𝘁𝟮()
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			continue L
			if 42 != 42 {
				𝗲𝗿𝗿𝟮 := 𝗽𝗼𝘀𝘁𝟭()
				if 𝗲𝗿𝗿𝟮 != nil {
					return 𝗲𝗿𝗿𝟮
				}
			}
		}
	}
	return nil
}
func for_try_in_post5() error {
outer:
	for {
		𝗽𝗼𝘀𝘁𝟮 := func() (_ E𝗿𝗿𝗼𝗿) {
			_, 𝗲𝗿𝗿𝟭 := func1[int, int](1)
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			return
		}
	inner:
		for {
			𝗽𝗼𝘀𝘁𝟭 := func() (_ E𝗿𝗿𝗼𝗿) {
				_, 𝗲𝗿𝗿𝟭 := func1[int, int](2)
				if 𝗲𝗿𝗿𝟭 != nil {
					return 𝗲𝗿𝗿𝟭
				}
				return
			}
			if true {
				𝗲𝗿𝗿𝟭 := 𝗽𝗼𝘀𝘁𝟭()
				if 𝗲𝗿𝗿𝟭 != nil {
					return 𝗲𝗿𝗿𝟭
				}
				continue inner
			} else {
				𝗲𝗿𝗿𝟮 := 𝗽𝗼𝘀𝘁𝟮()
				if 𝗲𝗿𝗿𝟮 != nil {
					return 𝗲𝗿𝗿𝟮
				}
				continue outer
			}
			𝗲𝗿𝗿𝟯 := 𝗽𝗼𝘀𝘁𝟭()
			if 𝗲𝗿𝗿𝟯 != nil {
				return 𝗲𝗿𝗿𝟯
			}
		}
		𝗲𝗿𝗿𝟰 := 𝗽𝗼𝘀𝘁𝟮()
		if 𝗲𝗿𝗿𝟰 != nil {
			return 𝗲𝗿𝗿𝟰
		}
		continue outer
	}
	return nil
}
func for_try_in_init_cond_post() error {
	𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, int](42)
	if 𝗲𝗿𝗿𝟭 != nil {
		return 𝗲𝗿𝗿𝟭
	}
	for i := 𝘃𝗮𝗹𝟭; ; {
		𝗽𝗼𝘀𝘁𝟭 := func() (_ E𝗿𝗿𝗼𝗿) {
			_, 𝗲𝗿𝗿𝟭 := func1[int, int](i)
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			return
		}
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := func1[int, int](i)
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
		if i >= 𝘃𝗮𝗹𝟮 {
			break
		}
		𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟯 := func1[int, int](i)
		if 𝗲𝗿𝗿𝟯 != nil {
			return 𝗲𝗿𝗿𝟯
		}
		println(𝘃𝗮𝗹𝟯)
		𝗲𝗿𝗿𝟰 := 𝗽𝗼𝘀𝘁𝟭()
		if 𝗲𝗿𝗿𝟰 != nil {
			return 𝗲𝗿𝗿𝟰
		}
	}
	return nil
}
func for_labeled_brk() error {
	𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, int](42)
	if 𝗲𝗿𝗿𝟭 != nil {
		return 𝗲𝗿𝗿𝟭
	}
L:
	for i := 𝘃𝗮𝗹𝟭; ; {
		𝗽𝗼𝘀𝘁𝟭 := func() (_ E𝗿𝗿𝗼𝗿) {
			_, 𝗲𝗿𝗿𝟭 := func1[int, int](i)
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			return
		}
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := func1[int, int](i)
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
		if i >= 𝘃𝗮𝗹𝟮 {
			break
		}
		break L
		𝗲𝗿𝗿𝟯 := 𝗽𝗼𝘀𝘁𝟭()
		if 𝗲𝗿𝗿𝟯 != nil {
			return 𝗲𝗿𝗿𝟯
		}
	}
	return nil
}
func for_labeled_continue() error {
	𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, int](42)
	if 𝗲𝗿𝗿𝟭 != nil {
		return 𝗲𝗿𝗿𝟭
	}
L:
	for i := 𝘃𝗮𝗹𝟭; ; {
		𝗽𝗼𝘀𝘁𝟭 := func() (_ E𝗿𝗿𝗼𝗿) {
			_, 𝗲𝗿𝗿𝟭 := func1[int, int](i)
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			return
		}
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := func1[int, int](i)
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
		if i >= 𝘃𝗮𝗹𝟮 {
			break
		}
		println(i)
		𝗲𝗿𝗿𝟯 := 𝗽𝗼𝘀𝘁𝟭()
		if 𝗲𝗿𝗿𝟯 != nil {
			return 𝗲𝗿𝗿𝟯
		}
		continue L
	}
	return nil
}
func for_labeled_cont_brk_goto() error {
L:
	for i := 0; i < 42; {
		𝗽𝗼𝘀𝘁𝟭 := func() (_ E𝗿𝗿𝗼𝗿) {
			_, 𝗲𝗿𝗿𝟭 := func1[int, int](i)
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			return
		}
		println(i)
		if i == 42 {
			𝗲𝗿𝗿𝟭 := 𝗽𝗼𝘀𝘁𝟭()
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			continue
		}
		if i == 42 {
			𝗲𝗿𝗿𝟮 := 𝗽𝗼𝘀𝘁𝟭()
			if 𝗲𝗿𝗿𝟮 != nil {
				return 𝗲𝗿𝗿𝟮
			}
			continue L
		}
		if i == 42 {
			goto L
		}
		if i == 42 {
			break
		} else {
			i++
		}
		𝗲𝗿𝗿𝟯 := 𝗽𝗼𝘀𝘁𝟭()
		if 𝗲𝗿𝗿𝟯 != nil {
			return 𝗲𝗿𝗿𝟯
		}
	}
	return nil
}
func for_nested_labeled() error {
	type (
		innerPost = int
		outerPost = int
	)
𝗟_𝗚𝗼𝘁𝗼_𝗼𝘂𝘁𝗲𝗿𝟭:
	{
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟭 := func1[int, int](1)
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
	outer:
		for i := 𝘃𝗮𝗹𝟮; ; {
			𝗽𝗼𝘀𝘁𝟮 := func() (_ E𝗿𝗿𝗼𝗿) {
				_, 𝗲𝗿𝗿𝟭 := func1[int, outerPost](i + 2)
				if 𝗲𝗿𝗿𝟭 != nil {
					return 𝗲𝗿𝗿𝟭
				}
				return
			}
			𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟮 := func1[int, int](i + 1)
			if 𝗲𝗿𝗿𝟮 != nil {
				return 𝗲𝗿𝗿𝟮
			}
			if i >= 𝘃𝗮𝗹𝟯 {
				break
			}
		𝗟_𝗚𝗼𝘁𝗼_𝗶𝗻𝗻𝗲𝗿𝟭:
			{
				𝘃𝗮𝗹𝟰, 𝗲𝗿𝗿𝟯 := func1[int, int](2)
				if 𝗲𝗿𝗿𝟯 != nil {
					return 𝗲𝗿𝗿𝟯
				}
			inner:
				for j := 𝘃𝗮𝗹𝟰; ; {
					𝗽𝗼𝘀𝘁𝟭 := func() (_ E𝗿𝗿𝗼𝗿) {
						_, 𝗲𝗿𝗿𝟭 := func1[int, innerPost](j + 4)
						if 𝗲𝗿𝗿𝟭 != nil {
							return 𝗲𝗿𝗿𝟭
						}
						return
					}
					𝘃𝗮𝗹𝟱, 𝗲𝗿𝗿𝟰 := func1[int, int](j + 3)
					if 𝗲𝗿𝗿𝟰 != nil {
						return 𝗲𝗿𝗿𝟰
					}
					if j >= 𝘃𝗮𝗹𝟱 {
						break
					}
					{
						𝘃𝗮𝗹𝟲, 𝗲𝗿𝗿𝟱 := func1[int, int](3)
						if 𝗲𝗿𝗿𝟱 != nil {
							return 𝗲𝗿𝗿𝟱
						}
						𝘃𝗮𝗹𝟭 := 𝘃𝗮𝗹𝟲
						𝘃𝗮𝗹𝟳, 𝗲𝗿𝗿𝟲 := func1[int, int](4)
						if 𝗲𝗿𝗿𝟲 != nil {
							return 𝗲𝗿𝗿𝟲
						}
						if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟳 {
							𝗲𝗿𝗿𝟳 := 𝗽𝗼𝘀𝘁𝟮()
							if 𝗲𝗿𝗿𝟳 != nil {
								return 𝗲𝗿𝗿𝟳
							}
							continue outer
						} else {
							𝘃𝗮𝗹𝟴, 𝗲𝗿𝗿𝟴 := func1[int, int](5)
							if 𝗲𝗿𝗿𝟴 != nil {
								return 𝗲𝗿𝗿𝟴
							}
							if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟴 {
								𝗲𝗿𝗿𝟵 := 𝗽𝗼𝘀𝘁𝟭()
								if 𝗲𝗿𝗿𝟵 != nil {
									return 𝗲𝗿𝗿𝟵
								}
								continue inner
							} else {
								𝘃𝗮𝗹𝟵, 𝗲𝗿𝗿𝟭𝟬 := func1[int, int](6)
								if 𝗲𝗿𝗿𝟭𝟬 != nil {
									return 𝗲𝗿𝗿𝟭𝟬
								}
								if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟵 {
									break outer
								} else {
									𝘃𝗮𝗹𝟭𝟬, 𝗲𝗿𝗿𝟭𝟭 := func1[int, int](7)
									if 𝗲𝗿𝗿𝟭𝟭 != nil {
										return 𝗲𝗿𝗿𝟭𝟭
									}
									if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟭𝟬 {
										break inner
									} else {
										𝘃𝗮𝗹𝟭𝟭, 𝗲𝗿𝗿𝟭𝟮 := func1[int, int](8)
										if 𝗲𝗿𝗿𝟭𝟮 != nil {
											return 𝗲𝗿𝗿𝟭𝟮
										}
										if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟭𝟭 {
											goto 𝗟_𝗚𝗼𝘁𝗼_𝗼𝘂𝘁𝗲𝗿𝟭
										} else {
											𝘃𝗮𝗹𝟭𝟮, 𝗲𝗿𝗿𝟭𝟯 := func1[int, int](9)
											if 𝗲𝗿𝗿𝟭𝟯 != nil {
												return 𝗲𝗿𝗿𝟭𝟯
											}
											if 𝘃𝗮𝗹𝟭 == 𝘃𝗮𝗹𝟭𝟮 {
												goto 𝗟_𝗚𝗼𝘁𝗼_𝗶𝗻𝗻𝗲𝗿𝟭
											} else {
												goto 𝗟_𝗕𝗿𝗸𝗧𝗼𝟭
											}
										}
									}
								}
							}
						}
					𝗟_𝗕𝗿𝗸𝗧𝗼𝟭:
					}
					𝗲𝗿𝗿𝟭𝟰 := 𝗽𝗼𝘀𝘁𝟭()
					if 𝗲𝗿𝗿𝟭𝟰 != nil {
						return 𝗲𝗿𝗿𝟭𝟰
					}
				}
			}
			𝗲𝗿𝗿𝟭𝟱 := 𝗽𝗼𝘀𝘁𝟮()
			if 𝗲𝗿𝗿𝟭𝟱 != nil {
				return 𝗲𝗿𝗿𝟭𝟱
			}
		}
	}
	return nil
}
