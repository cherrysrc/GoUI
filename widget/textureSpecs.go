package widget

import "github.com/veandco/go-sdl2/sdl"

//Generate textures to fit
//Using standard loaded textures will result in heavy stretching, which is ugly
//Using this you can generate a texture exactly fitting your ui element

//TextureSpecs struct
type TextureSpecs struct {
	BorderColor sdl.Color
	BodyColor   sdl.Color
}

//GenerateTexture function
func GenerateTexture(renderer *sdl.Renderer, specs TextureSpecs, width, height, borderThickness int) (*sdl.Texture, error) {
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STATIC, int32(width), int32(height))
	if err != nil {
		return nil, err
	}

	pixels := make([]uint8, width*height*4)
	pixelIdx := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if y <= borderThickness || x <= borderThickness || x >= width-1-borderThickness || y >= height-1-borderThickness {
				pixels[pixelIdx+0] = specs.BorderColor.A
				pixels[pixelIdx+1] = specs.BorderColor.B
				pixels[pixelIdx+2] = specs.BorderColor.G
				pixels[pixelIdx+3] = specs.BorderColor.R
			} else {
				pixels[pixelIdx+0] = specs.BodyColor.A
				pixels[pixelIdx+1] = specs.BodyColor.B
				pixels[pixelIdx+2] = specs.BodyColor.G
				pixels[pixelIdx+3] = specs.BodyColor.R
			}
			pixelIdx += 4
		}
	}

	err = texture.Update(nil, pixels, width*4)
	if err != nil {
		return nil, err
	}

	return texture, nil
}
