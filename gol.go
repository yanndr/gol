package main

import (
	"math/rand"
	"sync"
	"time"
)

func gol(grid *[gridWidth][gridHeight]bool, c, quit <-chan bool, speed *time.Duration) chan bool {

	initGrindRandom(grid, 0.15)
	update := make(chan bool)
	go func() {
		started := false
		for {
			select {
			case <-quit:
				return
			case started = <-c:
				go func() {
					for started {
						generateNextState(grid)
						update <- true
						time.Sleep(*speed)
					}
				}()
			}
		}
	}()

	return update
}

func generateNextState(grid *[gridWidth][gridHeight]bool) {
	copy := *grid
	wg := sync.WaitGroup{}
	for i := 0; i < len(grid)-1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < len(grid[0])-1; j++ {

				an := numberOfAliveNeigbour(&copy, i, j)
				if copy[i][j] {
					grid[i][j] = an > 1 && an < 4

				} else {
					grid[i][j] = an == 3
				}
			}
		}()
		wg.Wait()
	}
}

func initGrindRandom(grid *[gridWidth][gridHeight]bool, probability float32) {
	for i := 0; i < len(grid)-1; i++ {
		for j := 0; j < len(grid[0])-1; j++ {
			grid[i][j] = rand.Float32() < probability
		}
	}
}

func numberOfAliveNeigbour(grid *[gridWidth][gridHeight]bool, x, y int) int {
	num := 0
	for i := -1; i < 2; i++ {
		if x+i < 0 {
			continue
		}
		for j := -1; j < 2; j++ {
			if y+j < 0 {
				continue
			}
			if i == 0 && j == 0 {
				continue
			}
			if grid[x+i][y+j] {
				num++
			}
		}
	}
	return num
}
