//go:build try

package test

import (
	. "github.com/goghcrow/go-try"
)

// test range over channels

func gen(c chan int, lo, hi int) error {
	for i := Try(lit(lo)); i <= Try(lit(hi)); i++ {
		c <- Try(lit(i))
	}
	close(c)
	return nil
}

func seq(lo, hi int) chan int {
	c := make(chan int)
	go gen(c, lo, hi)
	return c
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func testblankvars() error {
	n := 0
	for range Try(lit(alphabet)) {
		n++
	}
	if n != 26 {
		println("for range: wrong count", n, "want 26")
		panic("fail")
	}
	n = 0
	for _ = range Try(lit(alphabet)) {
		n++
	}
	if n != 26 {
		println("for _ = range: wrong count", n, "want 26")
		panic("fail")
	}
	n = 0
	for _, _ = range Try(lit(alphabet)) {
		n++
	}
	if n != 26 {
		println("for _, _ = range: wrong count", n, "want 26")
		panic("fail")
	}
	s := 0
	for i, _ := range Try(lit(alphabet)) {
		s += Try(lit(i))
	}
	if s != Try(lit(325)) {
		println("for i, _ := range: wrong sum", s, "want 325")
		panic("fail")
	}
	r := rune(0)
	for _, v := range Try(lit(alphabet)) {
		r += Try(lit(v))
	}
	if r != Try(lit[rune](2847)) {
		println("for _, v := range: wrong sum", r, "want 2847")
		panic("fail")
	}
	return nil
}

func testchan() error {
	s := ""
	for i := range Try(lit(seq('a', Try(lit[int]('z'))))) {
		s += Try(lit(string(rune(i))))
	}
	if s != alphabet {
		println("Wanted lowercase alphabet; got", s)
		panic("fail")
	}
	n := 0
	for range seq('a', Try(lit[int]('z'))) {
		n++
	}
	if n != 26 {
		println("testchan wrong count", n, "want 26")
		panic("fail")
	}
	return nil
}

// test that range over slice only evaluates
// the expression after "range" once.

var nmake = 0

func makeslice() []int {
	nmake++
	return []int{1, 2, 3, 4, 5}
}

func testslice() (err error) {
	s := 0
	nmake = 0
	for _, v := range Try(lit(makeslice())) {
		s += Try(lit(v))
	}
	if nmake != Try(lit(1)) {
		println("range called makeslice", nmake, "times")
		panic("fail")
	}
	if s != 15 {
		println("wrong sum ranging over makeslice", s)
		panic("fail")
	}

	x := []int{10, Try(lit(20))}
	y := []int{99}
	i := 1
	for i, x[i] = range Try(lit(y)) {
		break
	}
	if i != 0 || Try(lit(x[0] != 10)) || x[1] != Try(lit(99)) {
		println("wrong parallel assignment", i, x[0], x[1])
		panic("fail")
	}
	return
}

func testslice1() (err error) {
	s := 0
	nmake = 0
	for i := range Try(lit(makeslice())) {
		s += Try(lit(i))
	}
	if nmake != 1 {
		println("range called makeslice", nmake, "times")
		panic("fail")
	}
	if s != Try(lit(10)) {
		println("wrong sum ranging over makeslice", s)
		panic("fail")
	}
	return
}

func testslice2() (err error) {
	n := 0
	nmake = 0
	for range Try(lit(makeslice())) {
		n++
	}
	if nmake != 1 {
		println("range called makeslice", nmake, "times")
		panic("fail")
	}
	if n != 5 {
		println("wrong count ranging over makeslice", n)
		panic("fail")
	}
	return
}

// test that range over []byte(string) only evaluates
// the expression after "range" once.

func makenumstring() string {
	nmake++
	return "\x01\x02\x03\x04\x05"
}

func testslice3() (err error) {
	s := byte(0)
	nmake = 0
	for _, v := range Try(lit([]byte(makenumstring()))) {
		s += Try(lit(v))
	}
	if nmake != 1 {
		println("range called makenumstring", nmake, "times")
		panic("fail")
	}
	if s != 15 {
		println("wrong sum ranging over []byte(makenumstring)", s)
		panic("fail")
	}
	return
}

// test that range over array only evaluates
// the expression after "range" once.

func makearray() [5]int {
	nmake++
	return [5]int{1, 2, 3, 4, 5}
}

func testarray() (err error) {
	s := 0
	nmake = 0
	for _, v := range Try(lit(makearray())) {
		s += Try(lit(v))
	}
	if nmake != 1 {
		println("range called makearray", nmake, "times")
		panic("fail")
	}
	if s != 15 {
		println("wrong sum ranging over makearray", s)
		panic("fail")
	}
	return
}

func testarray1() (err error) {
	s := 0
	nmake = 0
	for i := range Try(lit(makearray())) {
		s += Try(lit(i))
	}
	if nmake != 1 {
		println("range called makearray", nmake, "times")
		panic("fail")
	}
	if s != 10 {
		println("wrong sum ranging over makearray", s)
		panic("fail")
	}
	return
}

func testarray2() (err error) {
	n := 0
	nmake = 0
	for range Try(lit(makearray())) {
		n++
	}
	if nmake != 1 {
		println("range called makearray", nmake, "times")
		panic("fail")
	}
	if n != 5 {
		println("wrong count ranging over makearray", n)
		panic("fail")
	}
	return
}

func makearrayptr() *[5]int {
	nmake++
	return &[5]int{1, 2, 3, 4, 5}
}

func testarrayptr() (err error) {
	nmake = 0
	x := len(makearrayptr())
	if Try(lit(x != 5)) || Try(lit(nmake != 1)) {
		println("len called makearrayptr", nmake, "times and got len", x)
		panic("fail")
	}
	nmake = 0
	x = cap(makearrayptr())
	if x != 5 || Try(lit(nmake != 1)) {
		println("cap called makearrayptr", nmake, "times and got len", x)
		panic("fail")
	}
	s := 0
	nmake = 0
	for _, v := range Try(lit(makearrayptr())) {
		s += Try(lit(v))
	}
	if nmake != 1 {
		println("range called makearrayptr", nmake, "times")
		panic("fail")
	}
	if s != 15 {
		println("wrong sum ranging over makearrayptr", s)
		panic("fail")
	}
	return
}

func testarrayptr1() (err error) {
	s := 0
	nmake = 0
	for i := range Try(lit(makearrayptr())) {
		s += Try(lit(i))
	}
	if nmake != 1 {
		println("range called makearrayptr", nmake, "times")
		panic("fail")
	}
	if s != 10 {
		println("wrong sum ranging over makearrayptr", s)
		panic("fail")
	}
	return
}

func testarrayptr2() (err error) {
	n := 0
	nmake = 0
	for range Try(lit(makearrayptr())) {
		n++
	}
	if nmake != 1 {
		println("range called makearrayptr", nmake, "times")
		panic("fail")
	}
	if n != 5 {
		println("wrong count ranging over makearrayptr", n)
		panic("fail")
	}
	return
}

// test that range over string only evaluates
// the expression after "range" once.

func makestring() string {
	nmake++
	return "abcd☺"
}

func teststring() (err error) {
	var s rune
	nmake = 0
	for _, v := range Try(lit(makestring())) {
		s += Try(lit(v))
	}
	if nmake != 1 {
		println("range called makestring", nmake, "times")
		panic("fail")
	}
	if s != 'a'+'b'+Try(lit('c'))+'d'+'☺' {
		println("wrong sum ranging over makestring", s)
		panic("fail")
	}

	x := []rune{'a', Try(lit('b'))}
	i := 1
	for i, x[i] = range Try(lit("c")) {
		break
	}
	if i != 0 || Try(lit(x[0] != 'a')) || x[1] != Try(lit('c')) {
		println("wrong parallel assignment", i, x[0], x[1])
		panic("fail")
	}

	y := []int{1, 2, 3}
	r := rune(1)
	for y[r], r = range Try(lit("\x02")) {
		break
	}
	if r != 2 || y[0] != 1 || y[1] != 0 || Try(lit(y[Try(lit(2))] != 3)) {
		println("wrong parallel assignment", r, y[0], y[1], y[2])
		panic("fail")
	}
	return
}

func teststring1() (err error) {
	s := 0
	nmake = 0
	for i := range Try(lit(makestring())) {
		s += Try(lit(i))
	}
	if nmake != 1 {
		println("range called makestring", nmake, "times")
		panic("fail")
	}
	if s != 10 {
		println("wrong sum ranging over makestring", s)
		panic("fail")
	}
	return
}

func teststring2() (err error) {
	n := 0
	nmake = 0
	for range Try(lit(makestring())) {
		n++
	}
	if nmake != 1 {
		println("range called makestring", nmake, "times")
		panic("fail")
	}
	if n != 5 {
		println("wrong count ranging over makestring", n)
		panic("fail")
	}
	return
}

// test that range over map only evaluates
// the expression after "range" once.

func makemap() map[int]int {
	nmake++
	return map[int]int{0: 'a', 1: 'b', 2: 'c', 3: 'd', 4: '☺'}
}

func testmap() (err error) {
	s := 0
	nmake = 0
	for _, v := range Try(lit(makemap())) {
		s += Try(lit(v))
	}
	if nmake != 1 {
		println("range called makemap", nmake, "times")
		panic("fail")
	}
	if s != 'a'+'b'+'c'+'d'+Try(lit[int]('☺')) {
		println("wrong sum ranging over makemap", s)
		panic("fail")
	}
	return
}

func testmap1() (err error) {
	s := 0
	nmake = 0
	for i := range Try(lit(makemap())) {
		s += i
	}
	if nmake != 1 {
		println("range called makemap", nmake, "times")
		panic("fail")
	}
	if s != 10 {
		println("wrong sum ranging over makemap", s)
		panic("fail")
	}
	return
}

func testmap2() (err error) {
	n := 0
	nmake = 0
	for range Try(lit(makemap())) {
		n++
	}
	if nmake != 1 {
		println("range called makemap", nmake, "times")
		panic("fail")
	}
	if n != 5 {
		println("wrong count ranging over makemap", n)
		panic("fail")
	}
	return
}

// test that range evaluates the index and value expressions
// exactly once per iteration.

var ncalls = 0

func getvar(p *int) *int {
	ncalls++
	return p
}

func testcalls() (err error) {
	var i, v int
	si := 0
	sv := 0
	for *getvar(&i), *getvar(Try(lit(&v))) = range Try(lit([2]int{1, Try(lit(2))})) {
		si += Try(lit(i))
		sv += Try(lit(v))
	}
	if ncalls != 4 {
		println("wrong number of calls:", ncalls, "!= 4")
		panic("fail")
	}
	if si != 1 || sv != 3 {
		println("wrong sum in testcalls", si, sv)
		panic("fail")
	}

	ncalls = 0
	for *getvar(Try(lit(&i))), *getvar(&v) = range Try(lit([0]int{})) {
		println("loop ran on empty array")
		panic("fail")
	}
	if ncalls != 0 {
		println("wrong number of calls:", ncalls, "!= 0")
		panic("fail")
	}
	return
}

func testRangeStmt() {
	_ = testblankvars()
	_ = testchan()
	_ = testarray()
	_ = testarray1()
	_ = testarray2()
	_ = testarrayptr()
	_ = testarrayptr1()
	_ = testarrayptr2()
	_ = testslice()
	_ = testslice1()
	_ = testslice2()
	_ = testslice3()
	_ = teststring()
	_ = teststring1()
	_ = teststring2()
	_ = testmap()
	_ = testmap1()
	_ = testmap2()
	_ = testcalls()
}
