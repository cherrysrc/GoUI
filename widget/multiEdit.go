package widget

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

//TODO proper docs
//TODO proper encapsulation

//MultiLineEdit struct
type MultiLineEdit struct {
	rect        sdl.Rect
	baseTexture *sdl.Texture

	lineEdits []*SingleLineEdit
	lineCount int
	charWidth int
	textColor sdl.Color

	renderer *sdl.Renderer
	font     *ttf.Font

	parent   IWidget
	children []IWidget
}

//CreateMultiLineEdit function
func CreateMultiLineEdit(renderer *sdl.Renderer, baseTex *sdl.Texture, rect sdl.Rect, lineCount int, charWidth int, textColor sdl.Color, font *ttf.Font) *MultiLineEdit {
	return &MultiLineEdit{
		rect,
		baseTex,
		make([]*SingleLineEdit, lineCount),
		lineCount,
		charWidth,
		textColor,
		renderer,
		font,
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
	me.generateLines()
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
	for i := range me.lineEdits {
		if me.lineEdits[i].Contains(mx, my) {
			if me.lineEdits[i].ClickEvent(mx, my) {
				SelectedSingleLineEdit = me.lineEdits[i]
				return true
			}
		}
	}

	for i := range me.children {
		if me.children[i].Contains(mx, my) {
			if me.children[i].ClickEvent(mx, my) {
				return true
			}
		}
	}

	if me.Contains(mx, my) && me.IsClickable() {
		me.OnClick()
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
}

//Draw function. Draws children after itself, to ensure they're draw on top.
func (me *MultiLineEdit) Draw(renderer *sdl.Renderer) {
	//Child positions are relative to parent
	offset := me.GetAbsPosition()
	renderer.Copy(me.baseTexture, nil, &offset)

	for i := range me.lineEdits {
		me.lineEdits[i].Draw(renderer)
	}

	for i := range me.children {
		me.children[i].Draw(renderer)
	}
}

//---
//MultiLineEdit specific functions

func (me *MultiLineEdit) generateLines() error {
	absPos := me.GetAbsPosition()
	linePixelHeight := int(me.rect.H) / me.lineCount

	for i := range me.lineEdits {
		x := absPos.X                            //TODO use actual margin
		y := absPos.Y + int32(i*linePixelHeight) //TODO use actual margin
		w := absPos.W
		h := absPos.H / int32(me.lineCount)

		rect := sdl.Rect{X: x, Y: y, W: w, H: h}
		line, err := CreateSingleLineEdit(me.renderer, rect, nil, me.charWidth, me.textColor, me.font)
		if err != nil {
			return err
		}
		line.calculateCaretRect()
		line.generateCaretTexture()
		me.lineEdits[i] = line
	}

	return nil
}
