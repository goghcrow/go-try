//go:build !try

// Code generated by github.com/goghcrow/go-try DO NOT EDIT.
package test

import . "github.com/goghcrow/go-try/rt"

func swith1() error {
	type (
		A = int
		B = int
		C = int
		D = int
	)
	𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[A]()
	if 𝗲𝗿𝗿𝟭 != nil {
		return 𝗲𝗿𝗿𝟭
	}
	switch 𝘃𝗮𝗹𝟭 {
	case 0:
		𝗲𝗿𝗿𝟮 := ret0Err()
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}

	case 1:
		_, 𝗲𝗿𝗿𝟯 := ret1Err[B]()
		if 𝗲𝗿𝗿𝟯 != nil {
			return 𝗲𝗿𝗿𝟯
		}

	case 2:
		_, _, 𝗲𝗿𝗿𝟰 := ret2Err[C, C]()
		if 𝗲𝗿𝗿𝟰 != nil {
			return 𝗲𝗿𝗿𝟰
		}

	case 3:
		_, _, _, 𝗲𝗿𝗿𝟱 := ret3Err[D, D, D]()
		if 𝗲𝗿𝗿𝟱 != nil {
			return 𝗲𝗿𝗿𝟱
		}

	}
	return nil
}
func swith2() error {
	type (
		A = int
		B = int
		C = int
		D = int
		E = int
	)
	𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[A]()
	if 𝗲𝗿𝗿𝟭 != nil {
		return 𝗲𝗿𝗿𝟭
	}
	𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := ret1Err[B]()
	if 𝗲𝗿𝗿𝟮 != nil {
		return 𝗲𝗿𝗿𝟮
	}
	switch a, b := 𝘃𝗮𝗹𝟭, 𝘃𝗮𝗹𝟮; {
	case a == 1:
		𝗲𝗿𝗿𝟯 := ret0Err()
		if 𝗲𝗿𝗿𝟯 != nil {
			return 𝗲𝗿𝗿𝟯
		}

	case b == 2:
		_, 𝗲𝗿𝗿𝟰 := ret1Err[C]()
		if 𝗲𝗿𝗿𝟰 != nil {
			return 𝗲𝗿𝗿𝟰
		}

	case a == 3:
		_, _, 𝗲𝗿𝗿𝟱 := ret2Err[D, D]()
		if 𝗲𝗿𝗿𝟱 != nil {
			return 𝗲𝗿𝗿𝟱
		}

	case b == 4:
		_, _, _, 𝗲𝗿𝗿𝟲 := ret3Err[E, E, E]()
		if 𝗲𝗿𝗿𝟲 != nil {
			return 𝗲𝗿𝗿𝟲
		}

	}
	return nil
}
func switch3() error {
	𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[error]()
	if 𝗲𝗿𝗿𝟭 != nil {
		return 𝗲𝗿𝗿𝟭
	}
	switch 𝘃𝗮𝗹𝟭.(type) {
	}
	return nil
}
func if1() error {
	type (
		A = int
		B = int
		C = int
	)
	if true {
		𝗲𝗿𝗿𝟭 := ret0Err()
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}

		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟮 := ret1Err[A]()
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
		println(𝘃𝗮𝗹𝟭)
		𝘃𝗮𝗹𝟮, 𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟯 := ret2Err[B, string]()
		if 𝗲𝗿𝗿𝟯 != nil {
			return 𝗲𝗿𝗿𝟯
		}
		println(II(𝘃𝗮𝗹𝟮, 𝘃𝗮𝗹𝟯))
		𝘃𝗮𝗹𝟰, 𝘃𝗮𝗹𝟱, 𝘃𝗮𝗹𝟲, 𝗲𝗿𝗿𝟰 := ret3Err[C, string, rune]()
		if 𝗲𝗿𝗿𝟰 != nil {
			return 𝗲𝗿𝗿𝟰
		}
		println(III(𝘃𝗮𝗹𝟰, 𝘃𝗮𝗹𝟱, 𝘃𝗮𝗹𝟲))
	}
	return nil
}
func for1() error {
	type (
		A = int
		B = int
		C = int
	)
	for {
		𝗲𝗿𝗿𝟭 := ret0Err()
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}

		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟮 := ret1Err[A]()
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
		println(𝘃𝗮𝗹𝟭)
		𝘃𝗮𝗹𝟮, 𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟯 := ret2Err[B, string]()
		if 𝗲𝗿𝗿𝟯 != nil {
			return 𝗲𝗿𝗿𝟯
		}
		println(II(𝘃𝗮𝗹𝟮, 𝘃𝗮𝗹𝟯))
		𝘃𝗮𝗹𝟰, 𝘃𝗮𝗹𝟱, 𝘃𝗮𝗹𝟲, 𝗲𝗿𝗿𝟰 := ret3Err[C, string, rune]()
		if 𝗲𝗿𝗿𝟰 != nil {
			return 𝗲𝗿𝗿𝟰
		}
		println(III(𝘃𝗮𝗹𝟰, 𝘃𝗮𝗹𝟱, 𝘃𝗮𝗹𝟲))
	}
	return nil
}