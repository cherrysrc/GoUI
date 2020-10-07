package widget

import (
	"github.com/veandco/go-sdl2/sdl"
)

//TODO proper docs

//Panel struct
type Panel struct {
	rect     sdl.Rect
	texture  *sdl.Texture
	children []IWidget
}

//CreatePanel function
func CreatePanel(rect sdl.Rect, texture *sdl.Texture) *Panel {
	return &Panel{
		rect,
		texture,
		make([]IWidget, 0),
	}
}

//---
//IWidget implementation

//Contains function
func (panel *Panel) Contains(x int32, y int32) bool {
	return !(x < panel.rect.X || x >= panel.rect.X+panel.rect.W || y < panel.rect.Y || y >= panel.rect.Y+panel.rect.H)
}

//GetChildren function
func (panel *Panel) GetChildren() []IWidget {
	return panel.children
}

//AddChild function
func (panel *Panel) AddChild(child IWidget) {
	panel.children = append(panel.children, child)
}

//ClickEvent function. Checks children first, since they're drawn on top of the parent.
func (panel *Panel) ClickEvent(mx int32, my int32) bool {
	for i := range panel.children {
		if panel.children[i].Contains(mx, my) && panel.children[i].IsClickable() {
			return panel.children[i].ClickEvent(mx, my)
		}
	}

	if panel.Contains(mx, my) && panel.IsClickable() {
		panel.OnClick()
	}
	return true
}

//IsClickable function
func (panel *Panel) IsClickable() bool {
	return false
}

//OnClick function.
//Should be empty, if IsClickable returns false, since it will never be called anyways
func (panel *Panel) OnClick() {
}

//Draw function. Draws children after itself, to ensure they're draw on top.
func (panel *Panel) Draw(renderer *sdl.Renderer) {
	renderer.Copy(panel.texture, nil, &panel.rect)

	for i := range panel.children {
		panel.children[i].Draw(renderer)
	}
}
