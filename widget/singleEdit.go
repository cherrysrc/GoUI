package widget

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

//TODO proper docs
//TODO proper encapsulation

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

	caretPosition int
	caretRect     sdl.Rect

	parent   IWidget
	children []IWidget
}

//CreateSingleLineEdit function
func CreateSingleLineEdit(renderer *sdl.Renderer, rect sdl.Rect, baseTexture *sdl.Texture, charWidth int, textColor sdl.Color, font *ttf.Font) (*SingleLineEdit, error) {
	surface, err := font.RenderUTF8Solid(" ", textColor)
	if err != nil {
		return nil, err
	}

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, err
	}
	surface.Free()

	edit := &SingleLineEdit{
		rect,
		baseTexture,
		texture,
		renderer,
		font,
		textColor,
		"",
		charWidth,
		0,
		sdl.Rect{},
		nil,
		make([]IWidget, 0),
	}

	return edit, nil
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
	ta.calculateCaretRect()
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
		ta.calculateCaretPosition(mx)
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
	SelectedTextReceiver = ta
}

//Draw function. Draws children after itself, to ensure they're draw on top.
func (ta *SingleLineEdit) Draw(renderer *sdl.Renderer) {
	//Child positions are relative to parent
	offset := ta.GetAbsPosition()

	if ta.baseTexture != nil {
		renderer.Copy(ta.baseTexture, nil, &offset)
	}

	offset.W = int32(len(ta.text) * ta.charWidth)
	renderer.Copy(ta.textTexture, nil, &offset)

	if SelectedTextReceiver == ta {
		renderer.SetDrawColor(255, 255, 255, 64)
		renderer.FillRect(&ta.caretRect)
	}

	for i := range ta.children {
		ta.children[i].Draw(renderer)
	}
}

//---
//SingleLineEdit only functions

//AppendText function
func (ta *SingleLineEdit) AppendText(appendix string) {
	charLimit := int(ta.rect.W) / ta.charWidth
	if len(ta.text) >= charLimit {
		return
	}

	ta.text = ta.text[:ta.caretPosition] + appendix + ta.text[ta.caretPosition:]
	ta.caretPosition++
	ta.calculateCaretRect()
}

//PopText function
func (ta *SingleLineEdit) PopText() {
	if ta.caretPosition > 0 {
		ta.text = ta.text[:ta.caretPosition-1] + ta.text[ta.caretPosition:]
		ta.caretPosition--
		ta.calculateCaretRect()
	}
}

//EnterPressed function
func (ta *SingleLineEdit) EnterPressed() {
	return
}

//MoveCaretRight function
func (ta *SingleLineEdit) MoveCaretRight() {
	ta.caretPosition++
	ta.caretBounds()
	ta.calculateCaretRect()
}

//MoveCaretLeft function
func (ta *SingleLineEdit) MoveCaretLeft() {
	ta.caretPosition--
	ta.caretBounds()
	ta.calculateCaretRect()
}

//MoveCaretUp function
func (ta *SingleLineEdit) MoveCaretUp() {
}

//MoveCaretDown function
func (ta *SingleLineEdit) MoveCaretDown() {
}

func (ta *SingleLineEdit) caretBounds() {
	if ta.caretPosition < 0 {
		ta.caretPosition = 0
	} else if ta.caretPosition >= len(ta.text) {
		ta.caretPosition = len(ta.text)
	}
}

func (ta *SingleLineEdit) calculateCaretPosition(mx int32) {
	relX := mx - ta.GetAbsPosition().X
	idx := int(relX) / ta.charWidth
	if idx > len(ta.text)-1 {
		idx = len(ta.text)
	}
	ta.caretPosition = idx
	ta.calculateCaretRect()
}

func (ta *SingleLineEdit) calculateCaretRect() {
	absPos := ta.GetAbsPosition()

	x := absPos.X + int32(ta.caretPosition*ta.charWidth)
	y := absPos.Y + 1
	//w := int32(ta.charWidth)
	w := int32(ta.charWidth)
	h := absPos.H - 2

	ta.caretRect = sdl.Rect{X: x, Y: y, W: w, H: h}
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
	surface.Free()

	ta.textTexture = texture
	return nil
}
