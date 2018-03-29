package main

import (
	"log"
	"runtime"
	"time"
)

const (
	gridWidth  = 100
	gridHeight = 100
)

type pos struct {
	x, y int
}

func main() {

	var (
		grid  [gridWidth][gridHeight]bool
		start = make(chan bool)
		quit  = make(chan bool)
		speed = time.Second / 8
	)
	runtime.GOMAXPROCS(runtime.NumCPU())

	update := gol(&grid, start, quit, &speed)

	err := run(&grid, start, quit, update, &speed)
	if err != nil {
		log.Panicf("err running the SDL rendering: %v,", err)
	}

}
