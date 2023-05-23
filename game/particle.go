package game

import (
	"galaxiga/pkg/base"

	"github.com/go-gl/gl/v2.1/gl"
)

type Particle struct {
	Rect
	Rotation   float32
	Vel        Velocity
	Color      base.Color
	Life       int32
	Age        int32
	HasGravity bool
	Gravity    float32
}

type ParticleEmitter struct {
	Particles []Particle
	IsRunning bool
}

func (pe *ParticleEmitter) Start() {
	pe.IsRunning = true
}

func (pe *ParticleEmitter) Stop() {
	pe.IsRunning = false
	pe.Particles = pe.Particles[:0]
}

func (pe *ParticleEmitter) Pause() {
	pe.IsRunning = false
}

func (pe *ParticleEmitter) Emit(x, y, w, h float32, count int32, life int32, col base.Color) {
	for i := 0; i < int(count); i++ {
		p := Particle{
			Rect: Rect{
				X: x,
				Y: y,
				W: w,
				H: h,
			},
			Rotation:   0,
			Vel:        Velocity{base.RandomFloat(-1, 1), base.RandomFloat(-1, 1)},
			Color:      col,
			Life:       life,
			HasGravity: true,
			Gravity:    0.1,
		}

		pe.Particles = append(pe.Particles, p)
	}
}

func (pe *ParticleEmitter) EmitExplosion(x, y, w, h float32, count int32, life int32) {
	for i := 0; i < int(count); i++ {
		explosionColor := base.Color{R: uint8(255), G: uint8(255), B: uint8(255), A: 255}
		p := Particle{
			Rect: Rect{
				X: x,
				Y: y,
				W: h,
				H: h,
			},
			Rotation:   0,
			Vel:        Velocity{base.RandomFloat(-1, 1), base.RandomFloat(-1, 1)},
			Color:      explosionColor,
			Life:       life,
			Age:        40,
			HasGravity: false,
			Gravity:    0.1,
		}

		if i%2 == 0 {
			p.Vel.X = base.RandomFloat(-2, -1)
		} else {
			p.Vel.X = base.RandomFloat(1, 2)
		}

		pe.Particles = append(pe.Particles, p)
	}
}

func (pe *ParticleEmitter) Update(dt float32) {
	if !pe.IsRunning {
		return
	}

	for i := 0; i < len(pe.Particles); i++ {
		p := &pe.Particles[i]
		p.X += p.Vel.X
		p.Y += p.Vel.Y
		p.Age += 2
		p.Rotation += 1

		if p.Age >= p.Life {
			p.Color.A = uint8(255 - (255 * (float32(p.Age) / float32(p.Life))))
			pe.Particles = append(pe.Particles[:i], pe.Particles[i+1:]...)
			i -= 1
		}

		if p.HasGravity {
			p.Vel.Y += p.Gravity
		}
	}
}

func (pe *ParticleEmitter) Draw(dt float32) {

	if !pe.IsRunning {
		return
	}
	gl.Begin(gl.QUADS)
	for _, p := range pe.Particles {

		dt := float32(p.Age) / float32(p.Life)
		p.Color.A = uint8(255 - (255 * dt))
		gl.Color4ub(p.Color.R, p.Color.G, p.Color.B, p.Color.A)

		gl.Vertex2f(p.X, p.Y)
		gl.Vertex2f(p.X+p.W, p.Y)
		gl.Vertex2f(p.X+p.W, p.Y+p.H)
		gl.Vertex2f(p.X, p.Y+p.H)

	}
	gl.End()
}

func ClearUpParticles(pe *ParticleEmitter) {
	for _, p := range pe.Particles {
		if p.Age >= p.Life {
			pe.Particles = append(pe.Particles[:0], pe.Particles[1:]...)
		}
	}
}

func (pe *ParticleEmitter) Reset() {
	pe.Particles = pe.Particles[:0]
}
