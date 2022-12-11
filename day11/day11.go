package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aaronbee/aoc2022/fn"
)

type monkey struct {
	items []int
	op    func(old int) int
	test  int
	t     int
	f     int
	count int
}

func (m *monkey) clone() *monkey {
	mm := *m
	mm.items = append([]int(nil), m.items...)
	return &mm
}

func runRound(ms []*monkey, div3 bool, mod int) {
	for _, m := range ms {
		for _, it := range m.items {
			m.count++
			n := m.op(it) % mod // prevent overflow with '% mod'
			if n < 0 {
				panic(fmt.Errorf("overflow: Before: %d after: %d", it, n))
			}
			if div3 {
				n /= 3
			}
			if n%m.test == 0 {
				ms[m.t].items = append(ms[m.t].items, n)
			} else {
				ms[m.f].items = append(ms[m.f].items, n)
			}
		}
		m.items = m.items[:0]
	}
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func parseOperation(s string) func(old int) int {
	suffix := strings.TrimPrefix(s, "  Operation: new = old ")
	op, vStr, ok := strings.Cut(suffix, " ")
	if !ok {
		panic(fmt.Errorf("failed to parse: %q", suffix))
	}
	if vStr == "old" {
		switch op {
		case "*":
			return func(old int) int { return old * old }
		case "+":
			return func(old int) int { return old + old }
		default:
			panic(fmt.Errorf("failed to parse: %q", s))
		}
	}
	v := atoi(vStr)
	switch op {
	case "*":
		return func(old int) int { return old * v }
	case "+":
		return func(old int) int { return old + v }
	default:
		panic(fmt.Errorf("failed to parse: %q", s))
	}
}

func parseMonkey(b string) *monkey {
	var m monkey
	lines := strings.Split(b, "\n")
	m.items = fn.Map(
		strings.Split(strings.TrimPrefix(lines[1], "  Starting items: "), ", "),
		atoi,
	)
	m.op = parseOperation(lines[2])
	m.test = atoi(strings.TrimPrefix(lines[3], "  Test: divisible by "))
	m.t = atoi(strings.TrimPrefix(lines[4], "    If true: throw to monkey "))
	m.f = atoi(strings.TrimPrefix(lines[5], "    If false: throw to monkey "))
	return &m
}

func splitMonkey(data []byte, atEOF bool) (advance int, token []byte, err error) {
	i := bytes.Index(data, []byte("\n\n"))
	if i == -1 {
		if atEOF && len(data) > 0 {
			return len(data), bytes.TrimSpace(data), nil
		}
		return
	}
	return i + 2, data[:i], nil
}

func part1(ms []*monkey, mod int) int {
	for i := 0; i < 20; i++ {
		runRound(ms, true, mod)
	}
	counts := fn.Map(ms, func(m *monkey) int { return m.count })
	top2 := fn.TopN(counts, 2)
	return top2[0] * top2[1]
}

func part2(ms []*monkey, mod int) int {
	for i := 0; i < 10000; i++ {
		runRound(ms, false, mod)
	}
	counts := fn.Map(ms, func(m *monkey) int { return m.count })
	top2 := fn.TopN(counts, 2)
	return top2[0] * top2[1]
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	s.Split(splitMonkey)
	var ms []*monkey
	for s.Scan() {
		ms = append(ms, parseMonkey(s.Text()))
	}
	if s.Err() != nil {
		panic(s.Err())
	}
	// For all values of n and m.test: n % m.test == (n % mod) % m.test
	mod := fn.Reduce(ms, 1, func(acc int, m *monkey) int { return acc * m.test })
	fmt.Println("Part 1:", part1(fn.Map(ms, func(m *monkey) *monkey { return m.clone() }), mod))
	fmt.Println("Part 2:", part2(ms, mod))
}
