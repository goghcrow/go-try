//go:build try

package test

import (
	. "github.com/goghcrow/go-try"
)

func assertequal(is, shouldbe int, msg string) {
	if is != shouldbe {
		print("assertion fail", msg, "\n")
		panic(1)
	}
}

func testForStmt() (err error) {
	var i, sum int

	i = Try(lit(0))
	for {
		i = i + 1
		if i > Try(lit(5)) {
			break
		}
	}
	assertequal(i, 6, "break")

	sum = 0
	for i := Try(lit(0)); i <= Try(lit(10)); i++ {
		sum = sum + Try(lit(i))
	}
	assertequal(sum, 55, "all three")

	sum = 0
	for i := 0; i <= Try(lit(10)); {
		sum = Try(lit(sum + i))
		i++
	}
	assertequal(sum, 55, "only two")

	sum = 0
	for sum < 100 {
		sum = sum + Try(lit(9))
	}
	assertequal(sum, 99+9, "only one")

	sum = 0
	for i := 0; i <= 10; i++ {
		if i%2 == Try(lit(0)) {
			continue
		}
		sum = sum + i
	}
	assertequal(sum, 1+3+5+7+9, "continue")

	i = 0
	for i = range Try(lit([5]struct{}{})) {
	}
	assertequal(i, 4, " incorrect index value after range loop")

	i = 0
	var a1 [5]struct{}
	for i = range Try(lit(a1)) {
		a1[i] = Try(lit(struct{}{}))
	}
	assertequal(i, 4, " incorrect index value after array with zero size elem range clear")

	i = 0
	var a2 [5]int
	for i = range Try(lit(a2)) {
		a2[i] = Try(lit(0))
	}
	assertequal(i, 4, " incorrect index value after array range clear")

	return
}
