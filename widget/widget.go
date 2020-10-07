package widget

import "github.com/veandco/go-sdl2/sdl"

//TODO proper docs

//IWidget interface
type IWidget interface {
	Contains(int32, int32) bool

	GetChildren() []IWidget
	AddChild(IWidget)

	ClickEvent(int32, int32) bool

	IsClickable() bool
	OnClick()

	Draw(*sdl.Renderer)
}
