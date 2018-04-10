package main

import (
	"testing"
)

func TestGenerateNextState(t *testing.T) {

	tt := []struct {
		name string
		ig   [][]bool
		og   [][]bool
	}{
		{
			"basic",
			[][]bool{{false, false, false}, {false, true, false}, {false, false, false}},
			[][]bool{{false, false, false}, {false, false, false}, {false, false, false}},
		},
		{
			"clign",
			[][]bool{
				{false, true, false},
				{false, true, false},
				{false, true, false},
			},
			[][]bool{
				{false, false, false},
				{true, true, true},
				{false, false, false},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			generateNextState(tc.ig)

			for i := 0; i < len(tc.ig); i++ {
				for j := 0; j < len(tc.ig[i]); j++ {
					if tc.ig[i][j] != tc.og[i][j] {
						t.Fatalf("error in row %v col %v expected %v got %v", i, j, tc.og[i][j], tc.ig[i][j])
					}
				}
			}
		})
	}

}

func TestAliveNeighbour(t *testing.T) {

	tt := []struct {
		name string
		ig   [][]bool
		og   [][]int
	}{
		{
			"one in a center 3*3",
			[][]bool{{false, false, false}, {false, true, false}, {false, false, false}},
			[][]int{{1, 1, 1}, {1, 0, 1}, {1, 1, 1}},
		},
		{
			"one row 3*3",
			[][]bool{
				{false, true, false},
				{false, true, false},
				{false, true, false},
			},
			[][]int{
				{2, 1, 2},
				{3, 2, 3},
				{2, 1, 2},
			},
		},
		{
			"one 3 row 5*5",
			[][]bool{
				{false, false, false, false, false},
				{false, false, true, false, false},
				{false, false, true, false, false},
				{false, false, true, false, false},
				{false, false, false, false, false},
			},
			[][]int{
				{0, 1, 1, 1, 0},
				{0, 2, 1, 2, 0},
				{0, 3, 2, 3, 0},
				{0, 2, 1, 2, 0},
				{0, 1, 1, 1, 0},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			for i := 0; i < len(tc.ig); i++ {
				for j := 0; j < len(tc.ig[i]); j++ {
					an := numberOfAliveNeigbour(tc.ig, i, j)
					if an != tc.og[i][j] {
						t.Fatalf("error in row %v col %v expected %v got %v", i, j, tc.og[i][j], an)
					}
				}
			}
		})
	}

}
