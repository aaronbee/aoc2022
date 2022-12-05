package fn

import (
	"bufio"
	"io"

	"golang.org/x/exp/constraints"
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
