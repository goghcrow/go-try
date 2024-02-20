package rewriter

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ast/astutil"
)

type callCtx struct {
	c   *astutil.Cursor
	stk travalStack
	pkg *pkg

	tryFn    string
	callsite *ast.CallExpr /* Try() | TryErr()*/

	outerFun   ast.Node /*FuncLit | FuncDecl*/
	outerStmts *stmts
	outerScope *types.Scope
}

func (c *callCtx) assert(pos positioner, ok bool, format string, a ...any) {
	if !ok {
		loc := c.pkg.Fset.Position(pos.Pos()).String()
		panic(fmt.Sprintf(format, a...) + " in: " + loc)
	}
}
func (c *callCtx) assertCall(ok bool, format string, a ...any) {
	c.assert(c.callsite, ok, format, a...)
}
func (c *callCtx) brokenSemantics(construct string) {
	c.assert(c.callsite, false, "broken semantic when: "+construct)
}

func (c *callCtx) checkNonConsistentSemantics() {
	var parent ast.Node
outer:
	for i, n := range c.stk {
		if i < len(c.stk)-1 {
			parent = c.stk[i+1].node
		}
		switch n.node.(type) {
		case *ast.FuncLit, *ast.FuncDecl:
			break outer

		case *ast.GoStmt:
			c.brokenSemantics("go Try()")
		case *ast.DeferStmt:
			c.brokenSemantics("defer Try()")

		case ast.Stmt:
			switch p := parent.(type) {
			case *ast.ForStmt: // evaluation order
				if n.field == "Post" {
					c.brokenSemantics("for init; ; Try(...) { }")
				}
			case *ast.TypeSwitchStmt: // init dependency
				if n.field == "Assign" && p.Init != nil {
					c.brokenSemantics("switch init; Try(...) { }")
				}
			}
		case ast.Expr:
			switch p := parent.(type) {
			case *ast.CaseClause: // evaluation order
				if n.field == "List" {
					c.brokenSemantics("case Try(...)[,Try(...)...]:")
				}
			case *ast.IfStmt: // evaluation order
				if n.field == "Cond" {
					c.brokenSemantics("if ...; Try(...) { }")
				}

			case *ast.ForStmt: // evaluation order
				if n.field == "Cond" {
					c.brokenSemantics("if ...; Try(...) { }")
				}
			case *ast.SwitchStmt: // init dependency
				if n.field == "Tag" && p.Init != nil {
					c.brokenSemantics("switch init; Try(...) { }")
				}
			}
		}
	}
}

func (c *callCtx) checkSignature() (sig *types.Signature) {
	var (
		typeOf = c.pkg.TypesInfo.TypeOf
		argCnt = len(c.callsite.Args)
		retErr types.Type
		argErr types.Type
		errTy  = types.Universe.Lookup("error").Type()
		niled  = func(t types.Type) bool {
			// nil is a valid value for the following.
			// Pointers / Unsafe pointers / Interfaces / Channels / Maps / Slices / Functions
			switch t.Underlying().(type) {
			case *types.Pointer, *types.Interface, *types.Chan, *types.Map, *types.Slice, *types.Signature:
				return true
			default:
				return false
			}
		}
	)

	c.assertCall(argCnt > 0, "at least one arg")
	lastArg := c.callsite.Args[argCnt-1]

	var fnPos positioner
	switch n := c.outerFun.(type) {
	case *ast.FuncLit:
		sig, _ = typeOf(n).(*types.Signature)
		fnPos = n.Type
	case *ast.FuncDecl:
		sig, _ = typeOf(n.Name).(*types.Signature)
		fnPos = n.Type
	}
	c.assertCall(sig != nil, "<Try> must be in a Tryable func")
	retCnt := sig.Results().Len()
	c.assertCall(retCnt > 0, "expect at least one error return")
	retErr = sig.Results().At(retCnt - 1).Type()

	// c.assertCall(types.AssignableTo(retErr, errTy), "the last return type MUST assignable to error")
	c.assert(fnPos, types.Identical(retErr, errTy), "the last return type MUST be error, but %s", retErr)

	// Try(err)    				OR		func F() error ; 			Try(F())
	// Try(v1, err) 			OR 		func F() (A, error); 		Try(F())
	// Try(v1, v2, err) 		OR 		func F() (A, B, error); 	Try(F())
	// Try(v1, v2, v3, err) 	OR 		func F() (A, B, C, error); 	Try(F())
	for i, name := range funcTryNames {
		lastParam, paramCnt := i, i+1
		if c.tryFn == name {
			switch argCnt {
			case paramCnt:
				argErr = typeOf(lastArg)
			case 1:
				tup, _ := typeOf(c.callsite.Args[0]).(*types.Tuple)
				c.assertCall(tup != nil && tup.Len() == paramCnt, "invalid args, expect %d", paramCnt)
				argErr = tup.At(lastParam).Type()
			}
			break
		}
	}

	c.assertCall(argErr != nil, "invalid try args")
	c.assertCall(niled(argErr), "nil is not a valid value for %s", argErr.String())
	// c.assertCall(types.AssignableTo(argErr, retErr), "type mismatch, Try(..., ?) expect %v but %v", retErr, argErr)
	c.assertCall(types.AssignableTo(argErr, errTy), "type mismatch, Try(..., ?) expect %v but %v", retErr, argErr)
	return
}

func (c *callCtx) checkShadowedNil() {
	s := c.outerScope
	nilObj := s.Lookup("nil")
	if nilObj == nil {
		_, nilObj = s.LookupParent("nil", token.NoPos)
	}
	assert(nilObj != nil)
	c.assert(nilObj, nilObj.Type() == types.Typ[types.UntypedNil], "nil shadowed in fun scope, please rename it")
}
