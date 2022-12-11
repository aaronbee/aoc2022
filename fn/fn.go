package fn

import (
	"bufio"
	"io"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func ReduceReader[A any](r io.Reader, initial A, fn func(acc A, line string) A) A {
	acc := initial
	s := bufio.NewScanner(r)
	for s.Scan() {
		acc = fn(acc, s.Text())
	}
	if s.Err() != nil {
		panic(s.Err())
	}
	return acc
}

func Reduce[A any, E any](es []E, initial A, fn func(acc A, elem E) A) A {
	acc := initial
	for _, e := range es {
		acc = fn(acc, e)
	}
	return acc
}

func Map[A, B any](as []A, fn func(a A) B) []B {
	bs := make([]B, len(as))
	for i, a := range as {
		bs[i] = fn(a)
	}
	return bs
}

type addable interface {
	constraints.Integer | constraints.Float | constraints.Complex | ~string
}

func Sum[T addable](ns []T) T {
	var sum T
	for _, n := range ns {
		sum += n
	}
	return sum
}

func Max[T constraints.Integer](ns []T) T {
	if len(ns) == 0 {
		return 0
	}
	return Reduce(ns[1:], ns[0], func(acc T, n T) T {
		if n > acc {
			return n
		}
		return acc
	})
}

func Map2Func[K comparable, E any](m map[K]E) func(K) E {
	return func(k K) E { return m[k] }
}

func Slice2Set[T comparable](ts []T) map[T]bool {
	return Reduce(ts, make(map[T]bool), func(m map[T]bool, t T) map[T]bool {
		m[t] = true
		return m
	})
}

func Count[T any](ts []T, fn func(T) bool) int {
	var acc int
	for _, t := range ts {
		if fn(t) {
			acc++
		}
	}
	return acc
}

func Reverse[T any](ts []T) {
	for b, e := 0, len(ts)-1; b < e; b, e = b+1, e-1 {
		ts[b], ts[e] = ts[e], ts[b]
	}
}

func TopN[T constraints.Ordered](ts []T, n int) []T {
	slices.Sort(ts)
	return ts[len(ts)-n:]
}
