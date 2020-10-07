package main

import (
	"github.com/cherrysrc/GoUI/event"

	"github.com/cherrysrc/GoUI/widget"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	window, err := sdl.CreateWindow("TestUI", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	err = ttf.Init()
	if err != nil {
		panic(err)
	}

	font, err := ttf.OpenFont("res/JetBrainsMono-Regular.ttf", 50)
	if err != nil {
		panic(err)
	}

	baseTex, err := img.LoadTexture(renderer, "res/metalPanel.png")
	if err != nil {
		panic(err)
	}

	basePanel := widget.CreatePanel(sdl.Rect{10, 10, 200, 200}, baseTex)
	bLabel, err := widget.CreateLabel(renderer, sdl.Rect{10, 10, 200, 100}, "Labelo", sdl.Color{255, 255, 255, 255}, font)
	if err != nil {
		panic(err)
	}
	basePanel.AddChild(bLabel)

	mObserver := event.NewMouseObserver()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
				break
			}
		}

		if event, eventHappened := mObserver.Update(); eventHappened {
			if event == 1 {
				basePanel.ClickEvent(mObserver.MousePosition())
			}
		}

		renderer.SetDrawColor(50, 50, 50, 255)
		renderer.Clear()

		basePanel.Draw(renderer)

		renderer.Present()
	}

	ttf.Quit()
}
