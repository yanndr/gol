package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	gridWidth  = 300
	gridHeight = 200
	width      = 1280
	height     = 1084
)

func main() {
	var (
		space, copy    [gridWidth][gridHeight]bool
		bgColor        = sdl.Color{R: 0, G: 0, B: 0, A: 255}
		gridEdgesColor = sdl.Color{R: 0, G: 255, B: 255, A: 255}
		gridColor      = sdl.Color{R: 0, G: 255, B: 0, A: 255}
		cellColorAlone = sdl.Color{R: 0, G: 255, B: 255, A: 255}
		cellColorDying = sdl.Color{R: 255, G: 0, B: 0, A: 255}
		cellColor      = sdl.Color{R: 255, G: 255, B: 0, A: 255}
	)

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, r, err := sdl.CreateWindowAndRenderer(width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	initGrindRandom(&space, 0.1)

	copy = space
	var cellSize int
	if width/gridWidth < height/gridHeight {
		cellSize = width / gridWidth
	} else {
		cellSize = height / gridHeight
	}

	//scale := float32(1.0)
	start := false
	vpx, vpy := 0, 0
	iteration := 0
	go func() {
		for {

			wstart := (width - cellSize*gridWidth) / 2
			hstart := (height - cellSize*gridHeight) / 2
			wstart += vpx
			hstart += vpy
			wend := wstart + cellSize*gridWidth
			hend := hstart + cellSize*gridHeight

			gridpxW := cellSize * gridWidth
			gridpxH := cellSize * gridHeight

			copy = space
			r.SetDrawColor(bgColor.R, bgColor.G, bgColor.B, bgColor.A)
			r.Clear()

			drawGridEdges(r, int32(wstart), int32(hstart), int32(wend), int32(hend), gridEdgesColor)

			for i := 0; i < len(space)-1; i++ {
				r.SetDrawColor(gridColor.R, gridColor.G, gridColor.B, gridColor.A)
				if cellSize > 15 {
					r.DrawLine(int32(wstart+(i+1)*cellSize), int32(hstart), int32(wstart+(i+1)*cellSize), int32(hstart+gridpxH))
				}
				for j := 0; j < len(space[0])-1; j++ {
					r.SetDrawColor(gridColor.R, gridColor.G, gridColor.B, gridColor.A)
					if cellSize > 15 {
						r.DrawLine(int32(wstart), int32(hstart+(j+1)*cellSize), int32(wstart+gridpxW), int32(hstart+(j+1)*cellSize))
					}
					an := numberOfAliveNeigbour(&copy, i, j)
					if copy[i][j] {
						// cell := sdl.Rect{X: int32((i + 1) * size), Y: int32((j + 1) * size), H: int32(size), W: int32(size)}
						// r.SetDrawColor(0, 255, 0, 128)
						// r.DrawRect(&cell)
						if an == 1 {
							gfx.BoxColor(r, int32(wstart+(i+1)*cellSize), int32(hstart+(j+1)*cellSize), int32(wstart+(i+1)*cellSize+cellSize), int32(hstart+(j+1)*cellSize+cellSize), cellColorAlone)
						} else if an == 2 || an == 3 {
							gfx.BoxColor(r, int32(wstart+(i+1)*cellSize), int32(hstart+(j+1)*cellSize), int32(wstart+(i+1)*cellSize+cellSize), int32(hstart+(j+1)*cellSize+cellSize), cellColor)
						} else {
							gfx.BoxColor(r, int32(wstart+(i+1)*cellSize), int32(hstart+(j+1)*cellSize), int32(wstart+(i+1)*cellSize+cellSize), int32(hstart+(j+1)*cellSize+cellSize), cellColorDying)
						}
						if start {
							space[i][j] = an > 1 && an < 4
						}
					} else {
						if start {
							space[i][j] = an == 3
						}
					}
				}
			}
			r.Present()
			if start {
				iteration++
			}
			time.Sleep(time.Second / 8)
		}
	}()

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
				cellSize = int(int32(cellSize) + mwe.Y)
				println(cellSize)
			case *sdl.KeyboardEvent:
				ke := event.(*sdl.KeyboardEvent)
				switch ke.Keysym.Scancode {
				case 79:
					vpx = vpx - 10
				case 80:
					vpx = vpx + 10
				case 81:
					vpy = vpy - 10
				case 82:
					vpy = vpy + 10
				case 44:
					if ke.State == 1 {
						start = !start
					}
				default:
					fmt.Println("unknow key:", ke.Keysym.Scancode)
				}

			}
		}
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

func drawGridEdges(r *sdl.Renderer, x, y, w, h int32, c sdl.Color) {
	r.SetDrawColor(c.R, c.G, c.B, c.A)
	r.DrawLine(x, y, w, y)
	r.DrawLine(x, y, x, h)
	r.DrawLine(w, y, w, h)
	r.DrawLine(x, h, w, h)
}
