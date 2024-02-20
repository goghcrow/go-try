package rewriter

func prepend[E any, S ~[]E](xs S, y E) S {
	var z E
	xs = append(xs, z)
	copy(xs[1:], xs)
	xs[0] = y
	return xs
}

func grow[E any, S ~[]E](xs S, need int) S {
	cnt := len(xs)
	if need <= cnt {
		return xs
	}
	var z E
	for cap(xs) < need {
		xs = append(xs[:cap(xs)], z)
	}
	return xs[:cnt]
}

func insertBefore[E any, S ~[]E](xs S, idx int, ys ...E) S {
	need := len(xs) + len(ys)
	xs = grow(xs, need)[:need]
	copy(xs[idx+len(ys):], xs[idx:])
	copy(xs[idx:], ys)
	return xs
}

func reverseInplace[E any, S ~[]E](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func reverse[E any, S ~[]E](s S) S {
	a := make(S, len(s))
	copy(a, s)
	reverseInplace(a)
	return a
}

func mapto[E, R any, S ~[]E](s S, f func(int, E) R) []R {
	a := make([]R, len(s))
	for i, it := range s {
		a[i] = f(i, it)
	}
	return a
}

func filter[E any, S ~[]E](s S, f func(E) bool) []E {
	var a []E
	for _, it := range s {
		if f(it) {
			a = append(a, it)
		}
	}
	return a
}

func search[E comparable, S ~[]E](s S, x E) int {
	for i, it := range s {
		if it == x {
			return i
		}
	}
	return -1
}

func first[E any, S ~[]E](s S, f func(E) bool) *E {
	for _, it := range s {
		if f(it) {
			return &it
		}
	}
	return nil
}

func reverseFirst[E any, S ~[]E](s S, f func(E) bool) *E {
	for i := len(s) - 1; i >= 0; i-- {
		if f(s[i]) {
			return &s[i]
		}
	}
	return nil
}
