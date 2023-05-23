package game

import (
	"galaxiga/pkg/base"

	"github.com/go-gl/gl/v2.1/gl"
)

type Flags uint32

const (
	FLAG_NONE           Flags = 0
	FLAG_PLAYER         Flags = 1 << iota
	FLAG_ENEMY          Flags = 1 << iota
	FLAG_BULLET_PIERCER Flags = 1 << iota
	FLAG_SHIELD         Flags = 1 << iota
	FLAG_DOUBLE_SHOT    Flags = 1 << iota
)

type Rect struct {
	X, Y, W, H float32
}

type Velocity struct {
	X, Y float32
}

type Entity struct {
	Rect
	Flags         Flags
	Vel           Velocity
	Color         base.Color
	ShootDelay    int
	ShootDelayMax int
	Health        int
	MoveDelay     int
	Texture       uint32
}

func (e *Entity) Draw() {
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

	if e.Flags&FLAG_SHIELD == FLAG_SHIELD {
		gl.Begin(gl.QUADS)
		gl.Color4ub(0, 0, 255, 100)
		gl.Vertex2f(worldX-5, worldY-5)
		gl.Vertex2f(worldWidth+5, worldY-5)
		gl.Vertex2f(worldWidth+5, worldHeight+5)
		gl.Vertex2f(worldX-5, worldHeight+5)
		gl.End()
	}
}

func (e *Entity) SetPosition(x, y float32) {
	e.X = x
	e.Y = y
}

func (e *Entity) Collides(target *Entity) bool {
	return (e.X < target.X+target.W && e.X+e.W > target.X && e.Y < target.Y+target.H && e.Y+e.H > target.Y)
}

func (e *Entity) CollidesMap(mapWidth, mapHeight float32) bool {
	return (e.X < 0 || e.X+e.W > mapWidth || e.Y < 0 || e.Y+e.H > mapHeight)
}

func (e *Entity) HasHit(p *Projectile) bool {
	if e.Flags&FLAG_PLAYER == p.Flags&FLAG_PLAYER {
		return false
	}

	return (e.X < p.X+p.W && e.X+e.W > p.X && e.Y < p.Y+p.H && e.Y+e.H > p.Y)
}

func DrawShootDelayCoolDown(e *Entity) {

	gl.Begin(gl.QUADS)

	gl.Color4ub(0, 255, 0, 255)
	gl.Vertex2f(e.X, e.Y+e.H+5)
	gl.Vertex2f(e.X+e.W, e.Y+e.H+5)
	gl.Vertex2f(e.X+e.W, e.Y+e.H+10)
	gl.Vertex2f(e.X, e.Y+e.H+10)

	gl.Color4ub(255, 0, 0, 255)
	panelMaxWidth := 40

	coolDownBarWidth := float32(e.ShootDelay) * float32(panelMaxWidth) / float32(100)
	gl.Vertex2f(e.X, e.Y+e.H+5)
	gl.Vertex2f(e.X+coolDownBarWidth, e.Y+e.H+5)
	gl.Vertex2f(e.X+coolDownBarWidth, e.Y+e.H+10)
	gl.Vertex2f(e.X, e.Y+e.H+10)

	gl.End()
}

func DrawHealthBar(e *Entity) {
	offset := 10
	panelMaxWidth := 400

	gl.Begin(gl.QUADS)
	gl.Color4ub(255, 0, 0, 255)
	gl.Vertex2f(float32(offset), WindowHeight-30)
	gl.Vertex2f(float32(panelMaxWidth), WindowHeight-30)
	gl.Vertex2f(float32(panelMaxWidth), WindowHeight)
	gl.Vertex2f(float32(offset), WindowHeight)

	gl.Color4ub(0, 255, 0, 255)
	healthBarWidth := float32(e.Health) * float32(panelMaxWidth) / float32(100)
	gl.Vertex2f(float32(offset), WindowHeight-30)
	gl.Vertex2f(healthBarWidth, WindowHeight-30)
	gl.Vertex2f(healthBarWidth, WindowHeight)
	gl.Vertex2f(float32(offset), WindowHeight+float32(offset))
	gl.End()
}
