package main

import (
	"math/rand"
	"sync"
)

func gol(grid [][]bool, process, quit <-chan bool) chan bool {

	initGrindRandom(grid, 0.15)
	update := make(chan bool)
	go func() {
		for {
			select {
			case <-quit:
				close(update)
				return
			case <-process:
				generateNextState(grid)
				update <- true
			}
		}
	}()

	return update
}

func generateNextState(grid [][]bool) {
	cp := make([][]bool, len(grid))
	for i := 0; i < len(grid); i++ {
		cp[i] = make([]bool, len(grid[i]))
		copy(cp[i], grid[i])
	}

	wg := sync.WaitGroup{}
	for i := 0; i < len(grid); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < len(grid[0]); j++ {

				an := numberOfAliveNeigbour(cp, i, j)
				if cp[i][j] {
					grid[i][j] = an > 1 && an < 4

				} else {
					grid[i][j] = an == 3
				}
			}
		}()
		wg.Wait()
	}
}

func initGrindRandom(grid [][]bool, probability float32) {
	for i := 0; i < len(grid)-1; i++ {
		for j := 0; j < len(grid[0])-1; j++ {
			grid[i][j] = rand.Float32() < probability
		}
	}
}

func numberOfAliveNeigbour(grid [][]bool, x, y int) int {
	num := 0
	w := len(grid)
	h := len(grid[0])
	for i := -1; i < 2; i++ {
		if x+i < 0 {
			continue
		}
		if x+i >= w {
			continue
		}
		for j := -1; j < 2; j++ {
			if y+j < 0 {
				continue
			}
			if i == 0 && j == 0 {
				continue
			}
			if y+j >= h {
				continue
			}
			if grid[x+i][y+j] {
				num++
			}
		}
	}
	return num
}
