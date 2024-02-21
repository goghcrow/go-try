package rewriter

import (
	"fmt"
	"strings"

	"github.com/goghcrow/go-loader"
)

type (
	Option func(*option)
	option struct {
		fileSuffix string
		buildTag   string
	}
)

func WithFileSuffix(s string) Option { return func(opt *option) { opt.fileSuffix = s } }
func WithBuildTag(s string) Option   { return func(opt *option) { opt.buildTag = s } }

func Rewrite(dir string, opts ...Option) {
	opt := &option{
		fileSuffix: defaultFileSuffix,
		buildTag:   defaultBuildTag,
	}
	for _, o := range opts {
		o(opt)
	}

	var (
		endsWith       = strings.HasSuffix
		replace        = strings.ReplaceAll
		srcFileSuffix  = fmt.Sprintf("_%s.go", opt.fileSuffix)
		testFileSuffix = fmt.Sprintf("_%s_test.go", opt.fileSuffix)
		isTryFile      = func(filename string) bool {
			return endsWith(filename, srcFileSuffix) || endsWith(filename, testFileSuffix)
		}
	)

	mkRewriter(*opt, loader.MustNew(
		dir,
		loader.WithLoadDepts(),
		loader.WithLoadTest(),
		loader.WithBuildTag(opt.buildTag),
		loader.WithFileFilter(func(f *loader.File) bool {
			return isTryFile(f.Filename) && imported(f.File, pkgTryPath)
		}),
	)).rewriteAllFiles(func(filename string, f *loader.File) {
		filename = replace(filename, srcFileSuffix, ".go")
		filename = replace(filename, testFileSuffix, "_test.go")
		comment := fmt.Sprintf(fileComment, opt.buildTag)
		f.WriteWithComment(filename, comment)
	})

	mkOptimizer(*opt, loader.MustNew(
		dir,
		loader.WithLoadDepts(),
		loader.WithLoadTest(),
		loader.WithBuildTag("!"+opt.buildTag),
		loader.WithFileFilter(func(f *loader.File) bool {
			return imported(f.File, pkgRTPath)
		}),
	)).optimizeAllFiles(func(filename string, f *loader.File) {
		f.Write(filename)
	})
}
