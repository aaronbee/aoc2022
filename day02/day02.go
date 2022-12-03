package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aaronbee/aoc2022/fn"
)

type rps uint8

const (
	rock rps = iota
	paper
	scissors
)

type round struct {
	p1 rps
	p2 rps
}

func (r rps) String() string {
	switch r {
	case rock:
		return "rock"
	case paper:
		return "paper"
	case scissors:
		return "scissors"
	}
	panic(fmt.Errorf("unexpected rps value: %d", r))
}

func map1(s string) rps {
	switch s {
	case "A":
		return rock
	case "B":
		return paper
	case "C":
		return scissors
	}
	panic(fmt.Errorf("unexpected p1 play: %q", s))
}

func map2(s string) rps {
	switch s {
	case "X":
		return rock
	case "Y":
		return paper
	case "Z":
		return scissors
	}
	panic(fmt.Errorf("unexpected p2 play: %q", s))
}

func score(r round) int {
	var s int
	switch r.p2 {
	case rock:
		s += 1
	case paper:
		s += 2
	case scissors:
		s += 3
	}
	if r.p1 == r.p2 { // tie
		s += 3
	} else if (r.p1+1)%3 == r.p2 { // p2 wins
		s += 6
	}
	return s
}

func part1(r io.Reader) {
	rs := fn.ReduceReader(r, nil, func(acc []round, line string) []round {
		p1, p2, ok := strings.Cut(line, " ")
		if !ok {
			panic(fmt.Errorf("unexpected line: %q", line))
		}
		return append(acc, round{
			p1: map1(p1),
			p2: map2(p2),
		})
	})
	scores := fn.Map(rs, score)
	fmt.Println("Part 1:", fn.Sum(scores))
}

func mapStrat(p1 rps, strat string) rps {
	switch strat {
	case "X": // lose
		// avoid underflow by adding 2 instead of subtracting 1
		return (p1 + 2) % 3
	case "Y": // draw
		return p1
	case "Z": // win
		return (p1 + 1) % 3
	}
	panic(fmt.Errorf("unexpected strat: %q", strat))
}

func part2(r io.Reader) {
	rs := fn.ReduceReader(r, nil, func(acc []round, line string) []round {
		p1s, p2s, ok := strings.Cut(line, " ")
		if !ok {
			panic(fmt.Errorf("unexpected line: %q", line))
		}
		p1 := map1(p1s)
		p2 := mapStrat(p1, p2s)
		return append(acc, round{
			p1: p1,
			p2: p2,
		})
	})
	scores := fn.Map(rs, score)
	fmt.Println("Part 2:", fn.Sum(scores))
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	part1(f)
	f.Seek(0, 0)
	part2(f)
}
