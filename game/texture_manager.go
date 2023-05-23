package game

import (
	"fmt"
	"galaxiga/pkg/base"
)

type TextureManager struct {
	textures map[string]uint32
}

func NewTextureManager() *TextureManager {
	return &TextureManager{
		textures: make(map[string]uint32),
	}
}

func (tm *TextureManager) GetTexture(alias string) uint32 {
	if texture, ok := tm.textures[alias]; ok {
		return texture
	}

	panic(fmt.Sprintf("texture %s not loaded", alias))
}

func (tm *TextureManager) LoadTexture(alias, filename string) uint32 {
	if texture, ok := tm.textures[filename]; ok {
		return texture
	}

	texture, err := base.LoadTexture(filename)
	if err != nil {
		panic(err)
	}

	tm.textures[alias] = texture
	return texture
}
