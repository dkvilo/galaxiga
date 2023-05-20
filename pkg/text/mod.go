package text

import (
	"cgo/pkg/base"
	"image"
	"image/color"
	"io/ioutil"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/golang/freetype"
)

// Kids, don't do this. It's bad.
const filepath = "Karmina-Bold.ttf"

type Text struct {
	text        string
	size        int
	color       color.Color
	shadowColor color.Color
	texture     uint32
	shadow      uint32
	Width       int
	Height      int
	X, Y        int32
}

func createTextTexture(text string, size int, clr color.Color) (uint32, int, int, error) {

	fontBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return 0, 0, 0, err
	}

	fontType, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return 0, 0, 0, err
	}

	const DPI = 72
	img := image.NewRGBA(image.Rect(0, 0, 1024, 1024))

	c := freetype.NewContext()
	c.SetDPI(DPI)
	c.SetFont(fontType)
	c.SetFontSize(float64(size))
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(clr))

	pt := freetype.Pt(0, size)
	_, err = c.DrawString(text, pt)
	if err != nil {
		return 0, 0, 0, err
	}

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(img.Bounds().Dx()), int32(img.Bounds().Dy()), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))

	return texture, img.Bounds().Dx(), img.Bounds().Dy(), nil
}

func drawText(texture uint32, _x, _y int32, width, height int, opacity float32) {
	x, y := base.Unproject(float32(_x), float32(_y))

	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.Begin(gl.QUADS)
	gl.Color4f(1.0, 1.0, 1.0, opacity)

	gl.TexCoord2f(0.0, 0.0)
	gl.Vertex2f(x, y)

	gl.TexCoord2f(1.0, 0.0)
	gl.Vertex2f(x+float32(width), y)

	gl.TexCoord2f(1.0, 1.0)
	gl.Vertex2f(x+float32(width), y+float32(height))

	gl.TexCoord2f(0.0, 1.0)
	gl.Vertex2f(x, y+float32(height))

	gl.End()

	gl.Disable(gl.TEXTURE_2D)
}

func (t *Text) UpdateText(newText string) {
	t.text = newText
	gl.DeleteTextures(1, &t.texture)
	t.texture, t.Width, t.Height, _ = createTextTexture(t.text, t.size, t.color)

	gl.DeleteTextures(1, &t.shadow)
	t.shadow, t.Width, t.Height, _ = createTextTexture(t.text, t.size, t.shadowColor)
}

func (t *Text) UpdateColor(newColor color.Color) {
	t.color = newColor
	gl.DeleteTextures(1, &t.texture)
	t.texture, t.Width, t.Height, _ = createTextTexture(t.text, t.size, t.color)
}

func createTextTextureDouble(text string, size int, clr color.Color, shadowClr color.Color) (uint32, uint32, int, int, error) {
	texture, width, height, err := createTextTexture(text, size, clr)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	// Create the shadow texture
	shadow, _, _, err := createTextTexture(text, size, shadowClr)
	if err != nil {
		gl.DeleteTextures(1, &texture)
		return 0, 0, 0, 0, err
	}

	return texture, shadow, width, height, nil
}

func NewText(text string, size int, clr color.Color, shadowClr color.Color, x, y int32) (*Text, error) {
	texture, shadow, width, height, err := createTextTextureDouble(text, size, clr, shadowClr) // Updated function
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &Text{
		text:        text,
		size:        size,
		color:       clr,
		shadowColor: shadowClr,
		texture:     texture,
		shadow:      shadow,
		Width:       width,
		Height:      height,
		X:           x,
		Y:           y,
	}, nil
}

func (t *Text) Draw() {
	drawText(t.shadow, t.X+2, t.Y+2, t.Width, t.Height, 1.0)
	drawText(t.texture, t.X, t.Y, t.Width, t.Height, 1.0)
}
