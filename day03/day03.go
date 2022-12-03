package main

import (
	"fmt"
	"io"
	"os"

	"github.com/aaronbee/aoc2022/fn"
	"golang.org/x/exp/slices"
)

type rucksack struct {
	c1 []byte
	c2 []byte
}

func (r rucksack) common() byte {
	m := fn.Slice2Set(r.c1)
	i := slices.IndexFunc(r.c2, fn.Map2Func(m))
	if i == -1 {
		panic(fmt.Errorf("no matches found between %q and %q", r.c1, r.c2))
	}
	return r.c2[i]
}

func priority(b byte) int {
	if 'a' <= b && b <= 'z' {
		return 1 + int(b-'a')
	} else if 'A' <= b && b <= 'Z' {
		return 27 + int(b-'A')
	}
	panic(fmt.Errorf("unexpected value: %q", b))
}

func part1(r io.Reader) int {
	rs := fn.ReduceReader(r, nil, func(acc []rucksack, lineS string) []rucksack {
		line := []byte(lineS)
		mid := len(line) / 2
		return append(acc, rucksack{
			c1: line[:mid],
			c2: line[mid:],
		})
	})
	prs := fn.Map(rs, func(r rucksack) int { return priority(r.common()) })
	return fn.Sum(prs)
}

type group [][]byte

func newGroup() group { return make(group, 0, 3) }

func (g group) common() byte {
	m1 := fn.Slice2Set(g[0])
	m2 := fn.Slice2Set(g[1])
	i := slices.IndexFunc(g[2], func(b byte) bool { return m1[b] && m2[b] })
	if i == -1 {
		panic(fmt.Errorf("can't find common element in %q %q %q", g[0], g[1], g[2]))
	}
	return g[2][i]
}

func part2(r io.Reader) int {
	gs := fn.ReduceReader(r, []group{newGroup()}, func(acc []group, line string) []group {
		i := len(acc) - 1
		if len(acc[i]) == 3 {
			acc = append(acc, newGroup())
			i++
		}
		acc[i] = append(acc[i], []byte(line))
		return acc
	})
	prs := fn.Map(gs, func(g group) int { return priority(g.common()) })
	return fn.Sum(prs)
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 1:", part1(f))
	f.Seek(0, 0)
	fmt.Println("Part 2:", part2(f))
}
