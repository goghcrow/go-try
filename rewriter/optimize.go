package rewriter

import (
	"go/ast"
	"go/token"
	"go/types"
	"log"

	"github.com/goghcrow/go-loader"
	"github.com/goghcrow/go-matcher"
	"github.com/goghcrow/go-matcher/combinator"
)

type optimizer struct {
	opt   option
	l     *loader.Loader
	m     *matcher.Matcher
	rtFns []types.Object
}

func mkOptimizer(opt option, l *loader.Loader) *optimizer {
	return &optimizer{
		opt: opt,
		l:   l,
		m:   matcher.New(),
		rtFns: mapto(tupleNames, func(i int, n string) types.Object {
			return l.MustLookup(pkgRTPath + "." + n)
		}),
	}
}

func (r *optimizer) optimizeAllFiles(printer filePrinter) {
	rtPkg := r.l.LookupPackage(pkgRTPath)
	if rtPkg == nil {
		log.Printf("skipped: missing %s\n", pkgRTPath)
		return
	}

	r.l.VisitAllFiles(func(f *loader.File) {
		log.Printf("optimize file: %s\n", f.Filename)
		r.clearEmptyStmt(f)
		r.unwrapTuple(f)
		r.clearImport(f)
		printer(f.Filename, f)
	})
}

func (r *optimizer) clearEmptyStmt(f *loader.File) {
	Ø := &ast.ExprStmt{
		X: combinator.CalleeOf(r.m, func(ctx *combinator.MatchCtx, o types.Object) bool {
			return o == r.rtFns[0]
		}),
	}
	r.m.Match(f.Pkg, Ø, f.File, func(c *matcher.Cursor, ctx *matcher.MatchCtx) {
		c.Delete()
	})
}

func (r *optimizer) unwrapTuple(f *loader.File) {
	// iV, bV := II(𝘃𝗮𝗹𝟮, 𝘃𝗮𝗹𝟯)
	// iV, bV, sV := III(𝘃𝗮𝗹𝟮, 𝘃𝗮𝗹𝟯, 𝘃𝗮𝗹𝟰)
	var ident *ast.Ident // wildcard ident
	assignOrDef := matcher.MkPattern[matcher.TokenPattern](r.m, func(n ast.Node, ctx *matcher.MatchCtx) bool {
		var tok = token.Token(n.(matcher.TokenNode))
		return tok == token.ASSIGN || tok == token.DEFINE
	})
	tupleAssign := combinator.OrEx[matcher.NodePattern](r.m,
		&ast.AssignStmt{
			Lhs: []ast.Expr{ident, ident},
			Tok: assignOrDef,
			Rhs: []ast.Expr{
				combinator.CalleeOf(r.m, func(ctx *combinator.MatchCtx, f types.Object) bool {
					return f == r.rtFns[2]
				}),
			},
		},
		&ast.AssignStmt{
			Lhs: []ast.Expr{ident, ident, ident},
			Tok: assignOrDef,
			Rhs: []ast.Expr{
				combinator.CalleeOf(r.m, func(ctx *combinator.MatchCtx, f types.Object) bool {
					return f == r.rtFns[3]
				}),
			},
		},
	)
	r.m.Match(f.Pkg, tupleAssign, f.File, func(c *matcher.Cursor, ctx *matcher.MatchCtx) {
		assign := c.Node().(*ast.AssignStmt)
		c.Replace(&ast.AssignStmt{
			Lhs:    assign.Lhs,
			TokPos: assign.TokPos,
			Tok:    assign.Tok,
			Rhs:    assign.Rhs[0].(*ast.CallExpr).Args,
		})
	})
}

func (r *optimizer) clearImport(f *loader.File) {
	rtCall := combinator.CalleeOf(r.m, func(ctx *combinator.MatchCtx, f types.Object) bool {
		return first(r.rtFns, func(rtF types.Object) bool {
			return rtF == f
		}) != nil
	})
	if !r.m.Matched(f.Pkg, rtCall, f.File) {
		deleteImport(f.Pkg.Fset, f.File, pkgRTPath)
	}
}
