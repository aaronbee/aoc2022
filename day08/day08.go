package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aaronbee/aoc2022/fn"
)

type grid [][]int

func (g grid) String() string {
	buf := fn.Reduce(g, &strings.Builder{}, func(buf *strings.Builder, ts []int) *strings.Builder {
		strs := fn.Map(ts, strconv.Itoa)
		buf.WriteString(strings.Join(strs, " "))
		buf.WriteByte('\n')
		return buf
	})
	return buf.String()
}

func (g grid) visibleGrid() grid {
	visible := grid(fn.Map(g, func(ts []int) []int { return make([]int, len(ts)) }))
	// from left
	for y := 0; y < len(g); y++ {
		cur := -1
		for x := 0; x < len(g[y]); x++ {
			if t := g[y][x]; t > cur {
				visible[y][x] = 1
				cur = g[y][x]
			}
		}
	}
	// from right
	for y := 0; y < len(g); y++ {
		cur := -1
		for x := len(g[y]) - 1; x >= 0; x-- {
			if t := g[y][x]; t > cur {
				visible[y][x] = 1
				cur = t
			}
		}
	}
	// from top
	for x := 0; x < len(g[0]); x++ {
		cur := -1
		for y := 0; y < len(g); y++ {
			if t := g[y][x]; t > cur {
				visible[y][x] = 1
				cur = t
			}
		}
	}
	// from bottom
	for x := 0; x < len(g[0]); x++ {
		cur := -1
		for y := len(g) - 1; y >= 0; y-- {
			if t := g[y][x]; t > cur {
				visible[y][x] = 1
				cur = t
			}
		}
	}
	return visible
}

func (g grid) scenicScores() grid {
	out := grid(fn.Map(g, func(ts []int) []int { return make([]int, len(ts)) }))
	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[y]); x++ {
			out[y][x] = g.scenicScore(y, x)
		}
	}
	return out
}

func (g grid) scenicScore(y, x int) int {
	cur := g[y][x]
	var right, left, down, up int
	for x+right+1 < len(g[y]) {
		right++
		if g[y][x+right] >= cur {
			break
		}
	}
	for x-(left+1) >= 0 {
		left++
		if g[y][x-left] >= cur {
			break
		}
	}
	for y+down+1 < len(g) {
		down++
		if g[y+down][x] >= cur {
			break
		}
	}
	for y-(up+1) >= 0 {
		up++
		if g[y-up][x] >= cur {
			break
		}
	}

	return right * left * down * up
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	g := fn.ReduceReader(f, nil, func(acc grid, line string) grid {
		return append(acc, fn.Map([]byte(line), func(b byte) int {
			return int(b - '0')
		}))
	})
	visible := g.visibleGrid()
	fmt.Println("Part 1:", fn.Sum(fn.Map(visible, fn.Sum[int])))

	scores := g.scenicScores()
	fmt.Println("Part 2:", fn.Max(fn.Map(scores, fn.Max[int])))
}
