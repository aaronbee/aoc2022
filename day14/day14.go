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
	if c.floor == 0 {
		c.minX = min(c.minX, p.x)
		c.maxX = max(c.maxX, p.x)
		c.maxY = max(c.maxY, p.y)
	}
}

func (c *cave) setFloor() {
	c.floor = c.maxY + 2
}

func (c *cave) String() string {
	width := 3 + c.maxX - c.minX
	offset := -c.minX + 1
	depth := c.floor + 1
	buf := make([][]byte, depth)
	for i := range buf {
		if i == c.floor {
			buf[i] = bytes.Repeat([]byte{rock}, width)
			continue
		}
		buf[i] = bytes.Repeat([]byte{air}, width)
	}
	buf[0][500+offset] = 'S'
	for p, b := range c.m {
		x := p.x + offset
		if 0 <= x && x < width {
			buf[p.y][x] = b
		}
	}
	return string(bytes.Join(buf, []byte{'\n'}))
}

func (c *cave) dropSandGrain(s point) point {
	for {
		if s.y == c.floor-1 {
			break
		}
		next := s.add(point{0, 1})
		if _, ok := c.m[next]; !ok {
			s = next
			continue
		}
		next = s.add(point{-1, 1})
		if _, ok := c.m[next]; !ok {
			s = next
			continue
		}
		next = s.add(point{1, 1})
		if _, ok := c.m[next]; !ok {
			s = next
			continue
		}
		break
	}
	c.insert(s, sand)
	return s
}

func part1(c *cave) int {
	count := 0
	start := point{500, 0}
	for {
		grain := c.dropSandGrain(start)
		count++
		if grain.y >= c.maxY {
			break
		}
	}
	return count
}

func part2(c *cave) int {
	count := 0
	start := point{500, 0}
	for {
		grain := c.dropSandGrain(start)
		count++
		if grain == start {
			break
		}
	}
	return count
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
	fmt.Println("Part 1:", count-1) // don't count the grain that fell to the floor
	fmt.Println("Part 2:", count+part2(c))
}
