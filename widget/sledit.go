package widget

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

//TODO proper docs

//SelectedSingleLineEdit var
var SelectedSingleLineEdit *SingleLineEdit

//SingleLineEdit struct
type SingleLineEdit struct {
	rect        sdl.Rect
	baseTexture *sdl.Texture
	textTexture *sdl.Texture

	renderer *sdl.Renderer
	font     *ttf.Font
	color    sdl.Color

	text      string
	charWidth int

	parent   IWidget
	children []IWidget
}

//CreateSingleLineEdit function
func CreateSingleLineEdit(renderer *sdl.Renderer, rect sdl.Rect, baseTexture *sdl.Texture, text string, charWidth int, color sdl.Color, font *ttf.Font) (*SingleLineEdit, error) {
	surface, err := font.RenderUTF8Solid(text, color)
	if err != nil {
		return nil, err
	}

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, err
	}
	surface.Free()

	return &SingleLineEdit{
		rect,
		baseTexture,
		texture,
		renderer,
		font,
		color,
		text,
		charWidth,
		nil,
		make([]IWidget, 0),
	}, err
}

//GetRect function
func (ta *SingleLineEdit) GetRect() sdl.Rect {
	return ta.rect
}

//Contains function
func (ta *SingleLineEdit) Contains(x int32, y int32) bool {
	absRect := ta.GetAbsPosition()
	return !(x < absRect.X || x >= absRect.X+absRect.W || y < absRect.Y || y >= absRect.Y+absRect.H)
}

//GetAbsPosition function
func (ta *SingleLineEdit) GetAbsPosition() sdl.Rect {
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
func (ta *SingleLineEdit) SetParent(parent IWidget) {
	ta.parent = parent
}

//GetParent function
func (ta *SingleLineEdit) GetParent() IWidget {
	return ta.parent
}

//GetChildren function
func (ta *SingleLineEdit) GetChildren() []IWidget {
	return ta.children
}

//AddChild function
func (ta *SingleLineEdit) AddChild(child IWidget) {
	child.SetParent(ta)
	ta.children = append(ta.children, child)
}

//ClickEvent function. Checks children first, since they're drawn on top of the parent.
func (ta *SingleLineEdit) ClickEvent(mx int32, my int32) bool {
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
func (ta *SingleLineEdit) IsClickable() bool {
	return true
}

//OnClick function.
//Should be empty, if IsClickable returns false, since it will never be called anyways
func (ta *SingleLineEdit) OnClick() {
	SelectedSingleLineEdit = ta
}

//Draw function. Draws children after itself, to ensure they're draw on top.
func (ta *SingleLineEdit) Draw(renderer *sdl.Renderer) {
	//Child positions are relative to parent
	offset := ta.GetAbsPosition()
	renderer.Copy(ta.baseTexture, nil, &offset)

	offset.W = int32(len(ta.text) * ta.charWidth)
	renderer.Copy(ta.textTexture, nil, &offset)

	for i := range ta.children {
		ta.children[i].Draw(renderer)
	}
}

//AppendText function
func (ta *SingleLineEdit) AppendText(appendix string) {
	charLimit := int(ta.rect.W) / ta.charWidth
	if len(ta.text) >= charLimit {
		return
	}
	ta.text += appendix
	ta.RerenderText()
}

//PopText function
func (ta *SingleLineEdit) PopText() {
	if len(ta.text) > 0 {
		ta.text = ta.text[:len(ta.text)-1]
	}
}

//RerenderText function
func (ta *SingleLineEdit) RerenderText() error {
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
