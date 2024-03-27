package rewriter

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

	"github.com/goghcrow/go-ansi"
	"github.com/goghcrow/go-loader"
	"github.com/goghcrow/go-matcher"
	"github.com/goghcrow/go-matcher/combinator"
	"github.com/goghcrow/go-try/rewriter/helper"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/types/typeutil"
)

type (
	mctx   = *matcher.MatchCtx
	cursor = *matcher.Cursor
	fnName = string
	fnNode ast.Node // *ast.FuncLit | *ast.FuncDecl
)

type enclosingFn struct {
	sig   *types.Signature
	errTy types.Type
	errId *ast.Ident
}

type fileRewriter struct {
	pkg loader.Pkg
	f   *ast.File

	errTy types.Type

	tryFns   map[types.Object]fnName // try 函数对象 => 名称
	tryNodes map[ast.Node]bool       // 包含 try 调用的节点 set

	fnSig   map[fnNode]*enclosingFn     // 函数 签名, 返回 error 类型和变量名称
	fnZero  map[fnNode][]*ast.ValueSpec // 函数 zero value 返回值声明
	fnSyms  map[fnNode]symTbl           // 函数维度的 symbol 生成的计数表
	jmpTbls map[fnNode]helper.JumpTable // 函数内控制流表

	importRT bool // 是否需要导入 rt 包

	m        *matcher.Matcher
	tChecker *helper.TerminationChecker
}

func rewriteFile(tryFns map[types.Object]fnName, pkg loader.Pkg, f *ast.File) *ast.File {
	r := &fileRewriter{
		pkg:      pkg,
		f:        f,
		errTy:    types.Universe.Lookup("error").Type(),
		tryFns:   tryFns,
		tryNodes: map[ast.Node]bool{},
		fnSig:    map[fnNode]*enclosingFn{},
		fnZero:   map[fnNode][]*ast.ValueSpec{},
		fnSyms:   map[fnNode]symTbl{},
		jmpTbls:  map[fnNode]helper.JumpTable{},
		m:        matcher.New(),
	}
	r.tChecker = helper.NewTerminationChecker(r.collectPanicCalls())
	r.tryNodes = r.collectTryNodes()

	// 注意顺序:
	// 必须先改写 labeled, 再改写其他
	// 必须先改写 switch, 再改写 if, 因为 switch 会生成新的 if
	// 必须先改写 switch, 再改写 ||, 因为 switch 会生成新的 ||
	r.preRewriteLabeled()
	r.preRewriteSwitch()
	r.preRewriteIf()
	r.preRewriteFor()
	// typeSwitch 和 range 不会被 Try 影响
	// 		typeSwitch 的 case 必须是类型, 不会包含 try 调用
	// 		range 没 post 和 cond, 不会包含 try 调用

	r.rewriteFile()

	r.postRewriteImport()

	return r.f
}

func (r *fileRewriter) rewriteFile() {
	ctx := newWalkCtx(r.f)
	var xs []ast.Stmt
	r.f.Decls, xs = r.rewriteDeclList(ctx, r.f.Decls)
	// checkTryCall 已经检查过不在 func 内部的 try
	// 所以 file top level decl 不可能产生额外的 stmt
	assert(len(xs) == 0)
	return
}

func (r *fileRewriter) postRewriteImport() {
	helper.DeleteImport(r.pkg.Fset, r.f, pkgTryPath)
	if r.importRT {
		astutil.AddNamedImport(r.pkg.Fset, r.f, ".", pkgRTPath)
	}
}

func (r *fileRewriter) collectPanicCalls() map[*ast.CallExpr]bool {
	xs := map[*ast.CallExpr]bool{}
	panicCall := combinator.BuiltinCallee(r.m, "panic")
	r.match(panicCall, func(c cursor, ctx mctx) {
		xs[c.Node().(*ast.CallExpr)] = true
	})
	return xs
}

