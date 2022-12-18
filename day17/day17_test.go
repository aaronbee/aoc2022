package main

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	for _, tc := range []struct {
		c   chamber
		exp string
	}{{
		c: chamber{rocks: []*[width]rock{
			{none, none, none, none, none, none, none},
		}},
		exp: "+-------+",
	}, {
		c: chamber{rocks: []*[width]rock{
			{flat, none, none, none, none, none, none},
		}},
		exp: "|####...|\n" +
			"+-------+",
	}, {
		c: chamber{rocks: []*[width]rock{
			{none, flat, none, none, none, none, none},
		}},
		exp: "|.####..|\n" +
			"+-------+",
	}, {
		c: chamber{rocks: []*[width]rock{
			{none, none, none, flat, none, none, none},
		}},
		exp: "|...####|\n" +
			"+-------+",
	}, {
		c: chamber{rocks: []*[width]rock{
			{square, none, none, flat, none, none, none},
		}},
		exp: "|##.....|\n" +
			"|##.####|\n" +
			"+-------+",
	}, {
		c: chamber{rocks: []*[width]rock{
			{tall, none, none, flat, none, none, none},
		}},
		exp: "|#......|\n" +
			"|#......|\n" +
			"|#......|\n" +
			"|#..####|\n" +
			"+-------+",
	}, {
		c: chamber{rocks: []*[width]rock{
			{none, tall, none, flat, none, none, none},
		}},
		exp: "|.#.....|\n" +
			"|.#.....|\n" +
			"|.#.....|\n" +
			"|.#.####|\n" +
			"+-------+",
	}, {
		c: chamber{rocks: []*[width]rock{
			{none, none, tall, flat, none, none, none},
		}},
		exp: "|..#....|\n" +
			"|..#....|\n" +
			"|..#....|\n" +
			"|..#####|\n" +
			"+-------+",
	}, {
		c: chamber{rocks: []*[width]rock{
			{none, none, tall, flat, none, none, none},
			{none, none, none, cross, none, none, none},
		}},
		exp: "|..#.#..|\n" +
			"|..####.|\n" +
			"|..#.#..|\n" +
			"|..#####|\n" +
			"+-------+",
	}, {
		c: chamber{
			rocks: []*[width]rock{
				{none, none, tall, flat, none, none, none},
				{none, none, none, cross, none, none, none},
			},
			falling: elbow,
			y:       5,
		},
		exp: "|..@....|\n" +
			"|..@....|\n" +
			"|@@@....|\n" +
			"|.......|\n" +
			"|..#.#..|\n" +
			"|..####.|\n" +
			"|..#.#..|\n" +
			"|..#####|\n" +
			"+-------+",
	}} {
		got := tc.c.String()
		t.Logf("\n%s", got)
		if got != tc.exp {
			t.Errorf("Expected:\n%s\nGot:\n%s", tc.exp, got)
		}
	}
}

func TestCollides(t *testing.T) {
	for _, tc := range []struct {
		a      rock
		ax, ay int
		b      rock
		bx, by int
		exp    bool
	}{{
		a:   tall,
		b:   flat,
		exp: true,
	}, {
		a:   tall,
		b:   flat,
		bx:  1,
		exp: false,
	}, {
		a:   tall,
		b:   flat,
		by:  1,
		exp: true,
	}, {
		a:   flat,
		b:   tall,
		ax:  1,
		exp: false,
	}, {
		a:   flat,
		b:   tall,
		ay:  1,
		exp: true,
	}, {
		a:   tall,
		b:   flat,
		ax:  1,
		exp: true,
	}, {
		a:   tall,
		b:   flat,
		ay:  1,
		exp: false,
	}, {
		a:   tall,
		b:   cross,
		by:  2,
		exp: true,
	}, {
		a:   tall,
		b:   cross,
		by:  3,
		exp: false,
	}} {
		name := fmt.Sprintf("a:%s_ax:%d_ay:%d_b:%s_bx:%d_by:%d", tc.a, tc.ax, tc.ay, tc.b, tc.bx, tc.by)
		t.Run(name, func(t *testing.T) {
			got := collides(tc.a, tc.ax, tc.ay, tc.b, tc.bx, tc.by)
			if got != tc.exp {
				t.Errorf("Exp: %t Got: %t", tc.exp, got)
			}
		})
	}
}
