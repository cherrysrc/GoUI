package widget

import "github.com/veandco/go-sdl2/sdl"

//TODO proper docs
//TODO encapsulation

//SelectedTextReceiver var
var SelectedTextReceiver ITextReceiver

//IWidget interface
type IWidget interface {
	GetRect() sdl.Rect
	Contains(int32, int32) bool
	GetAbsPosition() sdl.Rect

	SetParent(IWidget)
	GetParent() IWidget

	GetChildren() []IWidget
	AddChild(IWidget)

	ClickEvent(int32, int32) bool

	IsClickable() bool
	OnClick()

	Draw(*sdl.Renderer)
}

//ITextReceiver interface
type ITextReceiver interface {
	AppendText(string)
	PopText()
	EnterPressed()
	RerenderText() error
}
