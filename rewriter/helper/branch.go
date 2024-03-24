package helper

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"

	"github.com/goghcrow/go-loader"
)

func isNil(n any) bool {
	if n == nil {
		return true
	}
	if v := reflect.ValueOf(n); v.Kind() == reflect.Ptr && v.IsNil() {
		return true
	}
	return false
}

type ErrorHandler func(pos loader.Positioner, msg string)

// JumpTable
// when branch is token.GOTO then stmt is *ast.LabeledStmt
// when branch is labeld token.CONTINUE then stmt is *ast.ForStmt, *ast.RangeStmt
// when branch is labeld token.BREAK then stmt is *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt, *ast.ForStmt, *ast.RangeStmt
type JumpTable map[*ast.BranchStmt]ast.Stmt

func (t JumpTable) JumpTo(s *ast.BranchStmt, target ast.Stmt) bool {
	stmt := t[s]
	if isNil(stmt) {
		return false
	}

	switch s.Tok {
	case token.GOTO:
		switch stmt := stmt.(type) {
		case *ast.LabeledStmt:
			return stmt.Stmt == target
		default:
			panic("unreachable")
		}
	case token.CONTINUE:
		switch stmt.(type) {
		case *ast.ForStmt, *ast.RangeStmt:
			return stmt == target
		default:
			panic("unreachable")
		}
	case token.BREAK:
		switch stmt.(type) {
		case *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt, *ast.ForStmt, *ast.RangeStmt:
			return stmt == target
		default:
			panic("unreachable")
		}
	default:
		panic("unreachable")
	}
	return false
}

func Target(body *ast.BlockStmt, errh ErrorHandler) JumpTable {
	if body == nil {
		return nil
	}
	// scope of all labels in this body
	ls := &labelScope{
		errh:         errh,
		branchTarget: map[*ast.BranchStmt]ast.Stmt{},
	}
	ls.blockBranches(nil, targets{}, nil, body.List)
	return ls.branchTarget
}

type labelScope struct {
	errh   ErrorHandler
	labels map[string]*label // all label declarations inside the function; allocated lazily

	// when goto then stmt is *ast.LabeledStmt
	// when labeld continue then stmt is *ast.ForStmt, *ast.RangeStmt
	// when labeld break then stmt is *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt, *ast.ForStmt, *ast.RangeStmt
	branchTarget map[*ast.BranchStmt]ast.Stmt
}

type label struct {
	parent *block           // block containing this label declaration
	lstmt  *ast.LabeledStmt // statement declaring the label
}

type block struct {
	parent *block           // immediately enclosing block, or nil
	lstmt  *ast.LabeledStmt // labeled statement associated with this block, or nil
}

func (ls *labelScope) err(pos loader.Positioner, format string, args ...interface{}) {
	ls.errh(pos, fmt.Sprintf(format, args...))
}

// declare declares the label introduced by s in block b
func (ls *labelScope) declare(b *block, s *ast.LabeledStmt) {
	name := s.Label.Name
	labels := ls.labels
	if labels == nil {
		labels = make(map[string]*label)
		ls.labels = labels
	} else if alt := labels[name]; alt != nil {
		ls.err(s.Label, "label %s already defined", name)
		// ls.err(s.Label, "label %s already defined at %s", name, alt.lstmt.Label.Pos().String())
	}
	labels[name] = &label{parent: b, lstmt: s}
}

// gotoTarget returns the labeled statement matching the given name and
// declared in block b or any of its enclosing blocks. The result is nil
// if the label is not defined, or doesn't match a valid labeled statement.
func (ls *labelScope) gotoTarget(b *block, name string) *ast.LabeledStmt {
	if l := ls.labels[name]; l != nil {
		for ; b != nil; b = b.parent {
			// b : goto 所在的 block 链
			// l.parent : 包含 label 的 block
			// 如果是先声明的 label, goto 往回跳, 返回 labedStmt
			if l.parent == b {
				return l.lstmt
			}
		}
	}
	return nil
}

