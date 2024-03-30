package rewriter

const fileComment = `//go:build !%s

// Code generated by github.com/goghcrow/go-try DO NOT EDIT.
`

const (
	defaultFileSuffix = "try"
	defaultBuildTag   = "try"
)

const (
	pkgTryPath = "github.com/goghcrow/go-try"
	pkgRTPath  = "github.com/goghcrow/go-try/rt" // 𝙧𝙩

)

const (
	labelPrefix        = "L_"
	valIdentPrefix     = "val"
	errIdentPrefix     = "err"
	valZeroIdentPrefix = "zero"
	postIdentPrefix    = "post"
)

var (
	tryFnNames     = []string{"Try0", "Try", "Try2", "Try3"}
	rtTupleFnNames = []string{"Ø", "Ƭ𝟭", "Ƭ2", "Ƭ3", "Ƭ4", "Ƭ5", "Ƭ6", "Ƭ7", "Ƭ8", "Ƭ9"}
	rtErrorTyName  = "E𝗿𝗿𝗼𝗿"
)

func retCntOfTryFn(tryFn string) int {
	for i, name := range tryFnNames {
		if tryFn == name {
			return i + 1 /*err*/
		}
	}
	panic("illegal try func: " + tryFn)
}
