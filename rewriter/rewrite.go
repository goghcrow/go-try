package rewriter

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"os"
	"strconv"
	"strings"

	matcher "github.com/goghcrow/go-ast-matcher"
	"github.com/goghcrow/go-ast-matcher/imports"
	"golang.org/x/tools/go/ast/astutil"
)

const (
	fileSuffix = "try"
	buildTag   = "try"

	pkgTryPath  = "github.com/goghcrow/go-try"
	funcTryName = "Try"

	valIdentPrefix = "𝘃𝗮𝗹" // 𝐯𝐚𝐥 𝕧𝕒𝕝 𝒗𝒂𝒍 𝘃𝗮𝗹 𝙫𝙖𝙡 𝘷𝘢𝘭 𝚟𝚊𝚕 𝗏𝖺𝗅 ᴠᴀʟ
	errIdentPrefix = "𝐞𝐫𝐫" // 𝘦𝘳𝘳 𝓮𝓻𝓻 𝗲𝗿𝗿 𝐞𝐫𝐫 𝕖𝕣𝕣 𝔢𝔯𝔯 𝖊𝖗𝖗 𝒆𝒓𝒓 𝚎𝚛𝚛 𝖾𝗋𝗋

	fileComment = `//go:build !try

// Code generated by github.com/goghcrow/go-try DO NOT EDIT.
`
)

var (
	srcFileSuffix  = fmt.Sprintf("_%s.go", fileSuffix)
	testFileSuffix = fmt.Sprintf("_%s_test.go", fileSuffix)
)

type FilePrinter func(filename string, file *ast.File)

func Rewrite(dir string) {
	tmpOutputDir := mkDir(dir + "_tmp")
	if !runningWithGoTest {
		//goland:noinspection GoUnhandledErrorResult
		defer os.RemoveAll(tmpOutputDir)
	}

	r := mkRewriter(matcher.NewMatcher(
		dir,
		matcher.PatternAll,
		matcher.WithLoadDepts(),
		matcher.WithLoadTest(),
		matcher.WithBuildTag(buildTag),
		matcher.WithFileFilter(func(filename string, file *ast.File) bool {
			return strings.HasSuffix(filename, srcFileSuffix) || strings.HasSuffix(filename, testFileSuffix)
		}),
	))
	r.rewriteAllFiles(func(filename string, file *ast.File) {
		filename = strings.ReplaceAll(filename, srcFileSuffix, ".go")
		filename = strings.ReplaceAll(filename, testFileSuffix, "_test.go")
		filename = strings.ReplaceAll(filename, dir, tmpOutputDir)
		r.m.WriteFileWithComment(filename, fileComment)
	})
}

type rewriter struct {
	m       *matcher.Matcher
	symCnt  int
	tryFunc types.Object
}

func mkRewriter(m *matcher.Matcher) *rewriter {
	return &rewriter{
		m:       m,
		tryFunc: m.MustLookup(pkgTryPath + "." + funcTryName),
	}
}

func (r *rewriter) gensym(prefix string) string {
	r.symCnt++
	return prefix + strconv.Itoa(r.symCnt)
}

type positioner interface{ Pos() token.Pos }

func (r *rewriter) assert(ok bool, pos positioner, format string, a ...any) {
	if !ok {
		loc := r.m.FSet.Position(pos.Pos()).String()
		panic(fmt.Sprintf(format, a...) + " in: " + loc)
	}
}

func (r *rewriter) rewriteAllFiles(printer FilePrinter) {
	tryPkg := r.m.All[pkgTryPath]
	if tryPkg == nil {
		log.Printf("skip rewrite: no import %s\n", pkgTryPath)
		return
	}

	r.m.VisitAllFiles(func(m *matcher.Matcher, f *ast.File) {
		if !imports.Uses(m, f, tryPkg.Types) {
			log.Printf("skip file: %s\n", r.m.Filename)
			return
		}

		// r.editFile(f)
		r.editFile_(f) // todo

		log.Printf("write file: %s\n", r.m.Filename)
		r.removeImport(f)
		r.clearComments(f)
		printer(r.m.Filename, f)
	})
}

