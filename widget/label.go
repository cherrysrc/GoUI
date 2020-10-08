package widget

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

//TODO proper docs

//Label struct
type Label struct {
	rect     sdl.Rect
	texture  *sdl.Texture
	parent   IWidget
	children []IWidget
}

//CreateLabel function
func CreateLabel(renderer *sdl.Renderer, rect sdl.Rect, text string, color sdl.Color, font *ttf.Font) (*Label, error) {
	surface, err := font.RenderUTF8Solid(text, color)
	if err != nil {
		return nil, err
	}

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, err
	}
	surface.Free()

	return &Label{
		rect,
		texture,
		nil,
		make([]IWidget, 0),
	}, nil
}

//GetRect function
func (label *Label) GetRect() sdl.Rect {
	return label.rect
}

//Contains function
func (label *Label) Contains(x int32, y int32) bool {
	absRect := label.GetAbsPosition()
	return !(x < absRect.X || x >= absRect.X+absRect.W || y < absRect.Y || y >= absRect.Y+absRect.H)
}

//GetAbsPosition function
func (label *Label) GetAbsPosition() sdl.Rect {
	offset := label.GetRect()

	//Sum up offsets
	for parent := label.GetParent(); parent != nil; parent = parent.GetParent() {
		parentRect := parent.GetRect()
		offset.X += parentRect.X
		offset.Y += parentRect.Y
	}

	return offset
}

//SetParent function
func (label *Label) SetParent(parent IWidget) {
	label.parent = parent
}

//GetParent function
func (label *Label) GetParent() IWidget {
	return label.parent
}

//GetChildren function
func (label *Label) GetChildren() []IWidget {
	return label.children
}

//AddChild function
func (label *Label) AddChild(child IWidget) {
	child.SetParent(label)
	label.children = append(label.children, child)
}

//ClickEvent function. Checks children first, since they're drawn on top of the parent.
func (label *Label) ClickEvent(mx int32, my int32) bool {
	for i := range label.children {
		if label.children[i].Contains(mx, my) && label.children[i].IsClickable() {
			return label.children[i].ClickEvent(mx, my)
		}
	}

	if label.Contains(mx, my) && label.IsClickable() {
		label.OnClick()
	}
	return true
}

//IsClickable function
func (label *Label) IsClickable() bool {
	return false
}

//OnClick function.
//Should be empty, if IsClickable returns false, since it will never be called anyways
func (label *Label) OnClick() {
}

//Draw function. Draws children after itself, to ensure they're draw on top.
func (label *Label) Draw(renderer *sdl.Renderer) {
	//Child positions are relative to parent
	offset := label.GetRect()

	//Sum up offsets
	for parent := label.GetParent(); parent != nil; parent = parent.GetParent() {
		parentRect := parent.GetRect()
		offset.X += parentRect.X
		offset.Y += parentRect.Y
	}

	renderer.Copy(label.texture, nil, &offset)

	for i := range label.children {
		label.children[i].Draw(renderer)
	}
}
