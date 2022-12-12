package main

import (
	"fmt"
	"os"

	"github.com/aaronbee/aoc2022/fn"
)

type coord struct{ x, y int }

func (c coord) up() coord    { return coord{c.x, c.y - 1} }
func (c coord) down() coord  { return coord{c.x, c.y + 1} }
func (c coord) left() coord  { return coord{c.x - 1, c.y} }
func (c coord) right() coord { return coord{c.x + 1, c.y} }
func (c coord) all() []coord { return []coord{c.up(), c.right(), c.down(), c.left()} }

type grid [][]uint8

func (g grid) get(c coord) uint8 {
	if c.y < 0 || c.y >= len(g) {
		return 255
	}
	if c.x < 0 || c.x >= len(g[c.y]) {
		return 255
	}
	return g[c.y][c.x]
}

func (g grid) startEnd() (start coord, end coord) {
	for y := range g {
		for x := range g[y] {
			if g[y][x] == 'S' {
				g[y][x] = 'a'
				start = coord{x, y}
			}
			if g[y][x] == 'E' {
				g[y][x] = 'z'
				end = coord{x, y}
			}
		}
	}

	return start, end
}

type path struct {
	loc   coord
	steps int
}

func nextSteps(g grid, seen map[coord]int, cur path) []path {
	ps := make([]path, 0, 4)
	curElev := g.get(cur.loc)
	for _, c := range cur.loc.all() {
		if g.get(c)-1 > curElev { // too high
			continue
		}
		if steps, ok := seen[c]; !ok || cur.steps+1 < steps {
			seen[c] = cur.steps + 1
		} else {
			continue
		}
		ps = append(ps, path{c, cur.steps + 1})
	}
	return ps
}

func shortestPath(g grid, seen map[coord]int, start, end coord) (int, bool) {
	seen[start] = 0
	fringe := nextSteps(g, seen, path{start, 0})
	for len(fringe) > 0 {
		cur := fringe[0]
		fringe = fringe[1:]
		if cur.loc == end {
			return cur.steps, true
		}
		fringe = append(fringe, nextSteps(g, seen, cur)...)
	}
	return 0, false
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	g := fn.ReduceReader(f, nil, func(acc grid, line string) grid {
		return append(acc, []uint8(line))
	})
	start, end := g.startEnd()
	seen := map[coord]int{}
	part1, ok := shortestPath(g, seen, start, end)
	if !ok {
		panic("I'm lost")
	}
	fmt.Println("Part 1:", part1)

	part2 := part1
	for y := range g {
		for x := range g[y] {
			c := coord{x, y}
			if g.get(c) != 'a' || c == start {
				continue
			}
			count, ok := shortestPath(g, seen, c, end)
			if ok && count < part2 {
				part2 = count
			}
		}
	}
	fmt.Println("Part 2:", part2)
}