// break / continue label, 一定是跳转到对应的 带 block 的标签,
// continue: for
// break: for/switch/typeswitch/select

var invalid = new(ast.LabeledStmt) // singleton to signal invalid enclosing target

// enclosingTarget returns the innermost enclosing labeled statement matching
// the given name. The result is nil if the label is not defined, and invalid
// if the label is defined but doesn't label a valid labeled statement.
func (ls *labelScope) enclosingTarget(b *block, name string) *ast.LabeledStmt {
	if l := ls.labels[name]; l != nil {
		for ; b != nil; b = b.parent {
			if l.lstmt == b.lstmt {
				return l.lstmt
			}
		}
		return invalid
	}
	return nil
}

// targets describes the target statements within which break
// or continue statements are valid.
type targets struct {
	breaks    ast.Stmt // *ForStmt, *RangeStmt, *SwitchStmt, *TypeSwitchStmt, *SelectStmt, or nil
	continues ast.Stmt // *ForStmt, *RangeStmt, or nil
	caseIndex int      // case index of immediately enclosing switch statement, or < 0
}

// blockBranches processes a block's body starting at start and returns the
// list of unresolved (forward) gotos. parent is the immediately enclosing
// block (or nil), ctxt provides information about the enclosing statements,
// and lstmt is the labeled statement associated with this block, or nil.
func (ls *labelScope) blockBranches(parent *block, ctxt targets, lstmt *ast.LabeledStmt, body []ast.Stmt) []*ast.BranchStmt {
	b := &block{parent: parent, lstmt: lstmt}

	var fwdGotos []*ast.BranchStmt

	innerBlock := func(ctxt targets, body []ast.Stmt) {
		// Unresolved forward gotos from the inner block
		// become forward gotos for the current block.
		fwdGotos = append(fwdGotos, ls.blockBranches(b, ctxt, lstmt, body)...)
	}

	// A fallthrough statement counts as last statement in a statement
	// list even if there are trailing empty statements; remove them.
	stmtList := trimTrailingEmptyStmts(body)
	for stmtIndex, stmt := range stmtList {
		lstmt = nil
	L:
		switch s := stmt.(type) {

		case *ast.LabeledStmt:
			// declare non-blank label
			if name := s.Label.Name; name != "_" {
				ls.declare(b, s)

				// resolve matching forward gotos
				i := 0
				for _, fwd := range fwdGotos {
					if fwd.Label.Name == name {
						ls.branchTarget[fwd] = s
					} else {
						// no match - keep forward goto
						fwdGotos[i] = fwd
						i++
					}
				}
				fwdGotos = fwdGotos[:i]

				lstmt = s
			}
			// process labeled statement
			stmt = s.Stmt
			goto L

		case *ast.BranchStmt:
			// unlabeled branch statement
			if s.Label == nil {
				switch s.Tok {
				case token.BREAK:
					if t := ctxt.breaks; t != nil {
						ls.branchTarget[s] = t
					} else {
						ls.err(s, "break is not in a loop, switch, or select")
					}
				case token.CONTINUE:
					if t := ctxt.continues; t != nil {
						ls.branchTarget[s] = t
					} else {
						ls.err(s, "continue is not in a loop")
					}
				case token.FALLTHROUGH:
					msg := "fallthrough statement out of place"
					switch t := ctxt.breaks.(type) {
					case *ast.TypeSwitchStmt:
						msg = "cannot fallthrough in type switch"
					case *ast.SwitchStmt:
						if ctxt.caseIndex < 0 || stmtIndex+1 < len(stmtList) {
							// fallthrough nested in a block or not the last statement
							// use msg as is
						} else if ctxt.caseIndex+1 == len(t.Body.List) {
							msg = "cannot fallthrough final case in switch"
						} else {
							continue // fallthrough ok
							// break // fallthrough ok
						}
					}

					// if t, _ := ctxt.breaks.(*ast.SwitchStmt); t != nil {
					// 	if _, ok := t.Tag.(*ast.TypeSwitchGuard); ok {
					// 		msg = "cannot fallthrough in type switch"
					// 	} else if ctxt.caseIndex < 0 || stmtIndex+1 < len(stmtList) {
					// 		// fallthrough nested in a block or not the last statement
					// 		// use msg as is
					// 	} else if ctxt.caseIndex+1 == len(t.Body.List) {
					// 		msg = "cannot fallthrough final case in switch"
					// 	} else {
					// 		break // fallthrough ok
					// 	}
					// }
					ls.err(s, msg)
				case token.GOTO:
					fallthrough // should always have a label
				default:
					panic("invalid BranchStmt")
				}
				break
			}

			// labeled branch statement
			name := s.Label.Name
			switch s.Tok {
			case token.BREAK:
				// spec: "If there is a label, it must be that of an enclosing
				// "for", "switch", or "select" statement, and that is the one
				// whose execution terminates."
				if t := ls.enclosingTarget(b, name); t != nil {
					switch t := t.Stmt.(type) {
					case *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt, *ast.ForStmt, *ast.RangeStmt:
						ls.branchTarget[s] = t
					default:
						ls.err(s.Label, "invalid break label %s", name)
					}
				} else {
					ls.err(s.Label, "break label not defined: %s", name)
				}

			case token.CONTINUE:
				// spec: "If there is a label, it must be that of an enclosing
				// "for" statement, and that is the one whose execution advances."
				if t := ls.enclosingTarget(b, name); t != nil {
					switch t := t.Stmt.(type) {
					case *ast.ForStmt, *ast.RangeStmt:
						ls.branchTarget[s] = t
					default:
						ls.err(s.Label, "invalid continue label %s", name)
					}
				} else {
					ls.err(s.Label, "continue label not defined: %s", name)
				}

			case token.GOTO:
				if t := ls.gotoTarget(b, name); t != nil {
					ls.branchTarget[s] = t
				} else {
					// label may be declared later - add goto to forward gotos
					fwdGotos = append(fwdGotos, s)
				}

			case token.FALLTHROUGH:
				fallthrough // should never have a label
			default:
				panic("invalid BranchStmt")
			}

		case *ast.BlockStmt:
			// 不能 break, 不能 continue
			inner := targets{ctxt.breaks, ctxt.continues, -1}
			innerBlock(inner, s.List)

		case *ast.IfStmt:
			// 不能 break, 不能 continue
			inner := targets{ctxt.breaks, ctxt.continues, -1}
			innerBlock(inner, s.Body.List)
			if s.Else != nil {
				innerBlock(inner, []ast.Stmt{s.Else})
			}

		case *ast.ForStmt:
			// 又能 break, 又能 continue
			inner := targets{s, s, -1}
			innerBlock(inner, s.Body.List)

		case *ast.RangeStmt:
			// 又能 break, 又能 continue
			inner := targets{s, s, -1}
			innerBlock(inner, s.Body.List)

		case *ast.SwitchStmt:
			// 只能 break, 不能 continue
			inner := targets{s, ctxt.continues, -1}
			for i, cc := range s.Body.List {
				inner.caseIndex = i
				innerBlock(inner, cc.(*ast.CaseClause).Body)
			}

		case *ast.TypeSwitchStmt:
			// 只能 break, 不能 continue
			inner := targets{s, ctxt.continues, -1}
			for i, cc := range s.Body.List {
				inner.caseIndex = i
				innerBlock(inner, cc.(*ast.CaseClause).Body)
			}

		case *ast.SelectStmt:
			// 只能 break, 不能 continue
			inner := targets{s, ctxt.continues, -1}
			for _, cc := range s.Body.List {
				innerBlock(inner, cc.(*ast.CommClause).Body)
			}
		}
	}

	return fwdGotos
}

func trimTrailingEmptyStmts(list []ast.Stmt) []ast.Stmt {
	for i := len(list); i > 0; i-- {
		if _, ok := list[i-1].(*ast.EmptyStmt); !ok {
			return list[:i]
		}
	}
	return nil
}

func assert(ok bool, msg string) {
	if !ok {
		panic(msg)
	}
}
