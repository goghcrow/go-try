package rewriter

const (
	defaultFileSuffix = "try"
	defaultBuildTag   = "try"

	pkgTryPath = "github.com/goghcrow/go-try"
	pkgRTPath  = "github.com/goghcrow/go-try/rt" // 𝙧𝙩

	valIdentPrefix = "val"
	errIdentPrefix = "err"
	valZero        = "zero"

	fileComment = `//go:build !try

// Code generated by github.com/goghcrow/go-try DO NOT EDIT.
`
)

var (
	funcTryNames = []string{"Try0", "Try", "Try2", "Try3"}
	tupleNames   = []string{"Ø", "I", "II", "III"}
)

// tryFnRetCnt returns the number of return values without error of tryFn.
func tryFnRetCnt(tryFn string) int {
	for i, name := range funcTryNames {
		if tryFn == name {
			return i
		}
	}
	panic("illegal tryFn")
}