func (r *rewriter) removeImport(f *ast.File) {
	for _, decl := range f.Decls {
		d, ok := decl.(*ast.GenDecl)
		if !ok || d.Tok != token.IMPORT {
			continue
		}
		specs := make([]ast.Spec, 0, len(d.Specs)-1)
		for _, spec := range d.Specs {
			s := spec.(*ast.ImportSpec)
			path := imports.SpecPath(s)
			if path != pkgTryPath {
				namePath := imports.Fmt(s)
				s := &ast.ImportSpec{Path: &ast.BasicLit{Value: namePath, Kind: d.Tok}}
				specs = append(specs, s)
			}
		}
		d.Specs = specs
	}
}

func (r *rewriter) clearComments(f *ast.File) {
	f.Comments = nil
}

func (r *rewriter) editFile(f *ast.File) {
	cache := map[ast.Node]bool{}
	ptn := matcher.FuncCallee(r.m, pkgTryPath, funcTryName)
	r.m.MatchNode(ptn, f, func(m *matcher.Matcher, c *astutil.Cursor, stack []ast.Node, binds matcher.Binds) {
		n := c.Node()
		if cache[n] {
			return
		}
		cache[n] = true

		call := n.(*ast.CallExpr)

		stk := nodeStack(stack)

		inFun := stk.nearestFunc()
		r.assert(inFun != nil, call, "Try must be in a tryable fun(...) (T, error)")

		s := stk.nearestScope(r.m)
		r.assert(s != nil, call, "missing scope")

		replace := r.rewriteTryCall(call, inFun, stk.nearestStmt(), s)
		c.Replace(replace)
	})
}

func (r *rewriter) checkSignature(callTry *ast.CallExpr, fn ast.Node) {
	typeOfFunc := func(n ast.Node) *types.Signature {
		switch n := n.(type) {
		case *ast.FuncLit:
			sig, _ := r.m.TypeOf(n).(*types.Signature)
			return sig
		case *ast.FuncDecl:
			sig, _ := r.m.TypeOf(n.Name).(*types.Signature)
			return sig
		default:
			return nil
		}
	}

	sig := typeOfFunc(fn)
	r.assert(sig != nil && sig.Results().Len() == 2, callTry,
		"Try must be in a tryable fun(...) (T, error)")

	var fst, snd types.Type
	switch len(callTry.Args) {
	case 1: // a,b := x
		tup, _ := r.m.TypeOf(callTry.Args[0]).(*types.Tuple)
		r.assert(tup.Len() == 2, callTry, "invalid args, expect 2")
		fst, snd = tup.At(0).Type(), tup.At(1).Type()
	case 2:
		fst, snd = r.m.TypeOf(callTry.Args[0]), r.m.TypeOf(callTry.Args[1])
	default:
		r.assert(false, callTry, "invalid args, expect 1 or 2")
	}

	retFst := sig.Results().At(0).Type()
	retSnd := sig.Results().At(1).Type()
	r.assert(types.AssignableTo(fst, retFst), callTry,
		"type mismatch, Try(?, ) expect %v but %v", retFst, fst)
	r.assert(types.AssignableTo(snd, retSnd), callTry,
		"type mismatch, Try(, ?) expect %v but %v", retSnd, snd)
	r.assert(types.AssignableTo(snd, types.Universe.Lookup("error").Type()), callTry,
		"type mismatch, Try(, ?) expect error but %v", snd)
}

func (r *rewriter) checkNilShadowed(s *types.Scope) {
	nilObj := s.Lookup("nil")
	if nilObj == nil {
		_, nilObj = s.LookupParent("nil", token.NoPos)
	}
	assert(nilObj != nil)
	r.assert(nilObj.Type() == types.Typ[types.UntypedNil], nilObj,
		"nil shadowed in fun scope, please rename it")
}

// v, e := ..
// if e != nil { return v, e }
// v
func (r *rewriter) rewriteTryCall(
	call *ast.CallExpr,
	outerFun ast.Node,
	outerStmts *stmts,
	outerScope *types.Scope,
) ast.Node {
	r.checkSignature(call, outerFun)
	r.checkNilShadowed(outerScope)

	valId := ast.NewIdent(r.gensym(valIdentPrefix))
	errId := ast.NewIdent(r.gensym(errIdentPrefix))
	lhs := []ast.Expr{valId, errId}

	outerStmts.insertAfter(
		&ast.AssignStmt{
			Lhs: lhs,
			Tok: token.DEFINE,
			Rhs: call.Args,
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
						Results: lhs,
					},
				},
			},
		},
	)

	return valId
}
