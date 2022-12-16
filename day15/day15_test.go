package main

import (
	"fmt"
	"testing"

	"golang.org/x/exp/slices"
)

func TestInterval(t *testing.T) {
	for _, tc := range []struct {
		is  []interval
		exp []interval
	}{{
		is:  []interval{{1, 2}, {3, 4}},
		exp: []interval{{1, 4}},
	}, {
		is:  []interval{{1, 4}, {1, 4}},
		exp: []interval{{1, 4}},
	}, {
		is:  []interval{{1, 4}, {1, 3}},
		exp: []interval{{1, 4}},
	}, {
		is:  []interval{{1, 3}, {1, 4}},
		exp: []interval{{1, 4}},
	}, {
		is:  []interval{{1, 2}, {4, 5}, {7, 8}},
		exp: []interval{{1, 2}, {4, 5}, {7, 8}},
	}, {
		is:  []interval{{7, 8}, {4, 5}, {1, 2}},
		exp: []interval{{1, 2}, {4, 5}, {7, 8}},
	}} {
		t.Run(fmt.Sprint(tc.is), func(t *testing.T) {
			var is intervals
			for _, i := range tc.is {
				is.add(i)
			}
			if !slices.Equal(is.is, tc.exp) {
				t.Errorf("Expected: %v Got: %v", tc.exp, is.is)
			}
		})
	}
}
