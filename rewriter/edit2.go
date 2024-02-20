package rewriter

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/ast/astutil"
)

func (r *rewriter) editFile2(f *ast.File, pkg *pkg) {
	type ctx struct {
		node travalNode

		scope *types.Scope

		stmts   *[]ast.Stmt // TO MODIFY, so ptr of slice
		stmtIdx int

		fn    ast.Node
		fnSym int
	}

	var (
		stk   = mkStack[*ctx]()
		enter = stk.push
		exit  = stk.pop

		nearestFun = func(pos positioner) ast.Node {
			fst := reverseFirst(stk.unwrap(), func(it *ctx) bool { return it.fn != nil })
			// assert(fst != nil)
			if fst == nil {
				return nil
			}
			return (*fst).fn
		}
		nearestScope = func() *types.Scope {
			fst := reverseFirst(stk.unwrap(), func(it *ctx) bool { return it.scope != nil })
			// assert(fst != nil)
			if fst == nil {
				return nil
			}
			return (*fst).scope
		}
		nearestStmts = func() *stmts {
			it := stk.unwrap()
			for i := len(it) - 1; i >= 1; i-- {
				if it[i].stmts != nil {
					// *SwitchStmt.Body, *TypeSwitchStmt.Body:  CaseClauses only, can't be inserted err-if stmt, ignored
					// *SelectStmt.Body:  CommClauses only, can't be inserted err-if stmt, ignored
					switch (*it[i].stmts)[0].(type) {
					case *ast.CaseClause, *ast.CommClause:
						continue
					}

					assert(it[i+1].stmtIdx >= 0)
					parentStmts := it[i].stmts
					childIdxPtr := &it[i+1].stmtIdx
					return newStmts(parentStmts, childIdxPtr)
				}
			}
			// panic("illegal state")
			return nil
		}
		travalStk = func() travalStack {
			return reverse(mapto(stk.unwrap(), func(_ int, it *ctx) travalNode {
				return it.node
			}))
		}
	)

	var (
		info  = pkg.TypesInfo
		cache = map[ast.Node]bool{}
	)
	astutil.Apply(f, func(c *astutil.Cursor) bool {
		n := c.Node()
		if n == nil {
			return false
		}
		if cache[n] {
			return false
		}
		cache[n] = true

		// The following node types may appear in Scopes:
		//	(update scope <1>)
		//     *ast.File
		//     *ast.TypeSpec
		//     *ast.BlockStmt
		//     *ast.IfStmt
		//     *ast.SwitchStmt, *ast.TypeSwitchStmt
		//     *ast.CaseClause, *ast.CommClause
		//     *ast.ForStmt, *ast.RangeStmt
		//  (update scope <2>)
		//     *ast.FuncType <- *ast.FuncDecl, *ast.FuncLit

		// The following node has []ast.Stmt field:
		//		*ast.BlockStmt.List 									(update stmts <1>)
		//		*ast.CaseClause.Body / *ast.CommClause.Body				(update stmts <2>)
		// The following node has BlockStmt field:  					(update stmts <1>)
		// 		*ast.FuncLit.Body / *ast.FuncDecl.Body
		//		*ast.IfStmt.Body
		//		*ast.SwitchStmt.Body / *ast.TypeSwitchStmt.Body
		//		*ast.SelectStmt.Body
		//		*ast.ForStmt.Body / *ast.RangeStmt.Body

		x := &ctx{
			node: travalNode{
				node:  n,
				field: c.Name(),
			},
			scope:   info.Scopes[n], // update scope <1>
			stmtIdx: c.Index(),
		}

		switch n := n.(type) {
		case *ast.BlockStmt:
			x.stmts = &n.List // update stmts <1>
		case *ast.FuncDecl, *ast.FuncLit:
			// NOTICE: (Type *FuncType) and (Body *BlockStmt) are sibling instead of parent-child
			// relation in travel stack, so, scoping by *FuncType when FuncDecl or FuncLit
			fnTy := getField[*ast.FuncType](n, "Type")
			x.scope = info.Scopes[fnTy] // update scope <2>

			x.fn, x.fnSym = n, r.symCnt // update fun <1>
			r.resetsym()
		}

		// (Body []Stmt) instead of (Body *ast.BlockStmt) in CaseClause/CommClause,
		if _, ok := n.(ast.Stmt); ok && c.Index() >= 0 {
			switch p := c.Parent().(type) {
			case *ast.CaseClause, *ast.CommClause:
				if c.Name() == "Body" {
					// attaching []Stmt to case/comm clause after
					// rewritting CaseCluase.List or CommClause.Comm (order by declaration)
					// for nearestStmts searching
					if stk.top().stmts == nil {
						stk.top().stmts = getFieldPtr[[]ast.Stmt](p, "Body") // update stmts <2>
					}
				}
			}
		}

		enter(x)
		return true
	}, func(c *astutil.Cursor) bool {
		n := c.Node()
		switch n := n.(type) {
		case *ast.FuncDecl, *ast.FuncLit: // case *ast.FuncType:
			r.symCnt = stk.top().fnSym
		case *ast.CallExpr:
			_, tryFn := r.tryCallee(info, n)
			if tryFn != "" {
				fun := nearestFun(n)
				stmts := nearestStmts()
				scope := nearestScope()
				tvlStk := travalStk()

				r.assert(pkg, n, fun != nil, "Try must be in a tryable fun(...) (T, error)")
				r.assert(pkg, n, stmts != nil, "missing stmts")
				r.assert(pkg, n, scope != nil, "missing scope")

				if _dbg {
					fun1 := tvlStk.nearestFunc()
					assert(fun == fun1)

					stmts1 := tvlStk.nearestStmts()
					assert((stmts == nil && stmts1 == nil) ||
						(*stmts.idx == *stmts1.idx && stmts.xs == stmts1.xs))

					scope1 := tvlStk.nearestScope(info)
					assert(scope == scope1)
				}

				r.rewriteTryCall(&callCtx{
					c:          c,
					stk:        tvlStk,
					pkg:        pkg,
					tryFn:      tryFn,
					callsite:   n,
					outerFun:   fun,
					outerStmts: stmts,
					outerScope: scope,
				})
			}
		}
		exit()
		return true
	})
}

// ↓↓↓↓↓↓↓↓↓↓↓↓↓↓ astutil.Apply helper ↓↓↓↓↓↓↓↓↓↓↓↓↓↓

type stack[T any] []T

func mkStack[T any]() *stack[T] {
	var z T
	return &stack[T]{z}
}

func (s *stack[T]) push(v T) { *s = append(*s, v) }
func (s *stack[T]) pop() T {
	assert(len(*s) > 1)
	v := s.top()
	*s = (*s)[:len(*s)-1]
	return v
}
func (s *stack[T]) top() T {
	assert(len(*s) > 1)
	return (*s)[len(*s)-1]
}
func (s *stack[T]) len() int    { return len(*s) - 1 }
func (s *stack[T]) unwrap() []T { return (*s)[1:] }
