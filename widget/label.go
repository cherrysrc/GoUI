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
		make([]IWidget, 0),
	}, nil
}

//Contains function
func (label *Label) Contains(x int32, y int32) bool {
	return !(x < label.rect.X || x >= label.rect.X+label.rect.W || y < label.rect.Y || y >= label.rect.Y+label.rect.H)
}

//GetChildren function
func (label *Label) GetChildren() []IWidget {
	return label.children
}

//AddChild function
func (label *Label) AddChild(child IWidget) {
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
	renderer.Copy(label.texture, nil, &label.rect)

	for i := range label.children {
		label.children[i].Draw(renderer)
	}
}
