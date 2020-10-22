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
				if widget.SelectedTextReceiver != nil {
					widget.SelectedTextReceiver.PopText()
					err := widget.SelectedTextReceiver.RerenderText()
					if err != nil {
						panic(err)
					}
				}
			} else if keydownEvent.Keysym.Sym == sdl.K_RETURN {
				if widget.SelectedTextReceiver != nil {
					widget.SelectedTextReceiver.EnterPressed()
				}
			} else if keydownEvent.Keysym.Sym == sdl.K_LEFT {
				widget.SelectedTextReceiver.MoveCaretLeft()
			} else if keydownEvent.Keysym.Sym == sdl.K_RIGHT {
				widget.SelectedTextReceiver.MoveCaretRight()
			} else if keydownEvent.Keysym.Sym == sdl.K_UP {
				widget.SelectedTextReceiver.MoveCaretUp()
			} else if keydownEvent.Keysym.Sym == sdl.K_DOWN {
				widget.SelectedTextReceiver.MoveCaretDown()
			}
			break
		case sdl.TEXTINPUT:
			if widget.SelectedTextReceiver != nil {
				widget.SelectedTextReceiver.AppendText(event.(*sdl.TextInputEvent).GetText())
				err := widget.SelectedTextReceiver.RerenderText()
				if err != nil {
					panic(err)
				}
			}
			break
		case sdl.TEXTEDITING:
			//TODO implement functionality for selecting text
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
