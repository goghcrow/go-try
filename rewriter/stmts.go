package rewriter

import "go/ast"

type stmts struct {
	idx int
	xs  *[]ast.Stmt
}

func newStmts(xs *[]ast.Stmt, idx int) *stmts {
	return &stmts{idx: idx, xs: xs}
}

func (s *stmts) insertAfter(stmts ...ast.Stmt) {
	idx := s.idx
	nxs := make([]ast.Stmt, len(*s.xs)+len(stmts))
	copy(nxs, (*s.xs)[:idx])
	copy(nxs[idx+2:], (*s.xs)[idx:])
	for i, x := range stmts {
		nxs[idx+i] = x
	}
	*s.xs = nxs
}
