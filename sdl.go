package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	width       = 1280
	height      = 1084
	minGridCell = 8
)

func run(grid *[gridWidth][gridHeight]bool, process, quitChan chan<- bool, update <-chan bool) error {

	var (
		bgColor        = sdl.Color{R: 0, G: 0, B: 0, A: 255}
		gridEdgesColor = sdl.Color{R: 0, G: 255, B: 255, A: 255}
		gridColor      = sdl.Color{R: 0, G: 255, B: 0, A: 255}
		cellColorAlone = sdl.Color{R: 0, G: 255, B: 255, A: 255}
		cellColorDying = sdl.Color{R: 255, G: 0, B: 0, A: 255}
		cellColor      = sdl.Color{R: 255, G: 255, B: 0, A: 255}
		cellColorNext  = sdl.Color{R: 255, G: 255, B: 0, A: 0}

		started   = false
		b1Click   = false
		actionAdd = false
		running   = true
		mouseMove = false
		iteration = 0
		alpha     = uint8(255)
		speed     = time.Microsecond
	)

	go func() {
		for {
			<-update
			iteration++
		}
	}()

	startAlphaChan := make(chan bool)
	quitAlphaChan := make(chan bool)
	go func() {
		for {
			select {
			case <-quitAlphaChan:
				return
			case run := <-startAlphaChan:
				for started && run {
					processGrid(&alpha, process, &speed)
				}
			}
		}

	}()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		log.Panicf("could not initialize TTF: %v", err)
	}
	defer ttf.Quit()

	window, r, err := sdl.CreateWindowAndRenderer(width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	var cellSize int
	if width/gridWidth < height/gridHeight {
		cellSize = width / gridWidth
	} else {
		cellSize = height / gridHeight
	}

	vpx, vpy := 0, 0

	var wstart, hstart, wend, hend, gridpxW, gridpxH int
	go func() {
		for running {

			wstart = (width - cellSize*gridWidth) / 2
			hstart = (height - cellSize*gridHeight) / 2
			wstart += vpx
			hstart += vpy
			wend = wstart + cellSize*gridWidth
			hend = hstart + cellSize*gridHeight

			gridpxW = cellSize * gridWidth
			gridpxH = cellSize * gridHeight

			r.SetDrawColor(bgColor.R, bgColor.G, bgColor.B, bgColor.A)
			r.Clear()

			drawGridEdges(r, int32(wstart), int32(hstart), int32(wend), int32(hend), gridEdgesColor)
			wg := sync.WaitGroup{}
			for i := 0; i < len(grid)-1; i++ {
				if cellSize > minGridCell {
					drawGridLine(r, int32(wstart+(i+1)*cellSize), int32(hstart), int32(wstart+(i+1)*cellSize), int32(hstart+gridpxH), gridColor)
				}
				wg.Add(1)
				go func() {
					defer wg.Done()
					for j := 0; j < len(grid[0])-1; j++ {
						if cellSize > minGridCell {
							drawGridLine(r, int32(wstart), int32(hstart+(j+1)*cellSize), int32(wstart+gridpxW), int32(hstart+(j+1)*cellSize), gridColor)
						}
						an := numberOfAliveNeigbour(grid, i, j)
						if grid[i][j] {
							drawCell(r, an, int32(wstart+(i)*cellSize), int32(hstart+(j)*cellSize), int32(wstart+(i)*cellSize+cellSize), int32(hstart+(j)*cellSize+cellSize), cellColorAlone, cellColor, cellColorDying, alpha)
						} else {
							if an == 3 {
								cellColorNext.A = 255 - alpha
								gfx.BoxColor(r, int32(wstart+(i)*cellSize), int32(hstart+(j)*cellSize), int32(wstart+(i)*cellSize+cellSize), int32(hstart+(j)*cellSize+cellSize), cellColorNext)
							}
						}
					}
				}()
				wg.Wait()
			}
			rect := sdl.Rect{X: 0, Y: 0, W: 100, H: 20}
			print(r, fmt.Sprintf("i: %v", iteration), cellColor, rect)

			r.Present()
		}
	}()

	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				quitAlphaChan <- true
				quitChan <- true

				break
			case *sdl.MouseWheelEvent:
				mwe := event.(*sdl.MouseWheelEvent)

				cellSize = int(int32(cellSize) + mwe.Y)
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
				case 40:
					if ke.State == 1 {
						go processGrid(&alpha, process, &speed)
					}
				case 44:
					if ke.State == 1 {
						started = !started
						startAlphaChan <- started
					}
				case 87:
					if ke.State == 0 {
						continue
					}
					if speed < time.Microsecond {
						continue
					}
					speed = speed / 2

				case 86:
					if ke.State == 0 {
						continue
					}
					if speed > time.Millisecond {
						continue
					}
					speed = speed * 2
				default:
					fmt.Println("unknow key:", ke.Keysym.Scancode)
				}
			case *sdl.MouseButtonEvent:
				me := event.(*sdl.MouseButtonEvent)
				b1Click = me.State == 1
				col := (int(me.X) - wstart) / cellSize
				row := (int(me.Y) - hstart) / cellSize
				if me.State == 1 {
					if !mouseMove {
						grid[col][row] = !grid[col][row]
					}
					mouseMove = false
				} else {

					if col >= 0 && row >= 0 && col < len(grid) && row < len(grid[0]) {
						actionAdd = !grid[col][row]
					}
				}

			case *sdl.MouseMotionEvent:
				me := event.(*sdl.MouseMotionEvent)
				col := (int(me.X) - wstart) / cellSize
				row := (int(me.Y) - hstart) / cellSize
				if b1Click && col >= 0 && row >= 0 && col < len(grid) && row < len(grid[0]) {
					mouseMove = true
					grid[col][row] = actionAdd
				}
			}
		}
	}
	return nil
}

func processGrid(alpha *uint8, process chan<- bool, duration *time.Duration) {
	*alpha = 255
	for *alpha != 0 {
		*alpha--
		//time.Sleep(*duration)
	}
	process <- true

}

func drawGridEdges(r *sdl.Renderer, x, y, w, h int32, c sdl.Color) {
	r.SetDrawColor(c.R, c.G, c.B, c.A)
	r.DrawLine(x, y, w, y)
	r.DrawLine(x, y, x, h)
	r.DrawLine(w, y, w, h)
	r.DrawLine(x, h, w, h)
}

func drawGridLine(r *sdl.Renderer, x, y, w, h int32, c sdl.Color) {
	r.SetDrawColor(c.R, c.G, c.B, c.A)
	r.DrawLine(x, y, w, h)
}

func drawCell(r *sdl.Renderer, aliveNeigbour int, x, y, w, h int32, c1, c2, c3 sdl.Color, alpha uint8) {
	var color sdl.Color
	if aliveNeigbour < 2 {
		color = c2
		color.A = alpha

	} else if aliveNeigbour == 2 || aliveNeigbour == 3 {
		color = c2
	} else {
		color = c2
		color.A = alpha
	}

	gfx.BoxColor(r, x, y, w, h, color)
}

func print(r *sdl.Renderer, text string, c sdl.Color, rect sdl.Rect) error {
	f, err := ttf.OpenFont("res/Roboto-Regular.ttf", 20)
	if err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	defer f.Close()

	s, err := f.RenderUTF8Solid(text, c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer t.Destroy()

	if err := r.Copy(t, nil, &rect); err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	return nil
}
