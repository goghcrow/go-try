//go:build !try

// Code generated by github.com/goghcrow/go-try DO NOT EDIT.
package test

import (
	"errors"
	. "github.com/goghcrow/go-try/rt"
)

type (
	Int int
	Str string
)

var helloErr = errors.New("hello error!")

func answer() (int, error) {
	var 𝘇𝗲𝗿𝗼𝟬 int
	𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟮 := ret1Err[int]()
	if 𝗲𝗿𝗿𝟮 != nil {
		return 𝘇𝗲𝗿𝗼𝟬, 𝗲𝗿𝗿𝟮
	}
	return 𝘃𝗮𝗹𝟭, nil
}

func assign() (_ error) {
	𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟮 := ret1Err[int]()
	if 𝗲𝗿𝗿𝟮 != nil {
		return 𝗲𝗿𝗿𝟮
	}
	x, y := 𝘃𝗮𝗹𝟭, 42
	consume2(x, y)
	𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟰 := 42, helloErr
	if 𝗲𝗿𝗿𝟰 != nil {
		return 𝗲𝗿𝗿𝟰
	}
	v1, v2 := 𝘃𝗮𝗹𝟯, 42
	consume2(v1, v2)
	return nil
}

func binary() (_ error) {
	𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟮 := ret1Err[int]()
	if 𝗲𝗿𝗿𝟮 != nil {
		return 𝗲𝗿𝗿𝟮
	}
	consume1(𝘃𝗮𝗹𝟭 + 1)
	return nil
}

func ret0() (_ error) {
	𝗲𝗿𝗿𝟭 := ret0Err()
	if 𝗲𝗿𝗿𝟭 != nil {
		return 𝗲𝗿𝗿𝟭
	}
	Ø()
	𝗲𝗿𝗿𝟮 := helloErr
	if 𝗲𝗿𝗿𝟮 != nil {
		return 𝗲𝗿𝗿𝟮
	}
	Ø()
	return
}

func ret1() (_ Int, _ error) {
	var 𝘇𝗲𝗿𝗼𝟬 Int
	_, 𝗲𝗿𝗿𝟭 := ret1Err[int]()
	if 𝗲𝗿𝗿𝟭 != nil {
		return 𝘇𝗲𝗿𝗼𝟬, 𝗲𝗿𝗿𝟭
	}
	Ø()
	𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟯 := ret1Err[int]()
	if 𝗲𝗿𝗿𝟯 != nil {
		return 𝘇𝗲𝗿𝗼𝟬, 𝗲𝗿𝗿𝟯
	}
	consume1(𝘃𝗮𝗹𝟮 + 1)
	_, 𝗲𝗿𝗿𝟰 := 42, helloErr
	if 𝗲𝗿𝗿𝟰 != nil {
		return 𝘇𝗲𝗿𝗼𝟬, 𝗲𝗿𝗿𝟰
	}
	Ø()
	𝘃𝗮𝗹𝟱, 𝗲𝗿𝗿𝟲 := 42, helloErr
	if 𝗲𝗿𝗿𝟲 != nil {
		return 𝘇𝗲𝗿𝗼𝟬, 𝗲𝗿𝗿𝟲
	}
	consume1(𝘃𝗮𝗹𝟱 + 1)

	return
}

