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
		grid  [gridWidth][gridHeight]bool
		start = make(chan bool)
		quit  = make(chan bool)
	)
	runtime.GOMAXPROCS(runtime.NumCPU())

	update := gol(&grid, start, quit)

	err := run(&grid, start, quit, update)
	if err != nil {
		log.Panicf("err running the SDL rendering: %v,", err)
	}

}
