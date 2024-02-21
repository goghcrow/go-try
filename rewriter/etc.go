package rewriter

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
)

var _dbg = true

func assert(ok bool) {
	if !ok {
		panic("illegal state")
	}
}

func instanceof[T any](x any) (ok bool) {
	_, ok = x.(T)
	return
}

type reflv = reflect.Value

// v must be a pointer to a struct or an interface
func getField[T any](v any, name string) T {
	rv, ok := v.(reflv)
	if !ok {
		rv = reflect.ValueOf(v)
	}
	return rv.Elem().FieldByName(name).Interface().(T)
}

func getFieldPtr[T any](v any, name string) *T {
	rv, ok := v.(reflv)
	if !ok {
		rv = reflect.ValueOf(v)
	}
	return rv.Elem().FieldByName(name).Addr().Interface().(*T)
}

// row polymorphism needed...
//
//	xs = []*struct {
//		Names   []*Ident
//		Type    Expr
//		...
//	}
func flattenNamesByType[T any](xs any, mk func(name *ast.Ident, expr ast.Expr) T) (ys []T) {
	var (
		V     = reflect.ValueOf
		Type  = func(v reflv) ast.Expr { return getField[ast.Expr](v, "Type") }
		Names = func(v reflv) []*ast.Ident { return getField[[]*ast.Ident](v, "Names") }
	)

	rxs := V(xs)
	for i := 0; i < rxs.Len(); i++ {
		x := rxs.Index(i)
		ns := Names(x)
		cnt := len(ns)
		switch cnt {
		case 0:
			ys = append(ys, mk(nil, Type(x)))
		case 1:
			ys = append(ys, mk(ns[0], Type(x)))
		default:
			for j := 0; j < cnt; j++ {
				ys = append(ys, mk(ns[j], Type(x)))
			}
		}
	}
	return
}

// row polymorphism needed...
//
//	xs = []*struct {
//		Names   []*Ident
//		Type    Expr
//		...
//	}
//
// x must assign to T
func groupNamesByType[T any](xs any, tyEq func(x, y ast.Expr) bool) (ys []T) {
	var (
		V        = reflect.ValueOf
		Type     = func(v reflv) ast.Expr { return getField[ast.Expr](v, "Type") }
		Names    = func(v reflv) []*ast.Ident { return getField[[]*ast.Ident](v, "Names") }
		SetNames = func(v reflv, ns []*ast.Ident) { v.Elem().FieldByName("Names").Set(V(ns)) }
	)

	rxs := V(xs)
	for i := 0; i < rxs.Len(); i++ {
		x := rxs.Index(i)
		pre := len(ys) - 1
		if pre >= 0 {
			y := V(ys[pre])
			if tyEq(Type(x), Type(y)) {
				SetNames(y, append(Names(y), Names(x)...))
				continue
			}
		}
		ys = append(ys, x.Interface().(T))
	}
	return ys
}

// importPath returns the unquoted import path of s
func importPath(s *ast.ImportSpec) string {
	t, err := strconv.Unquote(s.Path.Value)
	assert(err == nil) // assume well-syntaxed
	return t
}

// imported package by path, whatever the name is
func imported(f *ast.File, path string) bool {
	for _, decl := range f.Decls {
		d, ok := decl.(*ast.GenDecl)
		if !ok || d.Tok != token.IMPORT {
			continue
		}
		for _, spec := range d.Specs {
			s := spec.(*ast.ImportSpec)
			if importPath(s) == path {
				return true
			}
		}
	}
	return false
}

