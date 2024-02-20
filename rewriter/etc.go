package rewriter

import (
	"go/ast"
	"reflect"
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
