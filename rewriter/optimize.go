package rewriter

import (
	"go/ast"
	"go/token"
	"go/types"
	"log"

	"github.com/goghcrow/go-loader"
	"github.com/goghcrow/go-matcher"
	"github.com/goghcrow/go-matcher/combinator"
	"github.com/goghcrow/go-try/rewriter/helper"
)

type optimizer struct {
	m          *matcher.Matcher
	rtTupleFns []types.Object
	rtErrorTy  types.Object
}

func optimize(l *loader.Loader, printer filePrinter) {
	o := &optimizer{
		m: matcher.New(),
		rtTupleFns: sliceMap(rtTupleFnNames, func(i int, n string) types.Object {
			return l.Lookup(pkgRTPath + "." + n)
		}),
		rtErrorTy: l.Lookup(pkgRTPath + "." + rtErrorTyName),
	}
	rtPkg := l.LookupPackage(pkgRTPath)
	if rtPkg == nil {
		log.Printf("skipped: missing %s\n", pkgRTPath)
		return
	}
	l.VisitAllFiles(func(f *loader.File) {
		log.Printf("optimize file: %s\n", f.Filename)
		o.clearEmptyStmt(f)
		o.unwrapTuple(f)
		o.unwrapTupleAssign(f)
		o.mergeBlock(f)
		o.clearImport(f)
		printer(f.Filename, f)
	})
}

// remove?
func (r *optimizer) unwrapTuple(f *loader.File) {
	ptn := combinator.CalleeOf(r.m, func(ctx *combinator.MatchCtx, f types.Object) bool {
		n := ctx.Stack[0].(*ast.CallExpr)
		if _, ok := ctx.Stack[1].(*ast.ExprStmt); !ok {
			for _, it := range r.rtTupleFns[1:] {
				if it != nil && f == it {
					return n.Ellipsis == token.NoPos && len(n.Args) == 1
				}
			}
		}
		return false
	})
	r.m.Match(f.Pkg, ptn, f.File, func(c *matcher.Cursor, ctx *matcher.MatchCtx) {
		c.Replace(c.Node().(*ast.CallExpr).Args[0])
	})
}

func (r *optimizer) clearEmptyStmt(f *loader.File) {
	if r.rtTupleFns[0] == nil {
		return
	}
	type (
		cur = matcher.Cursor
		ctx = matcher.MatchCtx
	)
	// 可以手写一个 pass 移除
	Ø := &ast.ExprStmt{X: combinator.CalleeOf(r.m, func(ctx *combinator.MatchCtx, o types.Object) bool {
		return o == r.rtTupleFns[0]
	})}
	r.m.Match(f.Pkg, Ø, f.File, func(c *cur, ctx *ctx) {
		if c.Index() >= 0 {
			c.Delete()
		}
	})
	r.m.Match(f.Pkg, &ast.IfStmt{Init: Ø}, f.File, func(c *cur, ctx *ctx) {
		c.Node().(*ast.IfStmt).Init = nil
	})
	r.m.Match(f.Pkg, &ast.IfStmt{Else: Ø}, f.File, func(c *cur, ctx *ctx) {
		c.Node().(*ast.IfStmt).Else = nil
	})
	r.m.Match(f.Pkg, &ast.SwitchStmt{Init: Ø}, f.File, func(c *cur, ctx *ctx) {
		c.Node().(*ast.SwitchStmt).Init = nil
	})
	r.m.Match(f.Pkg, &ast.TypeSwitchStmt{Init: Ø}, f.File, func(c *cur, ctx *ctx) {
		c.Node().(*ast.TypeSwitchStmt).Init = nil
	})
	r.m.Match(f.Pkg, &ast.ForStmt{Init: Ø}, f.File, func(c *cur, ctx *ctx) {
		c.Node().(*ast.ForStmt).Init = nil
	})
	r.m.Match(f.Pkg, &ast.ForStmt{Post: Ø}, f.File, func(c *cur, ctx *ctx) {
		c.Node().(*ast.ForStmt).Post = nil
	})
}

func (r *optimizer) unwrapTupleAssign(f *loader.File) {
	// iV := T1(𝘃𝗮𝗹𝟮) // unwrapTuple 已经处理
	// iV, bV := T2(𝘃𝗮𝗹𝟮, 𝘃𝗮𝗹𝟯)
	// iV, bV, sV := T3(𝘃𝗮𝗹𝟮, 𝘃𝗮𝗹𝟯, 𝘃𝗮𝗹𝟰)
	// ...

	assignOrDef := matcher.MkPattern[matcher.TokenPattern](r.m, func(n ast.Node, ctx *matcher.MatchCtx) bool {
		var tok = token.Token(n.(matcher.TokenNode))
		return tok == token.ASSIGN || tok == token.DEFINE
	})
	tupleAssign := matcher.MkPattern[matcher.NodePattern](r.m, func(n ast.Node, ctx *matcher.MatchCtx) bool {
		return false
	})
	for i := 2; i < len(r.rtTupleFns); i++ {
		if r.rtTupleFns[i] == nil {
			continue
		}
		j := i
		tupleAssign = combinator.OrEx[matcher.NodePattern](r.m,
			tupleAssign,
			&ast.AssignStmt{
				Lhs: make([]ast.Expr, j), // wildcard ident
				Tok: assignOrDef,
				Rhs: []ast.Expr{
					combinator.CalleeOf(r.m, func(ctx *combinator.MatchCtx, f types.Object) bool {
						return f == r.rtTupleFns[j]
					}),
				},
			},
		)
	}
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

func (r *optimizer) mergeBlock(f *loader.File) {
	// for preRewriteLabel + preRewriteSwitch
	r.m.Match(f.Pkg, &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.BlockStmt{},
		},
	}, f.File, func(c *matcher.Cursor, ctx *matcher.MatchCtx) {
		block := c.Node().(*ast.BlockStmt)
		c.Replace(block.List[0])
	})
}

func (r *optimizer) clearImport(f *loader.File) {
	rtPtn := combinator.OrEx[matcher.NodePattern](r.m,
		combinator.CalleeOf(r.m, func(ctx *combinator.MatchCtx, f types.Object) bool {
			return findFirst(r.rtTupleFns, func(rtF types.Object) bool { return rtF == f }) != nil
		}),
		&ast.Field{
			Type: combinator.AndEx[matcher.ExprPattern](r.m, &ast.Ident{
				Name: rtErrorTyName,
			}, combinator.ObjectOf(r.m, func(ctx *combinator.MatchCtx, obj types.Object) bool {
				return r.rtErrorTy == obj
			})),
		},
	)
	if !r.m.Matched(f.Pkg, rtPtn, f.File) {
		helper.DeleteImport(f.Pkg.Fset, f.File, pkgRTPath)
	}
}
