//go:build !try

// Code generated by github.com/goghcrow/go-try DO NOT EDIT.
package test

func assertequal(is, shouldbe int, msg string) {
	if is != shouldbe {
		print("assertion fail", msg, "\n")
		panic(1)
	}
}
func testForStmt() (err error) {
	var i, sum int
	𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := lit(0)
	if 𝗲𝗿𝗿𝟭 != nil {
		err = 𝗲𝗿𝗿𝟭
		return
	}
	i = 𝘃𝗮𝗹𝟭
	for {
		i = i + 1
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := lit(5)
		if 𝗲𝗿𝗿𝟮 != nil {
			err = 𝗲𝗿𝗿𝟮
			return
		}
		if i > 𝘃𝗮𝗹𝟮 {
			break
		}
	}
	assertequal(i, 6, "break")
	sum = 0
	𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟯 := lit(0)
	if 𝗲𝗿𝗿𝟯 != nil {
		err = 𝗲𝗿𝗿𝟯
		return
	}
	for i := 𝘃𝗮𝗹𝟯; ; i++ {
		𝘃𝗮𝗹𝟰, 𝗲𝗿𝗿𝟰 := lit(10)
		if 𝗲𝗿𝗿𝟰 != nil {
			err = 𝗲𝗿𝗿𝟰
			return
		}
		if i > 𝘃𝗮𝗹𝟰 {
			break
		}
		𝘃𝗮𝗹𝟱, 𝗲𝗿𝗿𝟱 := lit(i)
		if 𝗲𝗿𝗿𝟱 != nil {
			err = 𝗲𝗿𝗿𝟱
			return
		}
		sum = sum + 𝘃𝗮𝗹𝟱
	}
	assertequal(sum, 55, "all three")
	sum = 0
	for i := 0; ; {
		𝘃𝗮𝗹𝟲, 𝗲𝗿𝗿𝟲 := lit(10)
		if 𝗲𝗿𝗿𝟲 != nil {
			err = 𝗲𝗿𝗿𝟲
			return
		}
		if i > 𝘃𝗮𝗹𝟲 {
			break
		}
		𝘃𝗮𝗹𝟳, 𝗲𝗿𝗿𝟳 := lit(sum + i)
		if 𝗲𝗿𝗿𝟳 != nil {
			err = 𝗲𝗿𝗿𝟳
			return
		}
		sum = 𝘃𝗮𝗹𝟳
		i++
	}
	assertequal(sum, 55, "only two")
	sum = 0
	for sum < 100 {
		𝘃𝗮𝗹𝟴, 𝗲𝗿𝗿𝟴 := lit(9)
		if 𝗲𝗿𝗿𝟴 != nil {
			err = 𝗲𝗿𝗿𝟴
			return
		}
		sum = sum + 𝘃𝗮𝗹𝟴
	}
	assertequal(sum, 99+9, "only one")
	sum = 0
	for i := 0; i <= 10; i++ {
		𝘃𝗮𝗹𝟵, 𝗲𝗿𝗿𝟵 := lit(0)
		if 𝗲𝗿𝗿𝟵 != nil {
			err = 𝗲𝗿𝗿𝟵
			return
		}
		if i%2 == 𝘃𝗮𝗹𝟵 {
			continue
		}
		sum = sum + i
	}
	assertequal(sum, 1+3+5+7+9, "continue")
	i = 0
	𝘃𝗮𝗹𝟭𝟬, 𝗲𝗿𝗿𝟭𝟬 := lit([5]struct{}{})
	if 𝗲𝗿𝗿𝟭𝟬 != nil {
		err = 𝗲𝗿𝗿𝟭𝟬
		return
	}
	for i = range 𝘃𝗮𝗹𝟭𝟬 {
	}
	assertequal(i, 4, " incorrect index value after range loop")
	i = 0
	var a1 [5]struct{}
	𝘃𝗮𝗹𝟭𝟭, 𝗲𝗿𝗿𝟭𝟭 := lit(a1)
	if 𝗲𝗿𝗿𝟭𝟭 != nil {
		err = 𝗲𝗿𝗿𝟭𝟭
		return
	}
	for i = range 𝘃𝗮𝗹𝟭𝟭 {
		𝘃𝗮𝗹𝟭𝟮, 𝗲𝗿𝗿𝟭𝟮 := lit(struct{}{})
		if 𝗲𝗿𝗿𝟭𝟮 != nil {
			err = 𝗲𝗿𝗿𝟭𝟮
			return
		}
		a1[i] = 𝘃𝗮𝗹𝟭𝟮
	}
	assertequal(i, 4, " incorrect index value after array with zero size elem range clear")
	i = 0
	var a2 [5]int
	𝘃𝗮𝗹𝟭𝟯, 𝗲𝗿𝗿𝟭𝟯 := lit(a2)
	if 𝗲𝗿𝗿𝟭𝟯 != nil {
		err = 𝗲𝗿𝗿𝟭𝟯
		return
	}
	for i = range 𝘃𝗮𝗹𝟭𝟯 {
		𝘃𝗮𝗹𝟭𝟰, 𝗲𝗿𝗿𝟭𝟰 := lit(0)
		if 𝗲𝗿𝗿𝟭𝟰 != nil {
			err = 𝗲𝗿𝗿𝟭𝟰
			return
		}
		a2[i] = 𝘃𝗮𝗹𝟭𝟰
	}
	assertequal(i, 4, " incorrect index value after array range clear")
	return
}
