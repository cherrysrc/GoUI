package widget

import (
	"github.com/veandco/go-sdl2/sdl"
)

//TODO proper docs

//Panel struct
type Panel struct {
	rect     sdl.Rect
	texture  *sdl.Texture
	parent   IWidget
	children []IWidget
}

//CreatePanel function
func CreatePanel(rect sdl.Rect, texture *sdl.Texture) *Panel {
	return &Panel{
		rect,
		texture,
		nil,
		make([]IWidget, 0),
	}
}

//---
//IWidget implementation

//GetRect function
func (panel *Panel) GetRect() sdl.Rect {
	return panel.rect
}

//Contains function
func (panel *Panel) Contains(x int32, y int32) bool {
	absRect := panel.GetAbsPosition()
	return !(x < absRect.X || x >= absRect.X+absRect.W || y < absRect.Y || y >= absRect.Y+absRect.H)
}

//GetAbsPosition function
func (panel *Panel) GetAbsPosition() sdl.Rect {
	offset := panel.GetRect()

	//Sum up offsets
	for parent := panel.GetParent(); parent != nil; parent = parent.GetParent() {
		parentRect := parent.GetRect()
		offset.X += parentRect.X
		offset.Y += parentRect.Y
	}

	return offset
}

//SetParent function
func (panel *Panel) SetParent(parent IWidget) {
	panel.parent = parent
}

//GetParent function
func (panel *Panel) GetParent() IWidget {
	return panel.parent
}

//GetChildren function
func (panel *Panel) GetChildren() []IWidget {
	return panel.children
}

//AddChild function
func (panel *Panel) AddChild(child IWidget) {
	child.SetParent(panel)
	panel.children = append(panel.children, child)
}

//ClickEvent function. Checks children first, since they're drawn on top of the parent.
func (panel *Panel) ClickEvent(mx int32, my int32) bool {
	for i := range panel.children {
		if panel.children[i].Contains(mx, my) {
			if panel.children[i].ClickEvent(mx, my) {
				return true
			}
		}
	}

	if panel.Contains(mx, my) && panel.IsClickable() {
		panel.OnClick()
		return true
	}
	return false
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
	//Child positions are relative to parent
	offset := panel.GetRect()

	//Sum up offsets
	for parent := panel.GetParent(); parent != nil; parent = parent.GetParent() {
		parentRect := parent.GetRect()
		offset.X += parentRect.X
		offset.Y += parentRect.Y
	}

	renderer.Copy(panel.texture, nil, &offset)

	for i := range panel.children {
		panel.children[i].Draw(renderer)
	}
}
