package rewriter

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
)

// https://stackoverflow.com/questions/14249217/how-do-i-know-im-running-within-go-test
var runningWithGoTest = flag.Lookup("test.v") != nil ||
	strings.HasSuffix(os.Args[0], ".test")

func mkDir(dir string) string {
	dir, err := filepath.Abs(dir)
	panicIf(err)
	err = os.MkdirAll(dir, os.ModePerm)
	panicIf(err)
	return dir
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func assert(ok bool) {
	if !ok {
		panic("illegal state")
	}
}
