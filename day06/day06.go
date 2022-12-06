package main

import (
	"bytes"
	"fmt"
	"os"
)

func uniq(byts []byte) bool {
	m := make(map[byte]bool, len(byts))
	for _, b := range byts {
		if m[b] {
			return false
		}
		m[b] = true
	}
	return true
}

func uniqSeq(n int, byts []byte) int {
	for i := 0; i < len(byts)-n; i++ {
		if uniq(byts[i : i+n]) {
			return i + n
		}
	}
	return -1
}

func main() {
	byts, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	byts = bytes.TrimSpace(byts)

	fmt.Println("Part 1:", uniqSeq(4, byts))
	fmt.Println("Part 2:", uniqSeq(14, byts))
}
