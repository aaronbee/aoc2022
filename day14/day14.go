package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aaronbee/aoc2022/fn"
)

const (
	air   = '.'
	rock  = '#'
	sand  = 'o'
	start = 'S'
)

type point struct{ x, y int }

func parsePoint(s string) point {
	xStr, yStr, ok := strings.Cut(s, ",")
	if !ok {
		panic(fmt.Errorf("failed to parse point: %q", s))
	}
	x, err := strconv.Atoi(xStr)
	if err != nil {
		panic(err)
	}
	y, err := strconv.Atoi(yStr)
	if err != nil {
		panic(err)
	}
	return point{x, y}
}

func (p point) add(other point) point {
	return point{p.x + other.x, p.y + other.y}
}

func (p point) line(end point, f func(p point)) {
	var dir point
	if p.x < end.x {
		dir = point{1, 0}
	} else if p.x > end.x {
		dir = point{-1, 0}
	} else if p.y < end.y {
		dir = point{0, 1}
	} else if p.y > end.y {
		dir = point{0, -1}
	}
	for p != end {
		p = p.add(dir)
		f(p)
	}
}

type cave struct {
	m          map[point]byte
	minX, maxX int
	maxY       int
	floor      int
}

func newCave() *cave {
	return &cave{
		m:    make(map[point]byte),
		minX: 500,
		maxX: 500,
	}
}

func min(x, y int) int {
	if x <= y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x >= y {
		return x
	}
	return y
}

func (c *cave) insert(p point, b byte) {
	c.m[p] = b
	c.minX = min(c.minX, p.x)
	c.maxX = max(c.maxX, p.x)
	if c.floor == 0 {
		c.maxY = max(c.maxY, p.y)
	}
}

func (c *cave) setFloor() {
	c.floor = c.maxY + 1
}

func (c *cave) String() string {
	width := 3 + c.maxX - c.minX
	offset := -c.minX + 1
	depth := 3 + c.maxY
	buf := make([][]byte, depth)
	for i := range buf {
		if i == depth-1 {
			buf[i] = bytes.Repeat([]byte{rock}, width)
			continue
		}
		buf[i] = bytes.Repeat([]byte{air}, width)
	}
	buf[0][500+offset] = 'S'
	for p, b := range c.m {
		buf[p.y][p.x+offset] = b
	}
	return string(bytes.Join(buf, []byte{'\n'}))
}

func (c *cave) dropSandOneStep(s point) point {
	if s.y == c.maxY+1 {
		return s
	}
	for _, next := range []point{s.add(point{0, 1}), s.add(point{-1, 1}), s.add(point{1, 1})} {
		if _, ok := c.m[next]; !ok {
			return next
		}
	}
	return s
}

func (c *cave) dropSand(stop func(s point) bool) bool {
	s := c.dropSandOneStep(point{500, 0})
	for !stop(s) {
		next := c.dropSandOneStep(s)
		if next == s {
			c.insert(next, sand)
			return true
		}
		s = next
	}
	return false
}

func part1(c *cave) int {
	i := 0
	for ; c.dropSand(func(s point) bool { return s.y >= c.maxY }); i++ {
	}
	return i
}

func part2(c *cave) int {
	i := 0
	for ; c.dropSand(func(s point) bool { return s == point{500, 0} }); i++ {
	}
	return i + 1
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	c := fn.ReduceReader(f, newCave(), func(c *cave, line string) *cave {
		ps := fn.Map(strings.Split(line, " -> "), parsePoint)
		cur := ps[0]
		c.insert(cur, rock)
		for _, p := range ps[1:] {
			cur.line(p, func(p point) { c.insert(p, rock) })
			cur = p
		}
		return c
	})
	c.setFloor()
	count := part1(c)
	fmt.Println("Part 1:", count)
	fmt.Println("Part 2:", count+part2(c))
}