func ret2() (_ Int, _ Str, _ error) {
	var (
		𝘇𝗲𝗿𝗼𝟬 Int
		𝘇𝗲𝗿𝗼𝟭 Str
	)
	_, _, 𝗲𝗿𝗿𝟭 := ret2Err[int, string]()
	if 𝗲𝗿𝗿𝟭 != nil {
		return 𝘇𝗲𝗿𝗼𝟬, 𝘇𝗲𝗿𝗼𝟭, 𝗲𝗿𝗿𝟭
	}
	Ø()
	𝘃𝗮𝗹𝟮, 𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟰 := ret2Err[int, string]()
	if 𝗲𝗿𝗿𝟰 != nil {
		return 𝘇𝗲𝗿𝗼𝟬, 𝘇𝗲𝗿𝗼𝟭, 𝗲𝗿𝗿𝟰
	}
	iV, bV := II(𝘃𝗮𝗹𝟮, 𝘃𝗮𝗹𝟯)
	consume2(iV, bV)
	𝘃𝗮𝗹𝟱, 𝘃𝗮𝗹𝟲, 𝗲𝗿𝗿𝟳 := ret2Err[int, string]()
	if 𝗲𝗿𝗿𝟳 != nil {
		return 𝘇𝗲𝗿𝗼𝟬, 𝘇𝗲𝗿𝗼𝟭, 𝗲𝗿𝗿𝟳
	}
	consume2(II(𝘃𝗮𝗹𝟱, 𝘃𝗮𝗹𝟲))
	_, _, 𝗲𝗿𝗿𝟴 := 42, "answer", helloErr
	if 𝗲𝗿𝗿𝟴 != nil {
		return 𝘇𝗲𝗿𝗼𝟬, 𝘇𝗲𝗿𝗼𝟭, 𝗲𝗿𝗿𝟴
	}
	Ø()
	𝘃𝗮𝗹𝟵, 𝘃𝗮𝗹𝟭𝟬, 𝗲𝗿𝗿𝟭𝟭 := 42, "answer", helloErr
	if 𝗲𝗿𝗿𝟭𝟭 != nil {
		return 𝘇𝗲𝗿𝗼𝟬, 𝘇𝗲𝗿𝗼𝟭, 𝗲𝗿𝗿𝟭𝟭
	}
	iV, bV = II(𝘃𝗮𝗹𝟵, 𝘃𝗮𝗹𝟭𝟬)
	consume2(iV, bV)
	𝘃𝗮𝗹𝟭𝟮, 𝘃𝗮𝗹𝟭𝟯, 𝗲𝗿𝗿𝟭𝟰 := 42, "answer", helloErr
	if 𝗲𝗿𝗿𝟭𝟰 != nil {
		return 𝘇𝗲𝗿𝗼𝟬, 𝘇𝗲𝗿𝗼𝟭, 𝗲𝗿𝗿𝟭𝟰
	}
	consume2(II(𝘃𝗮𝗹𝟭𝟮, 𝘃𝗮𝗹𝟭𝟯))
	return
}

func ret2_grouped_ret() (_, _ Int, _ error) {
	var 𝘇𝗲𝗿𝗼𝟬, 𝘇𝗲𝗿𝗼𝟭 Int
	_, _, 𝗲𝗿𝗿𝟭 := ret2Err[int, byte]()
	if 𝗲𝗿𝗿𝟭 != nil {
		return 𝘇𝗲𝗿𝗼𝟬, 𝘇𝗲𝗿𝗼𝟭, 𝗲𝗿𝗿𝟭
	}
	Ø()
	return
}

