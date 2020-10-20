package widget

import (
	"github.com/cherrysrc/GoUI/util"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

//TODO proper docs
//TODO proper encapsulation

//MultiLineEdit struct
type MultiLineEdit struct {
	rect        sdl.Rect
	baseTexture *sdl.Texture

	textLineTextures []*sdl.Texture
	lineLengths      []int

	renderer  *sdl.Renderer
	font      *ttf.Font
	textColor sdl.Color

	text string

	charWidth     int
	charHeight    int
	lineCharLimit int
	lineCount     int

	caretPosition int
	caretRect     sdl.Rect
	caretTexture  *sdl.Texture

	parent   IWidget
	children []IWidget
}

//CreateMultiLineEdit function
func CreateMultiLineEdit(renderer *sdl.Renderer, baseTex *sdl.Texture, rect sdl.Rect, charWidth int, charHeight int, textColor sdl.Color, font *ttf.Font) *MultiLineEdit {
	charLimitX := int(rect.W) / charWidth
	charLimitY := int(rect.H) / charHeight
	return &MultiLineEdit{
		rect,
		baseTex,
		make([]*sdl.Texture, charLimitY),
		make([]int, charLimitY),
		renderer,
		font,
		textColor,
		"",
		charWidth,
		charHeight,
		charLimitX,
		charLimitY,
		0,
		sdl.Rect{},
		nil,
		nil,
		make([]IWidget, 0),
	}
}

//GetRect function
func (me *MultiLineEdit) GetRect() sdl.Rect {
	return me.rect
}

//Contains function
func (me *MultiLineEdit) Contains(x int32, y int32) bool {
	absRect := me.GetAbsPosition()
	return !(x < absRect.X || x >= absRect.X+absRect.W || y < absRect.Y || y >= absRect.Y+absRect.H)
}

//GetAbsPosition function
func (me *MultiLineEdit) GetAbsPosition() sdl.Rect {
	offset := me.GetRect()

	//Sum up offsets
	for parent := me.GetParent(); parent != nil; parent = parent.GetParent() {
		parentRect := parent.GetRect()
		offset.X += parentRect.X
		offset.Y += parentRect.Y
	}

	return offset
}

//SetParent function
func (me *MultiLineEdit) SetParent(parent IWidget) {
	me.parent = parent
	me.calculateCaretRect()
	me.generateCaretTexture()
}

//GetParent function
func (me *MultiLineEdit) GetParent() IWidget {
	return me.parent
}

//GetChildren function
func (me *MultiLineEdit) GetChildren() []IWidget {
	return me.children
}

//AddChild function
func (me *MultiLineEdit) AddChild(child IWidget) {
	child.SetParent(me)
	me.children = append(me.children, child)
}

//ClickEvent function. Checks children first, since they're drawn on top of the parent.
func (me *MultiLineEdit) ClickEvent(mx int32, my int32) bool {
	for i := range me.children {
		if me.children[i].Contains(mx, my) {
			if me.children[i].ClickEvent(mx, my) {
				return true
			}
		}
	}

	if me.Contains(mx, my) && me.IsClickable() {
		me.OnClick()
		me.calculateCaretPosition(mx, my)
		return true
	}
	return false
}

//IsClickable function
func (me *MultiLineEdit) IsClickable() bool {
	return true
}

//OnClick function.
//Should be empty, if IsClickable returns false, since it will never be called anyways
func (me *MultiLineEdit) OnClick() {
	SelectedTextReceiver = me
}

//Draw function. Draws children after itself, to ensure they're draw on top.
func (me *MultiLineEdit) Draw(renderer *sdl.Renderer) {
	//Child positions are relative to parent
	offset := me.GetAbsPosition()

	if me.baseTexture != nil {
		renderer.Copy(me.baseTexture, nil, &offset)
	}

	offset.H = int32(me.charHeight)

	for lineTexIdx := range me.textLineTextures {
		offset.W = int32(me.lineLengths[lineTexIdx] * me.charWidth)
		renderer.Copy(me.textLineTextures[lineTexIdx], nil, &offset)
		offset.Y += int32(me.charHeight)
	}

	if SelectedTextReceiver == me {
		renderer.Copy(me.caretTexture, nil, &me.caretRect)
	}

	for i := range me.children {
		me.children[i].Draw(renderer)
	}
}

//---
//MultiLineEdit only functions

//AppendText function
func (me *MultiLineEdit) AppendText(appendix string) {
	if len(me.text) >= me.lineCharLimit*me.lineCount {
		return
	}
	me.text = me.text[:me.caretPosition] + appendix + me.text[me.caretPosition:]
	me.caretPosition++
	me.calculateCaretRect()
}

//PopText function
func (me *MultiLineEdit) PopText() {
	if me.caretPosition > 0 {
		me.text = me.text[:me.caretPosition-1] + me.text[me.caretPosition:]
		me.caretPosition--
		me.calculateCaretRect()
	}
}

//EnterPressed function
//TODO less hacky solution
func (me *MultiLineEdit) EnterPressed() {
	charsLeft := me.lineCharLimit - (me.caretPosition % me.lineCharLimit)
	for i := 0; i < charsLeft; i++ {
		me.AppendText(" ")
	}
}

func (me *MultiLineEdit) generateCaretTexture() error {
	texture, err := me.renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STATIC, int32(me.caretRect.W), int32(me.caretRect.H))
	if err != nil {
		return err
	}

	pixels := make([]uint8, me.caretRect.W*me.caretRect.H*4)
	pixelIdx := 0

	for y := 0; y < int(me.caretRect.H); y++ {
		for x := 0; x < int(me.caretRect.W); x++ {
			pixels[pixelIdx+0] = 128
			pixels[pixelIdx+1] = 0
			pixels[pixelIdx+2] = 0
			pixels[pixelIdx+3] = 0
			pixelIdx += 4
		}
	}

	err = texture.Update(nil, pixels, int(me.caretRect.W)*4)
	if err != nil {
		return err
	}

	me.caretTexture = texture
	return nil
}

