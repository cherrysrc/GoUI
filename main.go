package main

import (
	"fmt"

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

	font, err := ttf.OpenFont("res/JetBrainsMono-Regular.ttf", 128)
	if err != nil {
		panic(err)
	}

	baseTex, err := img.LoadTexture(renderer, "res/metalPanel.png")
	if err != nil {
		panic(err)
	}

	basePanel := widget.CreatePanel(sdl.Rect{10, 10, 780, 580}, baseTex)
	leftPanel := widget.CreatePanel(sdl.Rect{0, 50, 780 / 2, 550 - 20}, baseTex)
	rightPanel := widget.CreatePanel(sdl.Rect{780 / 2, 50, 780 / 2, 550 - 20}, baseTex)

	var button01 *widget.Button
	if button01, err = widget.CreateButton(renderer, sdl.Rect{10, 10, 780/2 - 20, 50}, baseTex, "Button 01", sdl.Color{0, 0, 0, 255}, font, func() {
		fmt.Println("PRESSED")
	}); err != nil {
		panic(err)
	}

	var titleLabel *widget.Label
	if titleLabel, err = widget.CreateLabel(renderer, sdl.Rect{400 - 50, 0, 100, 50}, "Title", sdl.Color{0, 0, 0, 255}, font); err != nil {
		panic(err)
	}

	basePanel.AddChild(leftPanel)
	basePanel.AddChild(rightPanel)
	basePanel.AddChild(titleLabel)
	leftPanel.AddChild(button01)

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