func ret3() (_ *Int, _ error) {
	var 𝘇𝗲𝗿𝗼𝟬 *Int
	_, _, _, 𝗲𝗿𝗿𝟭 := ret3Err[int, rune, string]()
	if 𝗲𝗿𝗿𝟭 != nil {
		return 𝘇𝗲𝗿𝗼𝟬, 𝗲𝗿𝗿𝟭
	}
	Ø()
	𝘃𝗮𝗹𝟮, 𝘃𝗮𝗹𝟯, 𝘃𝗮𝗹𝟰, 𝗲𝗿𝗿𝟱 := ret3Err[int, rune, string]()
	if 𝗲𝗿𝗿𝟱 != nil {
		return 𝘇𝗲𝗿𝗼𝟬, 𝗲𝗿𝗿𝟱
	}
	iV, bV, sV := III(𝘃𝗮𝗹𝟮, 𝘃𝗮𝗹𝟯, 𝘃𝗮𝗹𝟰)
	consume3(iV, bV, sV)
	𝘃𝗮𝗹𝟲, 𝘃𝗮𝗹𝟳, 𝘃𝗮𝗹𝟴, 𝗲𝗿𝗿𝟵 := ret3Err[int, rune, string]()
	if 𝗲𝗿𝗿𝟵 != nil {
		return 𝘇𝗲𝗿𝗼𝟬, 𝗲𝗿𝗿𝟵
	}
	consume3(III(𝘃𝗮𝗹𝟲, 𝘃𝗮𝗹𝟳, 𝘃𝗮𝗹𝟴))
	_, _, _, 𝗲𝗿𝗿𝟭𝟬 := 42, 'a', "hello", helloErr
	if 𝗲𝗿𝗿𝟭𝟬 != nil {
		return 𝘇𝗲𝗿𝗼𝟬, 𝗲𝗿𝗿𝟭𝟬
	}
	Ø()
	𝘃𝗮𝗹𝟭𝟭, 𝘃𝗮𝗹𝟭𝟮, 𝘃𝗮𝗹𝟭𝟯, 𝗲𝗿𝗿𝟭𝟰 := 42, 'a', "hello", helloErr
	if 𝗲𝗿𝗿𝟭𝟰 != nil {
		return 𝘇𝗲𝗿𝗼𝟬, 𝗲𝗿𝗿𝟭𝟰
	}
	iV, bV, sV = III(𝘃𝗮𝗹𝟭𝟭, 𝘃𝗮𝗹𝟭𝟮, 𝘃𝗮𝗹𝟭𝟯)
	consume3(iV, bV, sV)
	𝘃𝗮𝗹𝟭𝟱, 𝘃𝗮𝗹𝟭𝟲, 𝘃𝗮𝗹𝟭𝟳, 𝗲𝗿𝗿𝟭𝟴 := 42, 'a', "hello", helloErr
	if 𝗲𝗿𝗿𝟭𝟴 != nil {
		return 𝘇𝗲𝗿𝗼𝟬, 𝗲𝗿𝗿𝟭𝟴
	}
	consume3(III(𝘃𝗮𝗹𝟭𝟱, 𝘃𝗮𝗹𝟭𝟲, 𝘃𝗮𝗹𝟭𝟳))
	𝘃𝗮𝗹𝟭𝟵, 𝘃𝗮𝗹𝟮𝟬, 𝘃𝗮𝗹𝟮𝟭, 𝗲𝗿𝗿𝟮𝟮 := ret3Err[int, rune, string]()
	if 𝗲𝗿𝗿𝟮𝟮 != nil {
		return 𝘇𝗲𝗿𝗼𝟬, 𝗲𝗿𝗿𝟮𝟮
	}
	func(int, rune, string) {}(III(𝘃𝗮𝗹𝟭𝟵, 𝘃𝗮𝗹𝟮𝟬, 𝘃𝗮𝗹𝟮𝟭))
	return
}

func funcLit() {
	go func() {
		_ = func() error {
			𝗲𝗿𝗿𝟭 := ret0()
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			Ø()
			𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟯 := ret1Err[int]()
			if 𝗲𝗿𝗿𝟯 != nil {
				return 𝗲𝗿𝗿𝟯
			}
			consume1(𝘃𝗮𝗹𝟮)
			return nil
		}()
	}()

	defer func() {
		_ = func() error {
			𝗲𝗿𝗿𝟭 := ret0()
			if 𝗲𝗿𝗿𝟭 != nil {
				return 𝗲𝗿𝗿𝟭
			}
			Ø()
			𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟯 := ret1Err[int]()
			if 𝗲𝗿𝗿𝟯 != nil {
				return 𝗲𝗿𝗿𝟯
			}
			consume1(𝘃𝗮𝗹𝟮)
			return nil
		}()
	}()

	if func() int {
		func() (int, error) {
			var 𝘇𝗲𝗿𝗼𝟬 int
			𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟮 := ret1Err[int]()
			if 𝗲𝗿𝗿𝟮 != nil {
				return 𝘇𝗲𝗿𝗼𝟬, 𝗲𝗿𝗿𝟮
			}
			return id(𝘃𝗮𝗹𝟭), nil
		}()
		return 42
	}() == 42 {
	}
}