func (r *fileRewriter) collectTryNodes() map[ast.Node]bool {
	xs := map[ast.Node]bool{}
	var tryFnCalls = combinator.FuncCalleeOf(r.m, func(_ *combinator.MatchCtx, obj *types.Func) bool {
		return r.tryFns[obj] != ""
	})
	r.match(tryFnCalls, func(c cursor, ctx mctx) {
		// 跨越函数边界从叶子节点标记(传染)整个包含 try 调用的路径
		// 用来 walk 时候快速判断子树是否需要改写, 从而快速返回
		// 	即 if !ctx.try {  return n, nil }
		for _, n := range ctx.Stack {
			xs[n] = true
		}

		// 在尚未改写之前排查 try call 的 scope nil 是否被 shadow
		// 否则一旦 ast 被改写, scope 信息会不准确
		r.checkShadowedNil(ctx)
	})
	return xs
}

func (r *fileRewriter) checkShadowedNil(tryCallSiteCtx mctx) {
	stk := tryCallSiteCtx.Stack

	var s *types.Scope
	for _, n := range stk {
		s = r.pkg.ScopeFor(n)
		if s != nil {
			break
		}
	}
	assert(s != nil)
	nilObj := s.Lookup("nil")
	if nilObj == nil {
		_, nilObj = s.LookupParent("nil", token.NoPos)
	}
	assert(nilObj != nil)
	nilIsNil := nilObj.Type() == types.Typ[types.UntypedNil]
	tryPos := stk[0].Pos()
	nilPos := nilObj.Pos()
	if nilPos == token.NoPos {
		assert(nilIsNil)
		return
	}
	// 要么在 try 之后 shadow 了 nil
	// 要么 nil 没有被 shadow
	r.assert(nilObj, nilPos > tryPos || nilIsNil, "nil shadowed, please rename it")
}

func (r *fileRewriter) match(ptn ast.Node, f func(c cursor, ctx mctx)) {
	r.m.Match(r.pkg.Package, ptn, r.f, f)
	// cache := map[ast.Node]bool{}
	// r.m.Match(r.pkg.Package, ptn, r.f, func(c cursor, ctx mctx) {
	// 	n := c.Node()
	// 	if cache[n] {
	// 		return
	// 	}
	// 	cache[n] = true
	// 	f(c, ctx)
	// })
}

func (r *fileRewriter) assert(pos loader.Positioner, ok bool, format string, a ...any) {
	if !ok {
		panic(fmt.Sprintf(format, a...) + " in: " + r.pkg.ShowPos(pos))
	}
}

func (r *fileRewriter) jmpTbl(f fnNode) helper.JumpTable {
	tbl := r.jmpTbls[f]
	if tbl == nil {
		_, body := unpackFunc(f)
		tbl = helper.Target(body, func(pos loader.Positioner, msg string) {
			r.assert(pos, false, msg)
		})
		r.jmpTbls[f] = tbl
	}
	return tbl
}

func (r *fileRewriter) clearJmpTbl(f fnNode) {
	delete(r.jmpTbls, f)
}

func (r *fileRewriter) containsTryCall(n ast.Node, _ mctx) bool {
	return r.tryNodes[n]
}

func (r *fileRewriter) tryCallee(callsite *ast.CallExpr) string {
	return r.tryFns[typeutil.Callee(r.pkg.TypesInfo, callsite)]
}

func (r *fileRewriter) mkSym(f fnNode, s string) *ast.Ident       { return r.symTbl(f).mk(s) }
func (r *fileRewriter) genSym(f fnNode, prefix string) *ast.Ident { return r.symTbl(f).gen(prefix) }
func (r *fileRewriter) genValId(f fnNode) *ast.Ident              { return r.symTbl(f).gen(valIdentPrefix) }
func (r *fileRewriter) genErrId(f fnNode) *ast.Ident              { return r.symTbl(f).gen(errIdentPrefix) }

func (r *fileRewriter) symTbl(f fnNode) symTbl {
	sym := r.fnSyms[f]
	if sym == nil {
		sym = symTbl{}
		r.fnSyms[f] = sym
	}
	return sym
}

type symTbl map[string]int

func (r symTbl) gen(prefix string) *ast.Ident {
	r[prefix]++
	return ast.NewIdent(ansi.Transform("SansSerif-Bold", prefix+strconv.Itoa(r[prefix])))
}

func (r symTbl) mk(s string) *ast.Ident {
	return ast.NewIdent(ansi.Transform("SansSerif-Bold", s))
}
