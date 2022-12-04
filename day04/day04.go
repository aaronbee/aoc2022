package main

import (
	"fmt"
	"os"

	"github.com/aaronbee/aoc2022/fn"
)

type rnge struct {
	begin, end int
}

type pair struct {
	e1 rnge
	e2 rnge
}

func (p pair) overlap() bool {
	return (p.e1.begin <= p.e2.begin && p.e1.end >= p.e2.end) ||
		(p.e2.begin <= p.e1.begin && p.e2.end >= p.e1.end)
}

func (p pair) overlapPartial() bool {
	r1, r2 := p.e1, p.e2
	if r2.begin < r1.begin {
		r1, r2 = r2, r1
	}
	return r1.end >= r2.begin
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	pairs := fn.ReduceReader(f, nil, func(acc []pair, line string) []pair {
		var e1b, e1e, e2b, e2e int
		n, err := fmt.Sscanf(line, "%d-%d,%d-%d", &e1b, &e1e, &e2b, &e2e)
		if err != nil {
			panic(err)
		} else if n != 4 {
			panic(fmt.Errorf("unexpected line: %q", line))
		}
		return append(acc, pair{
			e1: rnge{e1b, e1e},
			e2: rnge{e2b, e2e},
		})
	})

	fmt.Println("Part 1:", fn.Count(pairs, pair.overlap))
	fmt.Println("Part 2:", fn.Count(pairs, pair.overlapPartial))
}
