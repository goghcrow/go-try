package rewriter

import (
	"go/ast"
	"go/types"

	"github.com/goghcrow/go-matcher"
	"github.com/goghcrow/go-matcher/combinator"
)

func (r *rewriter) editFile1(f *ast.File, pkg *pkg) {
	cache := map[ast.Node]bool{}

	m := matcher.New()
	var ptn *ast.CallExpr = combinator.FuncCalleeOf(m, func(_ *combinator.MatchCtx, obj *types.Func) bool {
		return r.tryFns[obj] != ""
	})
	m.Match(pkg, ptn, f, func(c *matcher.Cursor, ctx *matcher.MatchCtx) {
		n := c.Node()
		if cache[n] {
			return
		}
		cache[n] = true

		callsite := n.(*ast.CallExpr)

		_, tryFn := r.tryCallee(ctx.TypeInfo(), callsite)
		assert(tryFn != "")

		stk := newTravalStack(ctx.Stack, ctx.Names)

		fun := stk.nearestFunc()
		r.assert(pkg, callsite, fun != nil, "Try must be in a tryable fun(...) (T, error)")

		stmts := stk.nearestStmts()
		r.assert(pkg, callsite, stmts != nil, "missing stmts")

		scope := stk.nearestScope(ctx.TypeInfo())
		r.assert(pkg, callsite, scope != nil, "missing scope")

		r.rewriteTryCall(&callCtx{
			c:          c,
			stk:        stk,
			pkg:        pkg,
			tryFn:      tryFn,
			callsite:   callsite,
			outerFun:   fun,
			outerStmts: stmts,
			outerScope: scope,
		})
	})
}

// ↓↓↓↓↓↓↓↓↓↓↓↓↓↓ traval stack ↓↓↓↓↓↓↓↓↓↓↓↓↓↓

type travalNode struct {
	node  ast.Node
	field string // parent field
}

type travalStack []travalNode

func newTravalStack(stack []ast.Node, names []string) travalStack {
	assert(len(stack) == len(names))
	return mapto(stack, func(idx int, n ast.Node) travalNode {
		return travalNode{
			node:  n,
			field: names[idx],
		}
	})
}

func (s travalStack) nearestStmts() *stmts {
	_, ok := s[0].node.(*ast.CallExpr)
	assert(ok)

	indexOf := func(xs []ast.Stmt, x ast.Node) int {
		for idx, stmt := range xs {
			if x == stmt {
				return idx
			}
		}
		panic("illegal state")
	}

	var child ast.Node
	for i, n := range s {
		if i > 0 {
			child = s[i-1].node
		}

		switch n := n.node.(type) {
		case *ast.BlockStmt:
			switch n.List[0].(type) {
			case *ast.CaseClause, *ast.CommClause:
				continue
			}
			idx := indexOf(n.List, child)
			return newStmts(&n.List, &idx)

		case *ast.CaseClause:
			if x, ok := child.(ast.Expr); ok {
				if first(n.List, func(y ast.Expr) bool { return y == x }) != nil {
					continue
				}
			}
			body := getFieldPtr[[]ast.Stmt](n, "Body")
			idx := indexOf(*body, child)
			return newStmts(body, &idx)

		case *ast.CommClause:
			if x, ok := child.(ast.Stmt); ok {
				if n.Comm == x {
					continue
				}
			}
			body := getFieldPtr[[]ast.Stmt](n, "Body")
			idx := indexOf(*body, child)
			return newStmts(body, &idx)

		default:
			continue
		}
	}
	return nil
}

func (s travalStack) nearestScope(info *types.Info) *types.Scope {
	for _, n := range s {
		// The following node types may appear in Scopes:
		//     *ast.File
		//     *ast.FuncType
		//     *ast.TypeSpec
		//     *ast.BlockStmt
		//     *ast.IfStmt
		//     *ast.SwitchStmt
		//     *ast.TypeSwitchStmt
		//     *ast.CaseClause
		//     *ast.CommClause
		//     *ast.ForStmt
		//     *ast.RangeStmt
		switch n := n.node.(type) {
		case *ast.File:
			// panic("not support file level scope") // todo
			return info.Scopes[n]
		case *ast.TypeSpec:
			// panic("not support type level scope") // todo
			return info.Scopes[n]

		case *ast.BlockStmt, *ast.CaseClause, *ast.CommClause:
			s := info.Scopes[n]
			if s == nil {
				continue
			}
			return s

		// NOTICE:
		// (BlockStmt body) and (Body *BlockStmt) are sibling
		// NOT parent-child relationship in travel stack
		// so, scope by FuncType when FuncDecl or FuncLit
		// case *ast.FuncType:
		case *ast.FuncDecl, *ast.FuncLit:
			fnTy := getField[*ast.FuncType](n, "Type")
			return info.Scopes[fnTy]

		case *ast.IfStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.ForStmt, *ast.RangeStmt:
			return info.Scopes[n]

		default:
			continue
		}
	}
	return nil
}

func (s travalStack) nearestFunc() ast.Node {
	for _, it := range s {
		switch n := it.node.(type) {
		case *ast.FuncLit, *ast.FuncDecl:
			return n
		default:
			continue
		}
	}
	return nil
}
