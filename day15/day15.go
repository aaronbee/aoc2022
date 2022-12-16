package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aaronbee/aoc2022/fn"
)

type point struct{ x, y int }

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func dist(a, b point) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

type sensor struct {
	loc  point
	beac point
}

func (s sensor) intersection(yCoord int, m map[int]struct{}) {
	d := dist(s.loc, s.beac)
	distToLine := abs(yCoord - s.loc.y)
	leftOver := d - distToLine + 1
	if leftOver < 0 {
		return
	}
	m[s.loc.x] = struct{}{}
	for i := 0; i < leftOver; i++ {
		m[s.loc.x-i] = struct{}{}
		m[s.loc.x+i] = struct{}{}
	}
}

func part1(ss []sensor, yCoord int) int {
	intersects := make(map[int]struct{})
	for _, s := range ss {
		s.intersection(yCoord, intersects)
	}
	for _, s := range ss {
		if s.beac.y == yCoord {
			delete(intersects, s.beac.x)
		}
	}
	return len(intersects)
}

func part2(ss []sensor, maxC int) int {
	return 0
}

func main() {
	arg1 := flag.Int("y", 2000000, "part 1 argument")
	arg2 := flag.Int("box", 4000000, "part 2 argument")
	flag.Parse()
	f, err := os.Open(flag.Args()[0])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	ss := fn.ReduceReader(f, nil, func(acc []sensor, line string) []sensor {
		var s sensor
		n, err := fmt.Sscanf(line, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d",
			&s.loc.x, &s.loc.y, &s.beac.x, &s.beac.y)
		if err != nil || n != 4 {
			panic(fmt.Errorf("failed to parse line %q: %s", line, err))
		}
		return append(acc, s)
	})

	fmt.Println("Part 1:", part1(ss, *arg1))
	fmt.Println("Part 2:", part2(ss, *arg2))
}
