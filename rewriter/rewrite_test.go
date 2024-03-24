package rewriter

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
)

func TestPlayground(t *testing.T) {
	Rewrite("./test", WithFileSuffix("dbg"))
}

func TestRewrite(t *testing.T) {
	var (
		srcFileSuffix  = fmt.Sprintf("_%s.go", defaultFileSuffix)
		testFileSuffix = fmt.Sprintf("_%s_test.go", defaultFileSuffix)
		endsWith       = strings.HasSuffix
		replace        = strings.ReplaceAll
		isTestFile     = func(filename string) bool {
			return endsWith(filename, srcFileSuffix) || endsWith(filename, testFileSuffix)
		}
		diffFilenames = func(filename string) (got string, want string) {
			got = replace(replace(filename, srcFileSuffix, ".go"), testFileSuffix, "_test.go")
			want = replace(replace(filename, srcFileSuffix, ".want"), testFileSuffix, "_test.want")
			return
		}
	)

	dir := "test"
	Rewrite(dir)

	xs, _ := os.ReadDir(dir)
	for _, x := range xs {
		if x.IsDir() || !isTestFile(x.Name()) {
			continue
		}
		got, want := diffFilenames(x.Name())
		output, err := os.ReadFile(path.Join(dir, got))
		if err != nil {
			t.Fatal(err)
		}
		expect, err := os.ReadFile(path.Join(dir, want))
		if err != nil {
			t.Fatal(err)
		}
		if string(output) != string(expect) {
			t.Fatalf(x.Name())
		}
	}
}