// deleteImport deletes the import with the given path (ignoring the name) from the file f,
// if present. If there are duplicate import declarations, all matching ones are deleted.
// modified from https://github.com/golang/tools/blob/master/go/ast/astutil/imports.go#L210
func deleteImport(fset *token.FileSet, f *ast.File, path string) (deleted bool) {
	var delspecs []*ast.ImportSpec
	var delcomments []*ast.CommentGroup

	// Find the import nodes that import path, if any.
	for i := 0; i < len(f.Decls); i++ {
		decl := f.Decls[i]
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.IMPORT {
			continue
		}
		for j := 0; j < len(gen.Specs); j++ {
			spec := gen.Specs[j]
			impspec := spec.(*ast.ImportSpec)
			if /*importName(impspec) != parentField || */ importPath(impspec) != path {
				continue
			}

			// We found an import spec that imports path.
			// Delete it.
			delspecs = append(delspecs, impspec)
			deleted = true
			copy(gen.Specs[j:], gen.Specs[j+1:])
			gen.Specs = gen.Specs[:len(gen.Specs)-1]

			// If this was the last import spec in this decl,
			// delete the decl, too.
			if len(gen.Specs) == 0 {
				copy(f.Decls[i:], f.Decls[i+1:])
				f.Decls = f.Decls[:len(f.Decls)-1]
				i--
				break
			} else if len(gen.Specs) == 1 {
				if impspec.Doc != nil {
					delcomments = append(delcomments, impspec.Doc)
				}
				if impspec.Comment != nil {
					delcomments = append(delcomments, impspec.Comment)
				}
				for _, cg := range f.Comments {
					// Found comment on the same line as the import spec.
					if cg.End() < impspec.Pos() && fset.Position(cg.End()).Line == fset.Position(impspec.Pos()).Line {
						delcomments = append(delcomments, cg)
						break
					}
				}

				spec := gen.Specs[0].(*ast.ImportSpec)

				// Move the documentation right after the import decl.
				if spec.Doc != nil {
					for fset.Position(gen.TokPos).Line+1 < fset.Position(spec.Doc.Pos()).Line {
						fset.File(gen.TokPos).MergeLine(fset.Position(gen.TokPos).Line)
					}
				}
				for _, cg := range f.Comments {
					if cg.End() < spec.Pos() && fset.Position(cg.End()).Line == fset.Position(spec.Pos()).Line {
						for fset.Position(gen.TokPos).Line+1 < fset.Position(spec.Pos()).Line {
							fset.File(gen.TokPos).MergeLine(fset.Position(gen.TokPos).Line)
						}
						break
					}
				}
			}
			if j > 0 {
				lastImpspec := gen.Specs[j-1].(*ast.ImportSpec)
				lastLine := fset.PositionFor(lastImpspec.Path.ValuePos, false).Line
				line := fset.PositionFor(impspec.Path.ValuePos, false).Line

				// We deleted an entry but now there may be
				// a blank line-sized hole where the import was.
				if line-lastLine > 1 || !gen.Rparen.IsValid() {
					// There was a blank line immediately preceding the deleted import,
					// so there's no need to close the hole. The right parenthesis is
					// invalid after AddImport to an import statement without parenthesis.
					// Do nothing.
				} else if line != fset.File(gen.Rparen).LineCount() {
					// There was no blank line. Close the hole.
					fset.File(gen.Rparen).MergeLine(line)
				}
			}
			j--
		}
	}

	// Delete imports from f.Imports.
	for i := 0; i < len(f.Imports); i++ {
		imp := f.Imports[i]
		for j, del := range delspecs {
			if imp == del {
				copy(f.Imports[i:], f.Imports[i+1:])
				f.Imports = f.Imports[:len(f.Imports)-1]
				copy(delspecs[j:], delspecs[j+1:])
				delspecs = delspecs[:len(delspecs)-1]
				i--
				break
			}
		}
	}

	// Delete comments from f.Comments.
	for i := 0; i < len(f.Comments); i++ {
		cg := f.Comments[i]
		for j, del := range delcomments {
			if cg == del {
				copy(f.Comments[i:], f.Comments[i+1:])
				f.Comments = f.Comments[:len(f.Comments)-1]
				copy(delcomments[j:], delcomments[j+1:])
				delcomments = delcomments[:len(delcomments)-1]
				i--
				break
			}
		}
	}

	if len(delspecs) > 0 {
		panic(fmt.Sprintf("deleted specs from Decls but not Imports: %v", delspecs))
	}

	return
}
