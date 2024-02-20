package rewriter

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"strconv"
	"strings"

	"github.com/goghcrow/go-ansi"
	"github.com/goghcrow/go-loader"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
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

	l := loader.MustNew(
		dir,
		loader.WithLoadDepts(),
		loader.WithLoadTest(),
		loader.WithBuildTag(opt.buildTag),
		loader.WithFileFilter(func(f *loader.File) bool {
			return isTryFile(f.Filename) && imported(f.File, pkgTryPath)
		}),
	)
	r := mkRewriter(*opt, l)
	r.rewriteAllFiles(func(filename string, f *ast.File) {
		filename = replace(filename, srcFileSuffix, ".go")
		filename = replace(filename, testFileSuffix, "_test.go")
		l.WriteFileWithComment(filename, fileComment, f)
	})
}

type rewriter struct {
	opt         option
	l           *loader.Loader
	symCnt      int
	tryFns      map[types.Object]tryFnName
	waitingZero map[ast.Node]func() /*FuncLit|FuncDecl*/
	importRT    bool
}

func mkRewriter(opt option, l *loader.Loader) *rewriter {
	r := &rewriter{
		opt:         opt,
		l:           l,
		tryFns:      map[types.Object]tryFnName{},
		waitingZero: map[ast.Node]func(){},
	}
	for _, fnName := range funcTryNames {
		r.tryFns[l.MustLookup(pkgTryPath+"."+fnName)] = fnName
	}
	return r
}

func (r *rewriter) rewriteAllFiles(printer filePrinter) {
	tryPkg := r.l.LookupPackage(pkgTryPath)
	if tryPkg == nil {
		log.Printf("skipped: missing %s\n", pkgTryPath)
		return
	}

	r.l.VisitAllFiles(func(f *loader.File) {
		log.Printf("write file: %s\n", f.Filename)
		r.editFile(f)               // 1. rewrite try call
		r.editImport(f)             // 2. rewrite import
		f.File.Comments = nil       // 3. delete comments
		printer(f.Filename, f.File) // 4. writeback
	})
}

func (r *rewriter) editFile(f *loader.File) {
	if false {
		r.editFile1(f.File, f.Pkg)
	} else {
		r.editFile2(f.File, f.Pkg)
	}
	r.prependZeroVarDecl()
}

func (r *rewriter) prependZeroVarDecl() {
	for _, do := range r.waitingZero {
		do()
	}
}

func (r *rewriter) editImport(f *loader.File) {
	deleteImport(f.Pkg.Fset, f.File, pkgTryPath)
	if r.importRT {
		astutil.AddNamedImport(f.Pkg.Fset, f.File, ".", pkgRTPath)
	}
}

func (r *rewriter) rewriteTryCall(c *callCtx) {
	c.checkNonConsistentSemantics()

	sig := c.checkSignature()
	c.checkShadowedNil()
	retIds := r.waitToPrependZeroVarDecl(c.outerFun, sig)
	r.doRewriteTryCall(c.c, c.tryFn, c.callsite, c.outerStmts, retIds)
}

func (r *rewriter) waitToPrependZeroVarDecl(fn ast.Node, sig *types.Signature) (retIds []ast.Expr) {
	retCnt := sig.Results().Len()
	if retCnt <= 1 {
		return
	}

	specs := make([]*ast.ValueSpec, 0, retCnt-1)
	retIds = make([]ast.Expr, 0, retCnt)

	fnTy := getField[*ast.FuncType](fn, "Type").Results.List
	fnTy = flattenNamesByType(fnTy, func(n *ast.Ident, ty ast.Expr) *ast.Field {
		return &ast.Field{Names: []*ast.Ident{n}, Type: ty}
	})
	fnBody := getField[*ast.BlockStmt](fn, "Body")

	for i := 0; i < retCnt-1; i++ {
		zero := ast.NewIdent(r.mksym(fmt.Sprintf("%s%d", valZero, i)))
		ty := fnTy[i].Type
		specs = append(specs, &ast.ValueSpec{
			Names: []*ast.Ident{zero},
			Type:  ty,
		})
		retIds = append(retIds, zero)
	}

	r.waitingZero[fn] = func() {
		fnBody.List = prepend[ast.Stmt](fnBody.List, &ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: groupNamesByType[ast.Spec](specs, func(x, y ast.Expr) bool {
					return x == y
				}),
			},
		})
	}
	return
}

func (r *rewriter) doRewriteTryCall(
	c *astutil.Cursor,
	tryFn string, callsite *ast.CallExpr,
	outerStmts *stmts, retIds []ast.Expr,
) {
	vCnt := tryFnRetCnt(tryFn)
	valIds := make([]ast.Expr, 0, vCnt+1)

	switch {
	case vCnt == 0:
		r.deleteCurrent(c)
	case instanceof[*ast.ExprStmt](c.Parent()):
		for i := 0; i < vCnt; i++ {
			valId := ast.NewIdent("_")
			valIds = append(valIds, valId)
		}
		r.deleteCurrent(c)
	case vCnt == 1:
		valId := ast.NewIdent(r.gensym(valIdentPrefix))
		valIds = append(valIds, valId)
		c.Replace(valId)
	default:
		for i := 0; i < vCnt; i++ {
			valId := ast.NewIdent(r.gensym(valIdentPrefix))
			valIds = append(valIds, valId)
		}
		r.importRT = true
		// only one node can be replaced, so conv valIds to tuple call
		c.Replace(&ast.CallExpr{
			Fun:  ast.NewIdent(tupleNames[vCnt]),
			Args: valIds,
		})
	}

	errId := ast.NewIdent(r.gensym(errIdentPrefix))
	outerStmts.insertBefore(
		&ast.AssignStmt{
			Lhs: append(valIds, errId),
			Tok: token.DEFINE,
			Rhs: callsite.Args,
		},
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X:  errId,
				Op: token.NEQ,
				Y:  ast.NewIdent("nil"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: append(retIds, errId),
					},
				},
			},
		},
	)
}

func (r *rewriter) deleteCurrent(c *astutil.Cursor) {
	r.importRT = true
	// DELETE current Node
	c.Replace(&ast.CallExpr{
		Fun: ast.NewIdent(tupleNames[0]),
	})
}

func (r *rewriter) tryCallee(info *types.Info, callsite *ast.CallExpr) (
	callee types.Object, fnName string,
) {
	callee = typeutil.Callee(info, callsite)
	fnName = r.tryFns[callee]
	return
}

func (r *rewriter) gensym(prefix string) string {
	r.symCnt++
	return ansi.Transform("SansSerif-Bold", prefix+strconv.Itoa(r.symCnt))
}
func (r *rewriter) mksym(s string) string {
	return ansi.Transform("SansSerif-Bold", s)
}

func (r *rewriter) resetsym() {
	r.symCnt = 0
}

func (r *rewriter) assert(pkg *pkg, pos positioner, ok bool, format string, a ...any) {
	if !ok {
		loc := pkg.Fset.Position(pos.Pos()).String()
		panic(fmt.Sprintf(format, a...) + " in: " + loc)
	}
}

type (
	filePrinter func(filename string, file *ast.File)
	positioner  interface{ Pos() token.Pos }
	pkg         = packages.Package
	tryFnName   = string
)

type stmts struct {
	idx *int
	xs  *[]ast.Stmt
}

func newStmts(xs *[]ast.Stmt, idx *int) *stmts {
	return &stmts{idx: idx, xs: xs}
}

func (s *stmts) insertBefore(stmts ...ast.Stmt) {
	*s.xs = insertBefore(*s.xs, *s.idx, stmts...)
	*s.idx = *s.idx + len(stmts)
}
