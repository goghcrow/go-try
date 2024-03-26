package rewriter

import (
	"fmt"
	"go/types"
	"log"
	"strings"

	"github.com/goghcrow/go-loader"
	"github.com/goghcrow/go-try/rewriter/helper"
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
		toOptzFile = map[string]bool{}
	)

	rewrite(loader.MustNew(
		dir,
		loader.WithLoadDepts(),
		loader.WithLoadTest(),
		loader.WithBuildTag(opt.buildTag),
		loader.WithFileFilter(func(f *loader.File) bool {
			return isTryFile(f.Filename) && helper.Imported(f.File, pkgTryPath)
		}),
	), func(filename string, f *loader.File) {
		filename = replace(filename, srcFileSuffix, ".go")
		filename = replace(filename, testFileSuffix, "_test.go")
		comment := fmt.Sprintf(fileComment, opt.buildTag)
		f.WriteWithComment(filename, comment)
		toOptzFile[filename] = true
	})

	optimize(loader.MustNew(
		dir,
		loader.WithLoadDepts(),
		loader.WithLoadTest(),
		loader.WithBuildTag("!"+opt.buildTag),
		loader.WithFileFilter(func(f *loader.File) bool {
			return toOptzFile[f.Filename] && helper.Imported(f.File, pkgRTPath)
		}),
	), func(filename string, f *loader.File) {
		f.Write(filename)
	})
}

type filePrinter func(filename string, file *loader.File)

func rewrite(l *loader.Loader, printer filePrinter) {
	tryPkg := l.LookupPackage(pkgTryPath)
	if tryPkg == nil {
		log.Printf("skipped: missing %s\n", pkgTryPath)
		return
	}

	tryFns := map[types.Object]fnName{}
	for _, n := range funcTryNames {
		tryFns[l.MustLookup(pkgTryPath+"."+n)] = n
	}

	l.VisitAllFiles(func(f *loader.File) {
		log.Printf("write file: %s\n", f.Filename)
		pkg := loader.MkPkg(f.Pkg)
		f.File = rewriteFile(tryFns, pkg, f.File) // 1. rewrite try call
		f.File.Comments = nil                     // 2. delete comments
		f.File.Doc = nil                          // 3. delete pkg doc
		loader.ClearPos(f.File)                   // 4. clear position
		printer(f.Filename, f)                    // 5. writeback
	})
}
