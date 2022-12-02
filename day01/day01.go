package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aaronbee/aoc2022/fn"

	"golang.org/x/exp/slices"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	ns := fn.ReduceReader(f, [][]int{nil}, func(acc [][]int, line string) [][]int {
		if line == "" {
			acc = append(acc, []int(nil))
			return acc
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}
		i := len(acc) - 1
		acc[i] = append(acc[i], n)
		return acc
	})

	sums := fn.Map(ns, fn.Sum[int])
	slices.Sort(sums)
	fmt.Println("Part 1:", sums[len(sums)-1])
	fmt.Println("Part 2:", fn.Sum(sums[len(sums)-3:]))
}
