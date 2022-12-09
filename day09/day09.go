package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aaronbee/aoc2022/fn"
)

type direction byte

const (
	right direction = 'R'
	down  direction = 'D'
	left  direction = 'L'
	up    direction = 'U'
)

func (d direction) toCoord() coordinate {
	switch d {
	case right:
		return coordinate{1, 0}
	case down:
		return coordinate{0, -1}
	case left:
		return coordinate{-1, 0}
	case up:
		return coordinate{0, 1}
	}
	panic(fmt.Errorf("unexpected coordinate: %q", d))
}

type motion struct {
	dir direction
	n   int
}

type coordinate struct {
	x, y int
}

func (c coordinate) add(d coordinate) coordinate {
	return coordinate{c.x + d.x, c.y + d.y}
}

func oneCloser(head, tail int) int {
	if tail < head {
		return tail + 1
	} else if tail > head {
		return tail - 1
	}
	return tail
}

func follow(head, tail coordinate) coordinate {
	if head.x-tail.x > 1 {
		tail.x++
		tail.y = oneCloser(head.y, tail.y)
	} else if head.x-tail.x < -1 {
		tail.x--
		tail.y = oneCloser(head.y, tail.y)
	} else if head.y-tail.y > 1 {
		tail.y++
		tail.x = oneCloser(head.x, tail.x)
	} else if head.y-tail.y < -1 {
		tail.y--
		tail.x = oneCloser(head.x, tail.x)
	}
	return tail
}

func part1(ms []motion) int {
	seen := map[coordinate]struct{}{}
	head := coordinate{0, 0}
	tail := coordinate{0, 0}
	seen[tail] = struct{}{}
	for _, m := range ms {
		d := m.dir.toCoord()
		for i := 0; i < m.n; i++ {
			head = head.add(d)
			tail = follow(head, tail)
			seen[tail] = struct{}{}
		}
	}
	return len(seen)
}

func part2(ms []motion) int {
	seen := map[coordinate]struct{}{}
	knots := make([]coordinate, 10)
	seen[knots[9]] = struct{}{}
	for _, m := range ms {
		d := m.dir.toCoord()
		for i := 0; i < m.n; i++ {
			knots[0] = knots[0].add(d)
			for i := 1; i < len(knots); i++ {
				knots[i] = follow(knots[i-1], knots[i])
			}
			seen[knots[9]] = struct{}{}
		}
	}
	return len(seen)
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	ms := fn.ReduceReader(f, nil, func(acc []motion, line string) []motion {
		dStr, nStr, ok := strings.Cut(line, " ")
		if !ok {
			panic(fmt.Errorf("failed to cut line: %q", line))
		}
		n, err := strconv.Atoi(nStr)
		if err != nil {
			panic(fmt.Errorf("failed to parse n: %q", line))
		}
		return append(acc, motion{dir: direction(dStr[0]), n: n})
	})

	fmt.Println("Part 1:", part1(ms))
	fmt.Println("Part 2:", part2(ms))
}
