package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aaronbee/aoc2022/fn"
	"golang.org/x/exp/slices"
)

type rock uint16 // 4x4 bitmap, bottom-left justified

const (
	none rock = 0

	// ####
	flat rock = 0b0000_0000_0000_1111

	// .#.
	// ###
	// .#.
	cross rock = 0b0000_0100_1110_0100

	// ..#
	// ..#
	// ###
	elbow rock = 0b0000_0010_0010_1110

	// #
	// #
	// #
	// #
	tall rock = 0b1000_1000_1000_1000

	// ##
	// ##
	square rock = 0b0000_0000_1100_1100

	width int = 7
)

var order = [...]rock{flat, cross, elbow, tall, square}

func (s rock) String() string {
	switch s {
	case none:
		return "none"
	case flat:
		return "flat"
	case cross:
		return "cross"
	case elbow:
		return "elbow"
	case tall:
		return "tall"
	case square:
		return "square"
	}
	return fmt.Sprintf("rock(0b%b)", s)
}

func (s rock) height() int {
	switch {
	case s&0xF000 > 0:
		return 4
	case s&0x0F00 > 0:
		return 3
	case s&0x00F0 > 0:
		return 2
	case s&0x000F > 0:
		return 1
	default:
		return 0
	}
}

func (s rock) width() int {
	switch {
	case s&0x1111 > 0:
		return 4
	case s&0x2222 > 0:
		return 3
	case s&0x4444 > 0:
		return 2
	case s&0x8888 > 0:
		return 1
	default:
		return 0
	}
}

func collides(a rock, ax, ay int, b rock, bx, by int) bool {
	if a == none || b == none {
		return false
	}
	switch bx - ax {
	case 0:
	case 1:
		b = (b &^ 0x1111) >> 1
	case 2:
		b = (b &^ 0x3333) >> 2
	case 3:
		b = (b &^ 0x7777) >> 3
	case -1:
		a = (a &^ 0x1111) >> 1
	case -2:
		a = (a &^ 0x3333) >> 2
	case -3:
		a = (a &^ 0x7777) >> 3
	default:
		return false
	}
	switch by - ay {
	case 0:
	case 1:
		b <<= 4
	case 2:
		b <<= 8
	case 3:
		b <<= 12
	case -1:
		a <<= 4
	case -2:
		a <<= 8
	case -3:
		a <<= 12
	default:
		return false
	}
	return a&b > 0
}

func (s rock) writeTo(posX int, char byte, b [][]byte) {
	if s == none {
		return
	}
	for i := 0; i < 4; i++ {
		slice := s >> (i * 4)
		for j := 0; j < 4; j++ {
			mask := rock(1 << (3 - j))
			if slice&mask > 0 {
				b[i][j+posX] = char
			}
		}
	}
}

type chamber struct {
	// stable rocks
	rocks []*[width]rock

	falling rock
	x, y    int // coord for falling rock
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (c *chamber) highestOccupied() int {
	ret := c.falling.height() + c.y
	row := max(0, len(c.rocks)-4)
	for ; row < len(c.rocks); row++ {
		rocks := c.rocks[row]
		for _, rock := range rocks {
			ret = max(ret, rock.height()+row)
		}
	}
	return ret
}

func (c *chamber) String() string {
	const (
		floor    = "+-------+"
		nonFloor = "|.......|"
	)
	height := c.highestOccupied() + 1
	buf := make([][]byte, height)
	buf[0] = []byte(floor)
	for i := range buf[1:] {
		buf[i+1] = []byte(nonFloor)
	}
	for i, row := range c.rocks {
		for x, rock := range row {
			rock.writeTo(x+1, '#', buf[i+1:])
		}
	}
	c.falling.writeTo(c.x+1, '@', buf[c.y+1:])
	fn.Reverse(buf)
	return string(bytes.Join(buf, []byte{'\n'}))
}

func (c *chamber) collides(r rock, x, y int) bool {
	if x < 0 || y < 0 || r.width()+x > width {
		return true
	}
	for row := max(0, y-4); row < y+4 && row < len(c.rocks); row++ {
		for col := max(0, x-4); col < x+4 && col < width; col++ {
			if collides(c.rocks[row][col], col, row, r, x, y) {
				return true
			}
		}
	}
	return false
}

func (c *chamber) drop(r rock, pushes []byte, pushIndex int) int {
	c.y = c.highestOccupied() + 3
	c.falling = r
	c.x = 2
	defer func() {
		c.falling = none
		c.x = 0
		c.y = 0
	}()
	for {
		switch pushes[pushIndex] {
		case '<':
			if !c.collides(c.falling, c.x-1, c.y) {
				c.x--
			}
		case '>':
			if !c.collides(c.falling, c.x+1, c.y) {
				c.x++
			}
		default:
			panic(fmt.Errorf("unexpected push: %s", string(pushes[pushIndex])))
		}
		pushIndex = (pushIndex + 1) % len(pushes)
		if c.collides(c.falling, c.x, c.y-1) {
			// rock stops
			if c.y >= len(c.rocks) {
				c.rocks = slices.Grow(c.rocks, c.y+1-len(c.rocks))
			}
			for i := len(c.rocks); i < c.y+1; i++ {
				c.rocks = append(c.rocks, &[width]rock{})
			}
			c.rocks[c.y][c.x] = c.falling
			return pushIndex
		}
		c.y--
	}
}

func part1(pushes []byte) int {
	c := chamber{}
	var pushIndex int
	for i := 0; i < 2022; i++ {
		pushIndex = c.drop(order[i%len(order)], pushes, pushIndex)
	}
	return c.highestOccupied()
}

func main() {
	pushes, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	pushes = bytes.TrimSpace(pushes)
	fmt.Println("Part 1:", part1(pushes))
}
