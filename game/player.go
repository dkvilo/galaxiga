package game

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

func PlayerController(e *Entity) {
	const speed = 0.2

	if e.ShootDelay > 0 {
		e.ShootDelay--
	}

	key := sdl.GetKeyboardState()

	if key[sdl.SCANCODE_W] == 1 {
		e.Vel.Y -= 5
	}

	if key[sdl.SCANCODE_S] == 1 {
		e.Vel.Y += 5
	}

	if key[sdl.SCANCODE_A] == 1 {
		e.Vel.X -= 5
	}

	if key[sdl.SCANCODE_D] == 1 {
		e.Vel.X += 5
	}

	e.Vel.X *= 0.9
	e.Vel.Y *= 0.9

	e.X += e.Vel.X * speed
	e.Y += e.Vel.Y * speed

	if e.X < 0 {
		e.X = 0
	}

	if e.X > WindowWidth-e.W {
		e.X = WindowWidth - e.W
	}

	if e.Y < 0 {
		e.Y = 0
	}

	if e.Y > WindowHeight-e.H {
		e.Y = WindowHeight - e.H
	}

	if key[sdl.SCANCODE_SPACE] == 1 {

		if e.ShootDelay > 0 {
			return
		}

		if e.Flags&FLAG_DOUBLE_SHOT == FLAG_DOUBLE_SHOT {
			ShootProjectile(e, Velocity{X: -1, Y: -PROJECTILE_SPEED}, FLAG_PLAYER)
			ShootProjectile(e, Velocity{X: 1, Y: -PROJECTILE_SPEED}, FLAG_PLAYER)
		} else {
			if false {
				radius := 2
				for i := 0; i < 360; i += 32 {
					ShootProjectile(e, Velocity{X: float32(float64(radius) * math.Cos(float64(i))), Y: float32(float64(radius) * math.Sin(float64(i)))}, FLAG_PLAYER)
				}
			}
			ShootProjectile(e, Velocity{X: 0, Y: -PROJECTILE_SPEED}, FLAG_PLAYER)
		}
	}
}
