package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aaronbee/aoc2022/fn"
	"golang.org/x/exp/slices"
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

type interval struct{ beg, end int }

type intervals struct {
	is []interval
}

func (is *intervals) add(r interval) {
	if len(is.is) == 0 {
		is.is = append(is.is, r)
		return
	}
	pos, found := slices.BinarySearchFunc(is.is, r, func(a, b interval) int { return a.beg - b.beg })
	var cur *interval

	if found {
		cur = &is.is[pos]
		if cur.end < r.end {
			cur.end = r.end
		} else {
			return
		}
	} else {
		is.is = slices.Insert(is.is, pos, r)
		if pos > 0 {
			pos--
		}
		cur = &is.is[pos]
	}

	var toRemove int
	for i := pos + 1; i < len(is.is); i++ {
		next := is.is[i]
		if cur.end+1 < next.beg {
			break
		}
		toRemove++
		if cur.end < next.end {
			cur.end = next.end
		}
	}
	is.is = slices.Delete(is.is, pos+1, pos+1+toRemove)
}

func (is intervals) sum() int {
	return fn.Reduce(is.is, 0, func(acc int, i interval) int {
		return acc + i.end - i.beg + 1
	})
}

type sensor struct {
	loc  point
	beac point
}

func (s sensor) intersectionInterval(yCoord int) (interval, bool) {
	d := dist(s.loc, s.beac)
	distToLine := abs(yCoord - s.loc.y)
	leftOver := d - distToLine
	if leftOver < 0 {
		return interval{}, false
	}
	return interval{s.loc.x - leftOver, s.loc.x + leftOver}, true
}

func part1(ss []sensor, yCoord int) int {
	var is intervals
	for _, s := range ss {
		if i, ok := s.intersectionInterval(yCoord); ok {
			is.add(i)
		}
	}
	overlap := map[int]struct{}{}
	for _, s := range ss {
		if s.beac.y == yCoord {
			overlap[s.beac.x] = struct{}{}
		}
	}
	return is.sum() - len(overlap)
}

func part2(ss []sensor, maxC int) int {
	for y := 0; y <= maxC; y++ {
		var is intervals
		for _, s := range ss {
			if i, ok := s.intersectionInterval(y); ok {
				if i.beg > maxC || i.end < 0 {
					continue
				}
				if i.beg < 0 {
					i.beg = 0
				}
				if i.end > maxC {
					i.end = maxC
				}
				is.add(i)
			}
		}
		if len(is.is) == 2 {
			x := is.is[0].end + 1
			return x*4000000 + y
		}
	}
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
