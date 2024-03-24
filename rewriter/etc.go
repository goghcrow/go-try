package rewriter

import (
	"go/ast"
	"go/types"
	"reflect"
)

func assert(ok bool) {
	if !ok {
		panic("illegal state")
	}
}

func instanceof[T any](x any) (ok bool) {
	_, ok = x.(T)
	return
}

func isValidNil(t types.Type) bool {
	// nil is a valid value for the following.
	// Pointers / Unsafe pointers / Interfaces / Channels / Maps / Slices / Functions
	switch t.Underlying().(type) {
	case *types.Pointer, *types.Interface, *types.Chan, *types.Map, *types.Slice, *types.Signature:
		return true
	default:
		return false
	}
}

// ======================================================================

func sliceOf[T any](xs ...T) []T {
	return xs
}

func concat[E any, S ~[]E](xs ...S) S {
	return foldLeft(xs, make([]E, 0, foldLeft(xs, 0, func(i int, s S) int {
		return i + len(s)
	})), func(acc []E, el S) []E {
		return append(acc, el...)
	})
}

func prepend[E any, S ~[]E](xs S, ys ...E) S {
	if len(ys) == 0 {
		return xs
	}
	xs = append(xs, ys...)
	copy(xs[len(ys):], xs)
	copy(xs, ys)
	return xs
}

func sliceMap[E, R any, S ~[]E](s S, f func(int, E) R) []R {
	a := make([]R, len(s))
	for i, it := range s {
		a[i] = f(i, it)
	}
	return a
}

func findFirst[E any, S ~[]E](s S, f func(E) bool) *E {
	for _, it := range s {
		if f(it) {
			return &it
		}
	}
	return nil
}

func foldLeft[E any, S ~[]E, R any](xs S, z R, op func(R, E) R) R {
	for _, x := range xs {
		z = op(z, x)
	}
	return z
}

func foldRight[E any, S ~[]E, R any](xs S, z R, op func(E, R) R) R {
	for i := len(xs) - 1; i >= 0; i-- {
		z = op(xs[i], z)
	}
	return z
}

// ======================================================================

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
