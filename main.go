package main

import (
	"fmt"

	"github.com/cherrysrc/GoUI/event"
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

	basePanel := makeUI(renderer)

	mObserver := event.NewMouseObserver()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.GetType() {
			case sdl.QUIT:
				running = false
				break
			case sdl.KEYDOWN:
				keydownEvent := event.(*sdl.KeyboardEvent)
				if keydownEvent.Keysym.Sym == sdl.K_BACKSPACE {
					if widget.SelectedTextArea != nil {
						widget.SelectedTextArea.PopText()
						err = widget.SelectedTextArea.RedrawText()
						if err != nil {
							panic(err)
						}
					}
				}
				break
			case sdl.TEXTINPUT:
				if widget.SelectedTextArea != nil {
					widget.SelectedTextArea.AppendText(event.(*sdl.TextInputEvent).GetText())
					err = widget.SelectedTextArea.RedrawText()
					if err != nil {
						panic(err)
					}
				}
				break
			case sdl.TEXTEDITING:
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

	editRect := sdl.Rect{10, 70, 780/2 - 20, 200}
	editTexture, err := widget.GenerateTexture(renderer, textureSpecs, int(editRect.W), int(editRect.H), 1)
	if err != nil {
		panic(err)
	}
	editWidget, err := widget.CreateTextArea(renderer, editRect, editTexture, "Tas", sdl.Color{0, 0, 0, 255}, font)
	if err != nil {
		panic(err)
	}

	basePanel.AddChild(leftPanel)
	basePanel.AddChild(titleLable)
	leftPanel.AddChild(button)
	leftPanel.AddChild(editWidget)

	return basePanel
}
