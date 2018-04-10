package main

import (
	"log"
	"runtime"

	"github.com/yanndr/gol/engine/gol"
)

const (
	gridWidth  = 500
	gridHeight = 500
)

type pos struct {
	x, y int
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	g := gol.New(gridWidth, gridHeight)

	err := run(g)
	if err != nil {
		log.Panicf("err running the SDL rendering: %v,", err)
	}

}
