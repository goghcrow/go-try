//go:build !try

// Code generated by github.com/goghcrow/go-try DO NOT EDIT.
package test

func try_underlying_fun() error {
	{
		type F func() int
		var f F
		𝘃𝗮𝗹𝟭 := f()
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟭 := ret1Err[int]()
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		_ = 𝘃𝗮𝗹𝟭 + 𝘃𝗮𝗹𝟮
	}
	{
		type ErrF func() (int, error)
		var errF ErrF
		𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟮 := errF()
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
		_ = 𝘃𝗮𝗹𝟯
	}
	return nil
}
