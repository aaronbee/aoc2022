package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/aaronbee/aoc2022/fn"
)

type stacks [][]byte

func (ss stacks) String() string {
	return string(bytes.Join(ss, []byte{'\n'}))
}

func (ss stacks) clone() stacks {
	return fn.Map(ss, func(s []byte) []byte {
		return append([]byte(nil), s...)
	})
}

func (ss stacks) reverse() {
	for _, s := range ss {
		fn.Reverse(s)
	}
}

func (ss stacks) tops() string {
	bs := fn.Map(ss, func(s []byte) byte {
		if len(s) == 0 {
			return ' '
		}
		return s[len(s)-1]
	})
	return string(bs)
}

type move struct {
	n   int
	src int
	dst int
}

type parsed struct {
	ss stacks
	ms []move
}

func shuffle(p parsed, reverse bool) {
	for _, m := range p.ms {
		srcI := len(p.ss[m.src]) - m.n
		moved := p.ss[m.src][srcI:]
		if reverse {
			fn.Reverse(moved)
		}
		p.ss[m.dst] = append(p.ss[m.dst], moved...)
		p.ss[m.src] = p.ss[m.src][:srcI]
	}
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	p := fn.ReduceReader(f, parsed{}, func(acc parsed, line string) parsed {
		switch {
		case strings.Contains(line, "["):
			nStacks := (len(line) + 1) / 4
			if acc.ss == nil {
				acc.ss = make(stacks, nStacks)
			}
			for i := 0; i < nStacks; i++ {
				index := i*4 + 1
				crate := line[index]
				if crate == ' ' {
					continue
				}
				acc.ss[i] = append(acc.ss[i], crate)
			}
		case strings.HasPrefix(line, "move"):
			var m move
			n, err := fmt.Sscanf(line, "move %d from %d to %d", &m.n, &m.src, &m.dst)
			if err != nil {
				panic(fmt.Errorf("error scanning move line: %s", err))
			} else if n != 3 {
				panic(fmt.Errorf("unexpected number of args from scanning %q: %d", line, n))
			}
			acc.ms = append(acc.ms, move{n: m.n, src: m.src - 1, dst: m.dst - 1})
		}
		return acc
	})
	p.ss.reverse()

	p1 := p
	p2 := parsed{ss: p1.ss.clone(), ms: p1.ms}
	shuffle(p1, true)
	shuffle(p2, false)

	fmt.Println("Part 1:", p1.ss.tops())
	fmt.Println("Part 2:", p2.ss.tops())
}
