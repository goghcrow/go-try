package helper

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

// <<golang spec>>
// https://golang.org/ref/spec#Terminating_statements
// Terminating statements
// A terminating statement interrupts the regular flow of control in a block. The following statements are terminating:
//
// 1. A "return" or "goto" statement.
// 2. A call to the built-in function panic.
// 3. A block in which the statement list ends in a terminating statement.
// 4. An "if" statement in which:
// 		the "else" branch is present, and
// 		both branches are terminating statements.
// 5. A "for" statement in which:
// 		there are no "break" statements referring to the "for" statement, and
// 		the loop condition is absent, and
// 		the "for" statement does not use a range clause.
// 6. A "switch" statement in which:
// 		there are no "break" statements referring to the "switch" statement,
// 		there is a default case, and
// 		the statement lists in each case, including the default, end in a terminating statement, or a possibly labeled "fallthrough" statement.
// 7. A "select" statement in which:
// 		there are no "break" statements referring to the "select" statement, and
// 		the statement lists in each case, including the default if present, end in a terminating statement.
// 8. A labeled statement labeling a terminating statement.
//
// All other statements are not terminating.
//
// A statement list ends in a terminating statement if the list is not empty and its final non-empty statement is terminating.

// modified from go/src/go/types/return.go

type TerminationChecker struct {
	panicCallSites map[*ast.CallExpr]bool
}

func NewTerminationChecker(panicCallSites map[*ast.CallExpr]bool) *TerminationChecker {
	return &TerminationChecker{
		panicCallSites: panicCallSites,
	}
}

// IsTerminating reports if s is a terminating statement.
// If s is labeled, label is the label name; otherwise s
// is "".
func (t *TerminationChecker) IsTerminating(s ast.Stmt, label string) bool {
	switch s := s.(type) {
	default:
		panic("unreachable")

	case *ast.BadStmt, *ast.DeclStmt, *ast.EmptyStmt, *ast.SendStmt,
		*ast.IncDecStmt, *ast.AssignStmt, *ast.GoStmt, *ast.DeferStmt,
		*ast.RangeStmt:
		// no chance

	case *ast.LabeledStmt:
		return t.IsTerminating(s.Stmt, s.Label.Name)

	case *ast.ExprStmt:
		// calling the predeclared (possibly parenthesized) panic() function is terminating
		if call, ok := astutil.Unparen(s.X).(*ast.CallExpr); ok && t.panicCallSites[call] {
			return true
		}

	case *ast.ReturnStmt:
		return true

	case *ast.BranchStmt:
		if s.Tok == token.GOTO || s.Tok == token.FALLTHROUGH {
			return true
		}

	case *ast.BlockStmt:
		return t.isTerminatingList(s.List, "")

	case *ast.IfStmt:
		if s.Else != nil &&
			t.IsTerminating(s.Body, "") &&
			t.IsTerminating(s.Else, "") {
			return true
		}

	case *ast.SwitchStmt:
		return t.isTerminatingSwitch(s.Body, label)

	case *ast.TypeSwitchStmt:
		return t.isTerminatingSwitch(s.Body, label)

	case *ast.SelectStmt:
		for _, s := range s.Body.List {
			cc := s.(*ast.CommClause)
			if !t.isTerminatingList(cc.Body, "") || hasBreakList(cc.Body, label, true) {
				return false
			}

		}
		return true

	case *ast.ForStmt:
		if s.Cond == nil && !hasBreak(s.Body, label, true) {
			return true
		}
	}

	return false
}

func (t *TerminationChecker) isTerminatingList(list []ast.Stmt, label string) bool {
	// trailing empty statements are permitted - skip them
	for i := len(list) - 1; i >= 0; i-- {
		if _, ok := list[i].(*ast.EmptyStmt); !ok {
			return t.IsTerminating(list[i], label)
		}
	}
	return false // all statements are empty
}

func (t *TerminationChecker) isTerminatingSwitch(body *ast.BlockStmt, label string) bool {
	hasDefault := false
	for _, s := range body.List {
		cc := s.(*ast.CaseClause)
		if cc.List == nil {
			hasDefault = true
		}
		if !t.isTerminatingList(cc.Body, "") || hasBreakList(cc.Body, label, true) {
			return false
		}
	}
	return hasDefault
}

// TODO(gri) For nested breakable statements, the current implementation of hasBreak
// will traverse the same subtree repeatedly, once for each label. Replace
// with a single-pass label/break matching phase.

// hasBreak reports if s is or contains a break statement
// referring to the label-ed statement or implicit-ly the
// closest outer breakable statement.
func hasBreak(s ast.Stmt, label string, implicit bool) bool {
	switch s := s.(type) {
	default:
		panic("unreachable")

	case *ast.BadStmt, *ast.DeclStmt, *ast.EmptyStmt, *ast.ExprStmt,
		*ast.SendStmt, *ast.IncDecStmt, *ast.AssignStmt, *ast.GoStmt,
		*ast.DeferStmt, *ast.ReturnStmt:
		// no chance

	case *ast.LabeledStmt:
		return hasBreak(s.Stmt, label, implicit)

	case *ast.BranchStmt:
		if s.Tok == token.BREAK {
			if s.Label == nil {
				return implicit
			}
			if s.Label.Name == label {
				return true
			}
		}

	case *ast.BlockStmt:
		return hasBreakList(s.List, label, implicit)

	case *ast.IfStmt:
		if hasBreak(s.Body, label, implicit) ||
			s.Else != nil && hasBreak(s.Else, label, implicit) {
			return true
		}

	case *ast.CaseClause:
		return hasBreakList(s.Body, label, implicit)

	case *ast.SwitchStmt:
		if label != "" && hasBreak(s.Body, label, false) {
			return true
		}

	case *ast.TypeSwitchStmt:
		if label != "" && hasBreak(s.Body, label, false) {
			return true
		}

	case *ast.CommClause:
		return hasBreakList(s.Body, label, implicit)

	case *ast.SelectStmt:
		if label != "" && hasBreak(s.Body, label, false) {
			return true
		}

	case *ast.ForStmt:
		if label != "" && hasBreak(s.Body, label, false) {
			return true
		}

	case *ast.RangeStmt:
		if label != "" && hasBreak(s.Body, label, false) {
			return true
		}
	}

	return false
}

func hasBreakList(list []ast.Stmt, label string, implicit bool) bool {
	for _, s := range list {
		if hasBreak(s, label, implicit) {
			return true
		}
	}
	return false
}
