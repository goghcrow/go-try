//go:build !try

// Code generated by github.com/goghcrow/go-try DO NOT EDIT.
package test

func try_in_if_init_or_if_cond() error {
	𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := func1[int, bool](1)
	if 𝗲𝗿𝗿𝟭 != nil {
		return 𝗲𝗿𝗿𝟭
	}
	if 𝘃𝗮𝗹𝟭 {
	} else if false {
	} else {
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := func1[int, bool](2)
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
		if a := 𝘃𝗮𝗹𝟮; a {
		} else {
			𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟯 := func1[int, bool](3)
			if 𝗲𝗿𝗿𝟯 != nil {
				return 𝗲𝗿𝗿𝟯
			}
			if 𝘃𝗮𝗹𝟯 {
			} else if true {
			}
		}
	}
	return nil
}