package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aaronbee/aoc2022/fn"
)

// a value of 0 is a noop, else its an addx of the value.
type instruction int

func (i instruction) String() string {
	if i == 0 {
		return "noop"
	}
	return "addx " + strconv.Itoa(int(i))
}

func (i instruction) cycles() int {
	if i == 0 {
		return 1
	}
	return 2
}

type cpu struct {
	insn       []instruction
	pc         int
	cyclesLeft int
	x          int
}

func (c *cpu) reset() {
	c.pc = 0
	c.cyclesLeft = c.insn[0].cycles() - 1
	c.x = 1
}

func (c *cpu) run(n int) int {
	cur := c.x
	for i := 0; i < n; i++ {
		cur = c.x
		if c.cyclesLeft == 0 {
			c.x += int(c.insn[c.pc])
			c.pc++
			if c.pc == len(c.insn) {
				return cur
			}
			c.cyclesLeft = c.insn[c.pc].cycles()
		}
		c.cyclesLeft--
	}
	return cur
}

func part1(c *cpu) int {
	var count int
	var sum int
	for _, i := range []int{20, 40, 40, 40, 40, 40} {
		count += i
		x := c.run(i)
		sum += count * x
	}
	return sum
}

func part2(c *cpu) *strings.Builder {
	var buf strings.Builder
	for i := 0; i < 240; i++ {
		pos := i % 40
		if pos == 0 && i > 0 {
			buf.WriteByte('\n')
		}
		x := c.run(1)
		if x-1 <= pos && pos <= x+1 {
			buf.WriteByte('#')
		} else {
			buf.WriteByte(' ')
		}
	}
	return &buf
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	insns := fn.ReduceReader(f, nil, func(acc []instruction, line string) []instruction {
		if line == "noop" {
			return append(acc, instruction(0))
		}
		insn, nStr, ok := strings.Cut(line, " ")
		if !ok || insn != "addx" {
			panic(fmt.Errorf("failed to parse instruction: %q", line))
		}
		n, err := strconv.Atoi(nStr)
		if err != nil {
			panic(fmt.Errorf("failed to parse instruction: %q", line))
		} else if n == 0 {
			panic(fmt.Errorf("unsupported addx value: %d", n))
		}
		return append(acc, instruction(n))
	})
	c := cpu{insn: insns}
	c.reset()
	fmt.Println("Part 1:", part1(&c))
	c.reset()
	fmt.Println(part2(&c))
}
