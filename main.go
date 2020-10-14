package main

import (
	"fmt"

	"github.com/cherrysrc/GoUI/ui"

	"github.com/cherrysrc/GoUI/widget"

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
	sdl.StartTextInput()

	gui := ui.Create(makeUI(renderer))

	running := true
	for running {
		running = gui.PollEvent()

		renderer.SetDrawColor(50, 50, 50, 255)
		renderer.Clear()

		gui.Draw(renderer)

		renderer.Present()
	}

	sdl.StopTextInput()
	ttf.Quit()
	sdl.Quit()
}

func makeUI(renderer *sdl.Renderer) widget.IWidget {
	font, err := ttf.OpenFont("res/JetBrainsMono-Regular.ttf", 128)
	if err != nil {
		panic(err)
	}

	textureSpecs := widget.TextureSpecs{
		sdl.Color{225, 225, 255, 255},
		sdl.Color{90, 90, 90, 0},
	}

	//----
	//base panel
	basePanelRect := sdl.Rect{10, 10, 780, 580}
	bgeneratedTexture, err := widget.GenerateTexture(renderer, textureSpecs, int(basePanelRect.W), int(basePanelRect.H), 1)
	if err != nil {
		panic(err)
	}
	basePanel := widget.CreatePanel(basePanelRect, bgeneratedTexture)

	//----
	//Left Panel
	leftPanelRect := sdl.Rect{10, 50, 780 / 2, 520}
	lgeneratedTexture, err := widget.GenerateTexture(renderer, textureSpecs, int(leftPanelRect.W), int(leftPanelRect.H), 1)
	if err != nil {
		panic(err)
	}
	leftPanel := widget.CreatePanel(leftPanelRect, lgeneratedTexture)

	//----
	//Title label
	labelRect := sdl.Rect{basePanelRect.W/2 - 100, 0, 200, 50}
	titleLable, err := widget.CreateLabel(renderer, labelRect, "Title Here", sdl.Color{0, 0, 0, 255}, font)
	if err != nil {
		panic(err)
	}

	//----
	//First button in left sub panel
	buttonRect := sdl.Rect{10, 10, 780/2 - 20, 50}
	buttonTexture, err := widget.GenerateTexture(renderer, textureSpecs, int(buttonRect.W), int(buttonRect.H), 1)
	if err != nil {
		panic(err)
	}
	button, err := widget.CreateButton(renderer, buttonRect, buttonTexture, "Button Here", sdl.Color{0, 0, 0, 255}, font, func() {
		fmt.Println("Button Pressed")
	})
	if err != nil {
		panic(err)
	}

	editRect := sdl.Rect{10, 70, 780/2 - 20, 50}
	editTexture, err := widget.GenerateTexture(renderer, textureSpecs, int(editRect.W), int(editRect.H), 1)
	if err != nil {
		panic(err)
	}
	editWidget, err := widget.CreateSingleLineEdit(renderer, editRect, editTexture, 20, sdl.Color{0, 0, 0, 255}, font)
	if err != nil {
		panic(err)
	}
	editRect2 := sdl.Rect{10, 130, 780/2 - 20, 50}
	//nil to use no background texture
	editWidget2, err := widget.CreateSingleLineEdit(renderer, editRect2, editTexture, 20, sdl.Color{0, 0, 0, 255}, font)
	if err != nil {
		panic(err)
	}

	basePanel.AddChild(leftPanel)
	basePanel.AddChild(titleLable)
	leftPanel.AddChild(button)
	leftPanel.AddChild(editWidget)
	leftPanel.AddChild(editWidget2)

	//----
	//Right Panel
	rightPanelRect := sdl.Rect{780/2 + 5, 50, 780/2 - 15, 520}
	rTexture, err := widget.GenerateTexture(renderer, textureSpecs, int(rightPanelRect.W), int(rightPanelRect.H), 1)
	if err != nil {
		panic(err)
	}
	rightPanel := widget.CreatePanel(rightPanelRect, rTexture)

	//----
	//MultLineEdit
	meEditRect := sdl.Rect{10, 10, 780/2 - 35, 500}
	meTexture, err := widget.GenerateTexture(renderer, textureSpecs, int(rightPanelRect.W), int(rightPanelRect.H), 1)
	if err != nil {
		panic(err)
	}
	meEdit := widget.CreateMultiLineEdit(renderer, meTexture, meEditRect, 8, 20, sdl.Color{0, 0, 0, 255}, font)

	basePanel.AddChild(rightPanel)
	rightPanel.AddChild(meEdit)

	return basePanel
}
