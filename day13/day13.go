package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/aaronbee/aoc2022/fn"
	"golang.org/x/exp/slices"
)

type pair struct {
	l, r []any // l and r slices contain []any or float64
}

func (p pair) inOrder() bool {
	return compareLists(p.l, p.r) < 0
}

func compare(l, r any) int {
	switch l := l.(type) {
	case float64:
		switch r := r.(type) {
		case float64:
			if l < r {
				return -1
			} else if l > r {
				return 1
			}
			return 0
		case []any:
			return compareLists([]any{l}, r)
		}
	case []any:
		switch r := r.(type) {
		case float64:
			return compareLists(l, []any{r})
		case []any:
			return compareLists(l, r)
		}
	}
	panic(fmt.Errorf("unexpected values %T(%v) %T(%v)", l, l, r, r))
}

func compareLists(l, r []any) int {
	for i := range l {
		if i >= len(r) {
			return 1
		}
		if c := compare(l[i], r[i]); c != 0 {
			return c
		}
	}
	if len(l) == len(r) {
		return 0
	}
	return -1
}

func part1(ps []pair) int {
	var sum int
	for i, p := range ps {
		if p.inOrder() {
			sum += i + 1
		}
	}
	return sum
}

func part2(ps []pair) int {
	dividers := [][]any{
		[]any{[]any{float64(2)}},
		[]any{[]any{float64(6)}},
	}
	packets := fn.Reduce(ps, nil, func(acc [][]any, p pair) [][]any {
		return append(acc, p.l, p.r)
	})
	packets = append(packets, dividers...)

	less := func(a, b []any) bool { return compareLists(a, b) < 0 }
	slices.SortFunc(packets, less)

	equaler := func(a []any) func([]any) bool {
		return func(b []any) bool { return compareLists(a, b) == 0 }
	}
	div0 := slices.IndexFunc(packets, equaler(dividers[0]))
	div1 := slices.IndexFunc(packets, equaler(dividers[1]))

	return (div0 + 1) * (div1 + 1)
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	d := json.NewDecoder(f)
	var ps []pair
	for {
		var p pair
		if err := d.Decode(&p.l); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}
		if err := d.Decode(&p.r); err != nil {
			panic(err)
		}
		ps = append(ps, p)
	}
	fmt.Println("Part 1:", part1(ps))
	fmt.Println("Part 2:", part2(ps))
}
