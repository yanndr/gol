package main

import (
	"log"
	"runtime"
)

const (
	gridWidth  = 500
	gridHeight = 500
)

type pos struct {
	x, y int
}

func main() {

	var (
		grid  [][]bool
		start = make(chan bool)
		quit  = make(chan bool)
	)
	runtime.GOMAXPROCS(runtime.NumCPU())

	grid = make([][]bool, gridWidth)
	for i := 0; i < len(grid); i++ {
		grid[i] = make([]bool, gridHeight)
	}

	update := gol(grid, start, quit)

	err := run(grid, start, quit, update)
	if err != nil {
		log.Panicf("err running the SDL rendering: %v,", err)
	}

}
