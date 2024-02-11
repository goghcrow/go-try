package rewriter

import (
	"go/ast"
	"go/types"
	"reflect"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/types/typeutil"
)

func (r *rewriter) editFile_(f *ast.File) {
	type ctx struct {
		xs  *[]ast.Stmt
		idx int
		s   *types.Scope
	}
	type fnCtx struct {
		f      ast.Node
		symCnt int
	}
	var (
		fnStk   = mkStack[fnCtx]()
		enterFn = fnStk.push
		exitFn  = fnStk.pop

		ctxStk = mkStack[ctx]()
		enter  = ctxStk.push
		exit   = ctxStk.pop
	)

	var (
		nearestFun = func(pos positioner) ast.Node {
			it := fnStk.top()
			r.assert(it.f != nil, pos, "Try must be in a tryable fun(...) (T, error)")
			return it.f
		}
		nearestScope = func() *types.Scope {
			it := ctxStk.unwrap()
			for i := len(it) - 1; i >= 0; i-- {
				if it[i].s != nil {
					return it[i].s
				}
			}
			panic("illegal state")
		}
		nearestStmts = func() *stmts {
			it := ctxStk.unwrap()
			for i := len(it) - 1; i >= 1; i-- {
				if it[i].xs != nil {
					return newStmts(
						it[i].xs,
						ctxStk.unwrap()[i+1].idx,
					)
				}
			}
			panic("illegal state")
		}
	)

	cache := map[ast.Node]bool{}
	astutil.Apply(f, func(c *astutil.Cursor) bool {
		n := c.Node()
		if cache[n] {
			return false
		}
		cache[n] = true

		var xs *[]ast.Stmt
		switch n := n.(type) {
		case *ast.BlockStmt:
			xs = &n.List
		case *ast.CaseClause, *ast.CommClause:
			xs = reflect.ValueOf(n).Elem().FieldByName("Body").Addr().Interface().(*[]ast.Stmt)
		case *ast.FuncDecl, *ast.FuncLit,
			*ast.IfStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.ForStmt, *ast.RangeStmt:
			switch n.(type) {
			case *ast.FuncDecl, *ast.FuncLit:
				enterFn(fnCtx{f: n, symCnt: r.symCnt})
				r.symCnt = 0
			}
			body := reflect.ValueOf(n).Elem().FieldByName("Body").Interface().(*ast.BlockStmt)
			xs = &body.List
		}

		enter(ctx{
			xs:  xs,
			idx: c.Index(),
			s:   r.m.Scopes[n], // FuncType when FuncDecl or FuncLit
		})
		return true
	}, func(c *astutil.Cursor) bool {
		switch n := c.Node().(type) {
		case *ast.FuncDecl, *ast.FuncLit: // case *ast.FuncType:
			r.symCnt = exitFn().symCnt
		case *ast.CallExpr:
			isTryCall := typeutil.Callee(r.m.Info, n) == r.tryFunc
			if isTryCall {
				replace := r.rewriteTryCall(
					n,
					nearestFun(n),
					nearestStmts(),
					nearestScope(),
				)
				c.Replace(replace)
			}
		}
		exit()
		return true
	})
}
