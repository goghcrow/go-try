//go:build !try

// Code generated by github.com/goghcrow/go-try DO NOT EDIT.
package test

func emptyStmt() error {
	𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟮 := ret1Err[int]()
	if 𝗲𝗿𝗿𝟮 != nil {
		return 𝗲𝗿𝗿𝟮
	}
	switch 𝘃𝗮𝗹𝟭 {
	case 0:
		𝗲𝗿𝗿𝟯 := ret0Err()
		if 𝗲𝗿𝗿𝟯 != nil {
			return 𝗲𝗿𝗿𝟯
		}

	}
	return nil
}