func (me *MultiLineEdit) calculateCaretPosition(mx, my int32) {
	absPos := me.GetAbsPosition()

	charPosX := int(mx-absPos.X) / me.charWidth
	charPosY := int(my-absPos.Y) / me.charHeight

	charIdx := charPosX + charPosY*me.lineCharLimit
	if charIdx > len(me.text)-1 {
		charIdx = len(me.text)
	}
	me.caretPosition = charIdx

	me.calculateCaretRect()
}

func (me *MultiLineEdit) calculateCaretRect() {
	absPos := me.GetAbsPosition()

	charPosX := me.caretPosition % me.lineCharLimit
	charPosY := me.caretPosition / me.lineCharLimit

	me.caretRect = sdl.Rect{
		X: absPos.X + int32(me.charWidth*charPosX),
		Y: absPos.Y + int32(me.charHeight*charPosY),
		W: int32(2),
		H: int32(me.charHeight),
	}
}

//RerenderText function
func (me *MultiLineEdit) RerenderText() error {
	lines := util.Chunks(me.text, me.lineCharLimit)

	for lineIdx := range me.textLineTextures {
		if lineIdx >= len(lines) {
			me.textLineTextures[lineIdx] = nil
			continue
		}

		var renderString string
		if len(lines[lineIdx]) == 0 {
			renderString = " "
		} else {
			renderString = lines[lineIdx]
		}

		surface, err := me.font.RenderUTF8Blended(renderString, me.textColor)
		if err != nil {
			return err
		}

		texture, err := me.renderer.CreateTextureFromSurface(surface)
		if err != nil {
			return err
		}
		surface.Free()

		me.textLineTextures[lineIdx] = texture
		me.lineLengths[lineIdx] = len(lines[lineIdx])
	}

	return nil
}
