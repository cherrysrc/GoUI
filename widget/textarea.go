package widget

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

//TODO proper docs

//SelectedTextArea var
var SelectedTextArea *TextArea

//TextArea struct
type TextArea struct {
	rect        sdl.Rect
	baseTexture *sdl.Texture
	textTexture *sdl.Texture

	renderer *sdl.Renderer
	font     *ttf.Font
	color    sdl.Color

	text     string
	parent   IWidget
	children []IWidget
}

//CreateTextArea function
func CreateTextArea(renderer *sdl.Renderer, rect sdl.Rect, baseTexture *sdl.Texture, text string, color sdl.Color, font *ttf.Font) (*TextArea, error) {
	surface, err := font.RenderUTF8Solid(text, color)
	if err != nil {
		return nil, err
	}

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, err
	}
	surface.Free()

	return &TextArea{
		rect,
		baseTexture,
		texture,
		renderer,
		font,
		color,
		text,
		nil,
		make([]IWidget, 0),
	}, err
}

//GetRect function
func (ta *TextArea) GetRect() sdl.Rect {
	return ta.rect
}

//Contains function
func (ta *TextArea) Contains(x int32, y int32) bool {
	absRect := ta.GetAbsPosition()
	return !(x < absRect.X || x >= absRect.X+absRect.W || y < absRect.Y || y >= absRect.Y+absRect.H)
}

//GetAbsPosition function
func (ta *TextArea) GetAbsPosition() sdl.Rect {
	offset := ta.GetRect()

	//Sum up offsets
	for parent := ta.GetParent(); parent != nil; parent = parent.GetParent() {
		parentRect := parent.GetRect()
		offset.X += parentRect.X
		offset.Y += parentRect.Y
	}

	return offset
}

//SetParent function
func (ta *TextArea) SetParent(parent IWidget) {
	ta.parent = parent
}

//GetParent function
func (ta *TextArea) GetParent() IWidget {
	return ta.parent
}

//GetChildren function
func (ta *TextArea) GetChildren() []IWidget {
	return ta.children
}

//AddChild function
func (ta *TextArea) AddChild(child IWidget) {
	child.SetParent(ta)
	ta.children = append(ta.children, child)
}

//ClickEvent function. Checks children first, since they're drawn on top of the parent.
func (ta *TextArea) ClickEvent(mx int32, my int32) bool {
	for i := range ta.children {
		if ta.children[i].Contains(mx, my) {
			if ta.children[i].ClickEvent(mx, my) {
				return true
			}
		}
	}

	if ta.Contains(mx, my) && ta.IsClickable() {
		ta.OnClick()
		return true
	}
	return false
}

//IsClickable function
func (ta *TextArea) IsClickable() bool {
	return true
}

//OnClick function.
//Should be empty, if IsClickable returns false, since it will never be called anyways
func (ta *TextArea) OnClick() {
	SelectedTextArea = ta
}

//Draw function. Draws children after itself, to ensure they're draw on top.
func (ta *TextArea) Draw(renderer *sdl.Renderer) {
	//Child positions are relative to parent
	offset := ta.GetAbsPosition()

	renderer.Copy(ta.baseTexture, nil, &offset)
	renderer.Copy(ta.textTexture, nil, &offset)

	for i := range ta.children {
		ta.children[i].Draw(renderer)
	}
}

//AppendText function
func (ta *TextArea) AppendText(appendix string) {
	ta.text += appendix
	ta.RedrawText()
}

//PopText function
func (ta *TextArea) PopText() {
	if len(ta.text) > 0 {
		ta.text = ta.text[:len(ta.text)-1]
	}
}

//RedrawText function
func (ta *TextArea) RedrawText() error {
	var renderString string
	if len(ta.text) == 0 {
		renderString = " "
	} else {
		renderString = ta.text
	}

	surface, err := ta.font.RenderUTF8Solid(renderString, ta.color)
	if err != nil {
		return err
	}

	texture, err := ta.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return err
	}
	ta.textTexture = texture
	return nil
}
