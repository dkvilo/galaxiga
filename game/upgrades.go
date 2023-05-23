package game

import (
	"fmt"
	"galaxiga/pkg/base"
	"math/rand"

	"github.com/go-gl/gl/v2.1/gl"
)

type Upgrade struct {
	Rect
	Vel     Velocity
	Flags   int32
	Texture uint32
}

func (e *Upgrade) Draw() {
	worldX, worldY := base.Unproject(e.X, e.Y)
	worldWidth, worldHeight := base.Unproject(e.X+e.W, e.Y+e.H)

	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, e.Texture)

	gl.Begin(gl.QUADS)
	gl.Color3ub(255, 255, 255)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(worldX, worldY)
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(worldWidth, worldY)
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(worldWidth, worldHeight)
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(worldX, worldHeight)
	gl.End()
	gl.Disable(gl.TEXTURE_2D)
}

func SpawnRandomUpgrade(tm *TextureManager, x, y float32) {
	upgrade1 := Upgrade{
		Rect: Rect{
			X: x,
			Y: y,
			W: 32,
			H: 32,
		},
		Vel: Velocity{
			X: 0,
			Y: 1,
		},
		Flags:   int32(FLAG_SHIELD),
		Texture: tm.GetTexture("shiled"),
	}

	upgrade2 := Upgrade{
		Rect: Rect{
			W: 32,
			H: 32,
		},
		Vel: Velocity{
			X: 0,
			Y: 1,
		},
		Flags:   int32(FLAG_BULLET_PIERCER),
		Texture: tm.GetTexture("pierce"),
	}

	upgrade3 := Upgrade{
		Rect: Rect{
			W: 32,
			H: 32,
		},
		Vel: Velocity{
			X: 0,
			Y: 1,
		},
		Flags:   int32(FLAG_DOUBLE_SHOT),
		Texture: tm.GetTexture("double_shoot"),
	}

	randomUpgrade := rand.Intn(3)

	var upgrade Upgrade

	if randomUpgrade == 0 {
		upgrade = upgrade1
	} else if randomUpgrade == 1 {
		upgrade = upgrade2
	} else if randomUpgrade == 2 {
		upgrade = upgrade3
	}

	upgrade.X = rand.Float32() * WindowWidth
	upgrade.Y = rand.Float32() * WindowHeight

	GlobalUpgrades = append(GlobalUpgrades, upgrade)
}

func (u *Upgrade) Update() {
	if u.Y > WindowHeight {
		u.Vel.Y = 0
		u.Y = 0

		fmt.Println("Upgrade removed")

		u = nil
	} else {
		u.Vel.Y += 0.01
		u.Y += u.Vel.Y
	}
}

func RemoveOffScreenUpgrades() {
	for i := len(GlobalUpgrades) - 1; i >= 0; i-- {
		if GlobalUpgrades[i].Y > WindowHeight {
			GlobalUpgrades[i] = GlobalUpgrades[len(GlobalUpgrades)-1]
			GlobalUpgrades = GlobalUpgrades[:len(GlobalUpgrades)-1]
		}
	}
}

func (u *Upgrade) Collides(target *Entity) bool {
	return (u.X < target.X+target.W && u.X+u.W > target.X && u.Y < target.Y+target.H && u.Y+u.H > target.Y)
}

var GlobalUpgrades = make([]Upgrade, 0)
