package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

type node interface {
	isDir() bool
	size() int
}

type dir struct {
	sizeCache int
	children  map[string]node
	p         *dir
}

func (d *dir) isDir() bool { return true }
func (d *dir) size() int {
	if d.sizeCache != 0 {
		return d.sizeCache
	}
	for _, n := range d.children {
		d.sizeCache += n.size()
	}
	return d.sizeCache
}
func (d *dir) parent() *dir { return d.p }
func (d *dir) chdir(name string) *dir {
	return d.children[name].(*dir)
}

type file struct {
	s int
}

func (f *file) isDir() bool { return false }
func (f *file) size() int   { return f.s }

func walk(n node, f func(n node)) {
	f(n)
	if n.isDir() {
		for _, c := range n.(*dir).children {
			walk(c, f)
		}
	}
}

func splitCmd(data []byte, atEOF bool) (advance int, token []byte, err error) {
	i := bytes.Index(data, []byte("\n$"))
	if i == -1 {
		if atEOF && len(data) > 0 {
			return len(data), bytes.TrimSpace(data), nil
		}
		return
	}
	return i + 1, data[:i], nil
}

func parse(r io.Reader) *dir {
	root := &dir{children: make(map[string]node)}
	var cwd *dir
	s := bufio.NewScanner(r)
	s.Split(splitCmd)
	for s.Scan() {
		switch {
		case bytes.HasPrefix(s.Bytes(), []byte("$ cd")):
			d := string(bytes.TrimPrefix(s.Bytes(), []byte("$ cd ")))
			switch d {
			case "/":
				cwd = root
			case "..":
				cwd = cwd.parent()
			default:
				cwd = cwd.chdir(d)
			}
		case bytes.HasPrefix(s.Bytes(), []byte("$ ls")):
			ls := bytes.Split(s.Bytes(), []byte("\n"))[1:]
			for _, l := range ls {
				switch {
				case bytes.HasPrefix(l, []byte("dir")):
					d := string(bytes.TrimPrefix(l, []byte("dir ")))
					newDir := &dir{p: cwd, children: make(map[string]node)}
					cwd.children[d] = newDir
				default:
					sizeBytes, nameBytes, ok := bytes.Cut(l, []byte(" "))
					if !ok {
						panic(fmt.Errorf("failed to parse file listing: %q", l))
					}
					size, err := strconv.Atoi(string(sizeBytes))
					if err != nil {
						panic(fmt.Errorf("failed to parse size from: %q", l))
					}
					name := string(nameBytes)
					f := &file{s: size}
					cwd.children[name] = f
				}
			}
		}
	}
	return root
}

func part1(root *dir) int {
	const smallDirSize = 100000
	var sizeOfSmallDirs int
	walk(root, func(n node) {
		if n.isDir() && n.size() <= smallDirSize {
			sizeOfSmallDirs += n.size()
		}
	})
	return sizeOfSmallDirs
}

func part2(root *dir) int {
	const total = 70000000
	const need = 30000000
	have := total - root.size()
	want := need - have
	smallestDelete := root.size()
	walk(root, func(n node) {
		if n.isDir() && n.size() >= want && n.size() < smallestDelete {
			smallestDelete = n.size()
		}
	})
	return smallestDelete
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	root := parse(f)

	fmt.Println("Part 1:", part1(root))
	fmt.Println("Part 2:", part2(root))
}
