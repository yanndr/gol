package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	gridWidth  = 120
	gridHeight = 100
	width      = 800
	height     = 600
)

var space, copy [gridWidth][gridHeight]bool

func main() {

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, r, err := sdl.CreateWindowAndRenderer(width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	for i := 0; i < len(space)-1; i++ {
		for j := 0; j < len(space[0])-1; j++ {
			space[i][j] = rand.Float32() > 0.6
		}
	}
	c := sdl.Color{R: 255, G: 255, B: 0, A: 255}

	copy = space
	var cellSize int
	if width/gridWidth < height/gridHeight {
		cellSize = width / gridWidth
	} else {
		cellSize = height / gridHeight
	}
	fmt.Println(cellSize)

	wstart := (width - cellSize*gridWidth) / 2
	hstart := (height - cellSize*gridHeight) / 2
	wend := wstart + cellSize*gridWidth
	hend := hstart + cellSize*gridHeight

	gridpxW := cellSize * gridWidth
	gridpxH := cellSize * gridHeight

	scale := float32(1.0)

	go func() {

		for {
			copy = space
			r.SetDrawColor(0, 0, 0, 255)
			r.Clear()

			r.SetScale(scale, scale)
			r.SetDrawColor(0, 255, 255, 255)
			r.DrawLine(int32(wstart), int32(hstart), int32(wend), int32(hstart))
			r.DrawLine(int32(wstart), int32(hstart), int32(wstart), int32(hend))
			r.DrawLine(int32(wend), int32(hstart), int32(wend), int32(hend))
			r.DrawLine(int32(wstart), int32(hend), int32(wend), int32(hend))
			for i := 0; i < len(space)-1; i++ {
				r.SetDrawColor(0, 255, 0, 255)
				r.DrawLine(int32(wstart+(i+1)*cellSize), int32(hstart), int32(wstart+(i+1)*cellSize), int32(gridpxH))
				for j := 0; j < len(space[0])-1; j++ {
					r.SetDrawColor(0, 255, 0, 255)

					r.DrawLine(int32(wstart), int32(hstart+(j+1)*cellSize), int32(wstart+gridpxW), int32(hstart+(j+1)*cellSize))
					an := numberOfAliveNeigbour(i, j)
					if copy[i][j] {
						// cell := sdl.Rect{X: int32((i + 1) * size), Y: int32((j + 1) * size), H: int32(size), W: int32(size)}
						// r.SetDrawColor(0, 255, 0, 128)
						// r.DrawRect(&cell)

						gfx.BoxColor(r, int32(wstart+(i+1)*cellSize), int32(hstart+(j+1)*cellSize), int32(wstart+(i+1)*cellSize+cellSize), int32(hstart+(j+1)*cellSize+cellSize), c)
						space[i][j] = an > 1 && an < 4
					} else {
						space[i][j] = an == 3
					}
				}
			}
			r.Present()
			time.Sleep(time.Second / 8)
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
			case *sdl.MouseWheelEvent:
				mwe := event.(*sdl.MouseWheelEvent)
				fmt.Println(mwe.Y)
				scale = scale + float32(mwe.Y)/10
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
