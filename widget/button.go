package widget

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

//TODO proper docs

//Button struct
type Button struct {
	rect        sdl.Rect
	baseTexture *sdl.Texture
	textTexture *sdl.Texture
	parent      IWidget
	children    []IWidget
	onClick     func()
}

//CreateButton function
func CreateButton(renderer *sdl.Renderer, rect sdl.Rect, baseTexture *sdl.Texture, text string, color sdl.Color, font *ttf.Font, onClick func()) (*Button, error) {
	surface, err := font.RenderUTF8Solid(text, color)
	if err != nil {
		return nil, err
	}

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, err
	}
	surface.Free()

	return &Button{
		rect,
		baseTexture,
		texture,
		nil,
		make([]IWidget, 0),
		onClick,
	}, nil
}

//GetRect function
func (button *Button) GetRect() sdl.Rect {
	return button.rect
}

//Contains function
func (button *Button) Contains(x int32, y int32) bool {
	return !(x < button.rect.X || x >= button.rect.X+button.rect.W || y < button.rect.Y || y >= button.rect.Y+button.rect.H)
}

//SetParent function
func (button *Button) SetParent(parent IWidget) {
	button.parent = parent
}

//GetParent function
func (button *Button) GetParent() IWidget {
	return button.parent
}

//GetChildren function
func (button *Button) GetChildren() []IWidget {
	return button.children
}

//AddChild function
func (button *Button) AddChild(child IWidget) {
	child.SetParent(button)
	button.children = append(button.children, child)
}

//ClickEvent function. Checks children first, since they're drawn on top of the parent.
func (button *Button) ClickEvent(mx int32, my int32) bool {
	for i := range button.children {
		if button.children[i].Contains(mx, my) && button.children[i].IsClickable() {
			return button.children[i].ClickEvent(mx, my)
		}
	}

	if button.Contains(mx, my) && button.IsClickable() {
		button.OnClick()
	}
	return true
}

//IsClickable function
func (button *Button) IsClickable() bool {
	return true
}

//OnClick function.
//Should be empty, if IsClickable returns false, since it will never be called anyways
func (button *Button) OnClick() {
	button.onClick()
}

//Draw function. Draws children after itself, to ensure they're draw on top.
func (button *Button) Draw(renderer *sdl.Renderer) {
	//Child positions are relative to parent
	offset := button.GetRect()

	//Sum up offsets
	for parent := button.GetParent(); parent != nil; parent = parent.GetParent() {
		parentRect := parent.GetRect()
		offset.X += parentRect.X
		offset.Y += parentRect.Y
	}

	renderer.Copy(button.baseTexture, nil, &offset)
	renderer.Copy(button.textTexture, nil, &offset)

	for i := range button.children {
		button.children[i].Draw(renderer)
	}
}
