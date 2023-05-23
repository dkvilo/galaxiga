package game

import (
	"galaxiga/pkg/base"

	"github.com/go-gl/gl/v2.1/gl"
)

type Projectile struct {
	Rect
	Vel   Velocity
	Flags Flags
}

var GlobalProjectiles = make([]Projectile, 0)

func ShootProjectile(e *Entity, vel Velocity, flags Flags) {

	if flags&FLAG_PLAYER == FLAG_PLAYER {
		GlobalAudioManager.Play("shoot")
	}

	p := Projectile{
		Rect: Rect{
			X: e.X + (e.W / 2) - 4,
			Y: e.Y + (e.H / 2) - 4,
			W: 8,
			H: 8,
		},
		Flags: flags,
		Vel:   vel,
	}

	GlobalProjectiles = append(GlobalProjectiles, p)
	e.ShootDelay = PROJECTILE_SHOOT_DELAY
}

func (p *Projectile) CollidesWith(e *Projectile) bool {
	return p.X < e.X+e.W &&
		p.X+p.W > e.X &&
		p.Y < e.Y+e.H &&
		p.Y+p.H > e.Y
}

func UpdateProjectiles(pe *ParticleEmitter) {
	for i := 0; i < len(GlobalProjectiles); i++ {
		p := &GlobalProjectiles[i]
		p.X += p.Vel.X
		p.Y += p.Vel.Y

		if p.Y < 0 || p.Y > WindowHeight { // off screen
			GlobalProjectiles = append(GlobalProjectiles[:i], GlobalProjectiles[i+1:]...)
			i--
		}
	}
}

func DrawProjectiles(pe *ParticleEmitter) {
	const size = 5
	const life = 100

	for _, p := range GlobalProjectiles {

		if p.Flags&FLAG_PLAYER == FLAG_PLAYER {
			color := base.Color{R: uint8(base.RandomInt(0, 255)), G: uint8(base.RandomInt(0, 255)), B: uint8(base.RandomInt(0, 255)), A: 255}
			pe.Emit(p.X, p.Y, size, size, 1, life, color)
		}

		if false && p.Flags&FLAG_ENEMY == FLAG_ENEMY {
			pe.Emit(p.X, p.Y, size, size, 1, life, base.Color{R: 255, G: 0, B: 0, A: 255})
		}
	}

	gl.Begin(gl.QUADS)
	for _, p := range GlobalProjectiles {
		if p.Flags&FLAG_ENEMY == FLAG_ENEMY {
			gl.Color4ub(255, 0, 0, 255)
		} else {
			gl.Color4ub(255, 255, 255, 255)
		}
		gl.Vertex2f(p.X, p.Y)
		gl.Vertex2f(p.X+p.W, p.Y)
		gl.Vertex2f(p.X+p.W, p.Y+p.H)
		gl.Vertex2f(p.X, p.Y+p.H)
	}
	gl.End()
}
