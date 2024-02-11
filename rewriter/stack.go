package rewriter

import (
	"go/ast"
	"go/types"
	"reflect"

	matcher "github.com/goghcrow/go-ast-matcher"
)

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

// ↓↓↓↓↓↓↓↓↓↓↓↓↓↓ node stack ↓↓↓↓↓↓↓↓↓↓↓↓↓↓

type nodeStack []ast.Node

func (s nodeStack) nearestStmt() *stmts {
	stk := []ast.Node(s)
	_, ok := stk[0].(*ast.CallExpr)
	assert(ok)

	indexOf := func(xs []ast.Stmt, x ast.Node) int {
		for idx, stmt := range xs {
			if x == stmt {
				return idx
			}
		}
		panic("illegal state")
	}

	for i, n := range stk {
		switch n := n.(type) {

		case *ast.BlockStmt:
			return newStmts(&n.List, indexOf(n.List, stk[i-1]))

		case *ast.CaseClause, *ast.CommClause:
			body := reflect.ValueOf(n).Elem().FieldByName("Body").Addr().Interface().(*[]ast.Stmt)
			return newStmts(body, indexOf(*body, stk[i-1]))

		case *ast.FuncDecl, *ast.FuncLit, // case *ast.FuncType:
			*ast.IfStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.ForStmt, *ast.RangeStmt:
			body := reflect.ValueOf(n).Elem().FieldByName("Body").Interface().(*ast.BlockStmt)
			return newStmts(&body.List, indexOf(body.List, stk[i-1]))

		default:
			continue
		}
	}
	panic("missing stmts")
}

func (s nodeStack) nearestScope(m *matcher.Matcher) *types.Scope {
	stk := []ast.Node(s)
	for _, n := range stk {
		switch n := n.(type) {
		case *ast.File:
			panic("not support file level scope")
		case *ast.TypeSpec:
			panic("not support type level scope")

		case *ast.BlockStmt, *ast.CaseClause, *ast.CommClause:
			s := m.Scopes[n]
			if s == nil {
				continue
			}
			return s

		// NOTICE: scope by FuncType when FuncDecl or FuncLit
		case *ast.FuncDecl, *ast.FuncLit: // case *ast.FuncType:
			funTy := reflect.ValueOf(n).Elem().FieldByName("Type").Interface().(*ast.FuncType)
			return m.Scopes[funTy]

		case *ast.IfStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.ForStmt, *ast.RangeStmt:
			return m.Scopes[n]

		default:
			continue
		}
	}
	return nil
}

func (s nodeStack) nearestFunc() ast.Node {
	stk := []ast.Node(s)
	for _, it := range stk {
		switch n := it.(type) {
		case *ast.FuncLit, *ast.FuncDecl:
			return n
		default:
			continue
		}
	}
	return nil
}
