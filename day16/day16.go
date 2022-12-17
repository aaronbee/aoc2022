package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aaronbee/aoc2022/fn"
	"golang.org/x/exp/slices"
)

type valve struct {
	name        string
	rate        int
	connections []*valve
}

func (v *valve) String() string {
	return fmt.Sprintf("{name: %s rate: %d connections: %v}",
		v.name, v.rate, fn.Map(v.connections, func(v *valve) string { return v.name }))
}

type path struct {
	v    *valve
	cost int
}

func shortestPathsFromV(v *valve, nonzeroValves map[*valve]struct{}, out map[*valve]int) {
	expectedConnections := len(nonzeroValves)
	if v.rate > 0 {
		expectedConnections--
	}
	seen := map[*valve]struct{}{v: {}}
	fringe := fn.Map(v.connections, func(v *valve) path { return path{v, 1} })
	for {
		p := fringe[0]
		fringe = fringe[1:]
		seen[p.v] = struct{}{}
		if p.v.rate > 0 {
			out[p.v] = p.cost
		}
		if len(out) == expectedConnections {
			return
		}
		for _, next := range p.v.connections {
			if _, ok := seen[next]; !ok {
				fringe = append(fringe, path{next, p.cost + 1})
			}
		}
	}
}

func shortestPaths(g map[string]*valve, nonzeroValves map[*valve]struct{}) map[*valve]map[*valve]int {
	out := map[*valve]map[*valve]int{}
	for _, v := range g {
		m, ok := out[v]
		if !ok {
			m = make(map[*valve]int)
			out[v] = m
		}
		shortestPathsFromV(v, nonzeroValves, m)
	}
	return out
}

type action struct {
	v        *valve
	timeleft int
}

func maxRelease(shortestPaths map[*valve]map[*valve]int, on []*valve, v *valve, timeleft int) (int, []action) {
	var myRelease int
	var myActions []action
	if v.rate > 0 {
		timeleft-- // time to turn on
		myRelease = v.rate * timeleft
		myActions = []action{{v, timeleft}}
	}
	best := myRelease
	actions := myActions

	for next, cost := range shortestPaths[v] {
		if next.rate == 0 {
			continue
		}
		if cost > timeleft {
			continue
		}
		if slices.Contains(on, next) {
			continue
		}
		release, acts := maxRelease(shortestPaths, append(on, next), next, timeleft-cost)
		totalRelease := myRelease + release
		if totalRelease > best {
			best = totalRelease
			actions = append(acts, myActions...)
		}
	}

	return best, actions
}

func part1(vs map[string]*valve, nonzeroValves map[*valve]struct{}) int {
	sps := shortestPaths(vs, nonzeroValves)
	rate, actions := maxRelease(sps, nil, vs["AA"], 30)
	fn.Reverse(actions)
	for _, a := range actions {
		fmt.Printf("At time %d turn on %s at rate %d for total %d\n", 30-a.timeleft, a.v.name, a.v.rate, a.timeleft*a.v.rate)
	}
	return rate
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	vs := fn.ReduceReader(f, map[string]*valve{}, func(acc map[string]*valve, line string) map[string]*valve {
		var rate int
		var name string
		n, err := fmt.Sscanf(line, "Valve %s has flow rate=%d;", &name, &rate)
		if err != nil || n != 2 {
			panic(fmt.Errorf("failed to parse line %q: n: %d err: %v", line, n, err))
		}
		v, ok := acc[name]
		if !ok {
			v = &valve{name: name}
			acc[name] = v
		}
		v.rate = rate
		i := strings.Index(line, "valve")
		i += len("valve")
		if line[i] == 's' {
			i++
		}
		i += 1 // skip space
		for _, c := range strings.Split(line[i:], ", ") {
			cv, ok := acc[c]
			if !ok {
				cv = &valve{name: c}
				acc[c] = cv
			}
			v.connections = append(v.connections, cv)
		}
		return acc
	})
	nonzeroValves := map[*valve]struct{}{}
	for _, v := range vs {
		if v.rate > 0 {
			nonzeroValves[v] = struct{}{}
		}
	}

	fmt.Println("Part 1:", part1(vs, nonzeroValves))
}
