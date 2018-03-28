package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

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
		color = c1
		color.A = alpha

	} else if aliveNeigbour == 2 || aliveNeigbour == 3 {
		color = c2
	} else {
		color = c3
		color.A = alpha
	}

	gfx.BoxColor(r, x, y, w, h, color)
}

func print(r *sdl.Renderer, text string, c sdl.Color, rect sdl.Rect) error {
	f, err := ttf.OpenFont("res/Roboto-Regular.ttf", 20)
	if err != nil {
		fmt.Errorf("could not load font: %v", err)
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
