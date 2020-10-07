package main

import (
	"github.com/cherrysrc/GoUI/event"

	"github.com/cherrysrc/GoUI/widget"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
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

	baseTex, err := img.LoadTexture(renderer, "res/metalPanel.png")
	if err != nil {
		panic(err)
	}

	panel := widget.CreatePanel(sdl.Rect{10, 10, 200, 200}, baseTex)
	cPanel := widget.CreatePanel(sdl.Rect{20, 20, 100, 100}, baseTex)
	panel.AddChild(cPanel)

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
				panel.ClickEvent(mObserver.MousePosition())
			}
		}

		renderer.SetDrawColor(50, 50, 50, 255)
		renderer.Clear()

		panel.Draw(renderer)

		renderer.Present()
	}
}
