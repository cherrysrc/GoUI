package ui

import (
	"github.com/cherrysrc/GoUI/event"
	"github.com/cherrysrc/GoUI/widget"
	"github.com/veandco/go-sdl2/sdl"
)

//TODO proper docs
//TODO nuklear style auto fitting of components

//UI struct
type UI struct {
	masterWidget  widget.IWidget
	mouseObserver *event.MouseStateObserver
}

//Create function
func Create(master widget.IWidget) *UI {
	return &UI{
		master,
		event.NewMouseObserver(),
	}
}

//PollEvent function
func (ui *UI) PollEvent() bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.GetType() {
		case sdl.QUIT:
			return false
		case sdl.KEYDOWN:
			keydownEvent := event.(*sdl.KeyboardEvent)
			if keydownEvent.Keysym.Sym == sdl.K_BACKSPACE {
				if widget.SelectedSingleLineEdit != nil {
					widget.SelectedSingleLineEdit.PopText()
					err := widget.SelectedSingleLineEdit.RerenderText()
					if err != nil {
						panic(err)
					}
				}
			}
			break
		case sdl.TEXTINPUT:
			if widget.SelectedSingleLineEdit != nil {
				widget.SelectedSingleLineEdit.AppendText(event.(*sdl.TextInputEvent).GetText())
				err := widget.SelectedSingleLineEdit.RerenderText()
				if err != nil {
					panic(err)
				}
			}
			break
		case sdl.TEXTEDITING:
			//TODO implement functionality
			break
		}
	}

	if event, eventHappened := ui.mouseObserver.Update(); eventHappened {
		if event == 1 {
			ui.masterWidget.ClickEvent(ui.mouseObserver.MousePosition())
		}
	}

	return true
}

//Draw function
func (ui *UI) Draw(renderer *sdl.Renderer) {
	ui.masterWidget.Draw(renderer)
}
