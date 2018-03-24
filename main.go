package main

import (
	"time"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	gridSize = 40
	widht    = 800
	height   = 800
)

var space, copy [gridSize][gridSize]bool

func main() {

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, r, err := sdl.CreateWindowAndRenderer(widht, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	space[15][15] = true
	space[16][15] = true
	space[17][15] = true
	space[18][15] = true
	space[19][15] = true

	//for it := 0; it < 1; it++ {
	c := sdl.Color{R: 255, G: 255, B: 0, A: 128}

	copy = space

	size := widht / gridSize
	go func() {
		for {
			copy = space
			r.SetDrawColor(0, 0, 0, 255)
			r.Clear()
			r.SetDrawColor(0, 255, 0, 255)
			for i := 0; i < len(space)-1; i++ {
				r.SetDrawColor(0, 255, 0, 255)
				r.DrawLine(int32((i+1)*size), int32(0), int32((i+1)*size), int32(height))
				for j := 0; j < len(space[0])-1; j++ {
					r.SetDrawColor(0, 255, 0, 255)
					r.DrawLine(int32(0), int32((j+1)*size), int32(widht), int32((j+1)*size))
					an := numberOfAliveNeigbour(i, j)
					if copy[i][j] {
						// cell := sdl.Rect{X: int32((i + 1) * size), Y: int32((j + 1) * size), H: int32(size), W: int32(size)}
						// r.SetDrawColor(0, 255, 0, 128)
						// r.DrawRect(&cell)
						gfx.BoxColor(r, int32((i+1)*size), int32((j+1)*size), int32((i+1)*size+size), int32((j+1)*size+size), c)
						space[i][j] = an > 1 && an < 4
					} else {
						space[i][j] = an == 3
					}
				}
			}
			r.Present()
			time.Sleep(time.Second / 2)
		}
	}()

	// fmt.Println(space[15][15])
	// fmt.Println(space[15][16])
	// fmt.Println(space[15][17])
	// time.Sleep(time.Second)
	//}

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}

func numberOfAliveNeigbour(x, y int) int {
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
			if copy[x+i][y+j] {
				num++
			}
		}
	}

	return num
}
