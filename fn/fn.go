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
