package base

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"math/rand"
	"os"

	"github.com/go-gl/gl/v2.1/gl"
)

func Project(worldX, worldY float32) (screenX, screenY float32) {
	// Scale factors
	var scaleX float32 = 1.0 // 1080.0 / 1080.0
	var scaleY float32 = 1.0 // 720.0 / 720.0

	// Project world coordinates to screen coordinates
	screenX = worldX * scaleX
	screenY = worldY * scaleY

	return screenX, screenY
}

func Unproject(screenX, screenY float32) (worldX, worldY float32) {
	// Scale factors
	var scaleX float32 = 1.0 // 1080.0 / 1080.0
	var scaleY float32 = 1.0 // 720.0 / 720.0

	// Unproject screen coordinates to world coordinates
	worldX = screenX / scaleX
	worldY = screenY / scaleY

	return worldX, worldY
}

type Color struct {
	R, G, B, A uint8
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func RandomFloat(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func Lerp(a, b, t float32) float32 {
	return a + (b-a)*t
}

func Clamp(a, min, max float32) float32 {
	if a < min {
		return min
	}
	if a > max {
		return max
	}
	return a
}

func ClampInt(a, min, max int) int {
	if a < min {
		return min
	}
	if a > max {
		return max
	}
	return a
}

func Min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func Max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func LoadTexture(filename string) (uint32, error) {
	// Open the image file
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return 0, err
	}

	// Convert the image to RGBA format
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	// Generate a new texture ID
	var textureID uint32
	gl.GenTextures(1, &textureID)

	// Bind the texture
	gl.BindTexture(gl.TEXTURE_2D, textureID)

	fmt.Println("Texture ID:", textureID)

	// Upload the image data to the GPU
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)

	// Set texture parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	// Unbind the texture
	gl.BindTexture(gl.TEXTURE_2D, 0)

	return textureID, nil
}
