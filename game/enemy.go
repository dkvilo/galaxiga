package game

import (
	"fmt"
	"galaxiga/pkg/base"
)

var GlobalEnemies = make([]Entity, 0)

func SpawnEnemy(tm *TextureManager, count int) {
	for i := 0; i < count; i++ {

		x := float32(base.RandomInt(0, WindowWidth/64)) * 64
		y := float32(64)

		randDel := base.RandomInt(0, 100)
		e := Entity{
			Rect: Rect{
				X: x,
				Y: y,
				W: 42,
				H: 42,
			},
			Health:        100,
			Vel:           Velocity{X: 0, Y: 0},
			ShootDelay:    randDel, // randomize shoot delay, so they don't all shoot at once
			ShootDelayMax: 100,
			MoveDelay:     100,
			Flags:         FLAG_ENEMY,
			Texture:       tm.GetTexture(fmt.Sprintf("enemy_%d", base.RandomInt(1, 3))),
		}
		GlobalEnemies = append(GlobalEnemies, e)
	}
}

func UpdateEnemies(player *Entity, pe *ParticleEmitter, tm *TextureManager) {
	for i := 0; i < len(GlobalEnemies); i++ {
		e := &GlobalEnemies[i]
		e.X += e.Vel.X
		e.Y += e.Vel.Y

		if e.ShootDelay > 0 {
			e.ShootDelay--
		}

		if e.MoveDelay > 0 {
			e.MoveDelay--
		}

		if e.Health <= 0 {
			GlobalEnemies = append(GlobalEnemies[:i], GlobalEnemies[i+1:]...)
			i--
		}

		if e.ShootDelay == 0 {
			ShootProjectile(e, Velocity{X: 0, Y: PROJECTILE_SPEED / 2}, FLAG_ENEMY)
			e.ShootDelay = 100
		}

		// move randomly but do't go off screen
		if e.MoveDelay == 0 {
			e.Vel.X = base.RandomFloat(-1, 1)
			e.Vel.Y = base.RandomFloat(-1, 1)
			e.MoveDelay = 100

			if e.X < 0 {
				e.X = 0
			}

			if e.X > float32(WindowWidth-e.W) {
				e.X = float32(WindowWidth - e.W)
			}

			if e.Y < 0 {
				e.Y = 0
			}

			if e.Y > float32(WindowHeight-e.H) {
				e.Y = float32(WindowHeight - e.H)
			}

			e.Vel.X *= 2
			e.Vel.Y *= 2
		}

		if e.Y > float32(WindowHeight) {
			player.Health -= 20
			e.Health = 0
		}

		if e.Collides(player) {
			e.Health = 0
			player.Health = 0
			pe.EmitExplosion(e.X, e.Y, 10, 10, 40, 100)
			pe.EmitExplosion(player.X, player.Y, 10, 10, 40, 100)
			GlobalAudioManager.Play("explosion")
		}

		for j := 0; j < len(GlobalProjectiles); j++ {
			p := &GlobalProjectiles[j]
			if e.HasHit(p) {
				GlobalAudioManager.Play("explosion")
				e.Health = 0
				GlobalScore += 10
				pe.EmitExplosion(p.X, p.Y, 10, 10, 40, 100)
				// chance to spawn a powerup when an enemy dies (1 in 20)
				if base.RandomInt(0, 20) == 0 {
					SpawnRandomUpgrade(tm, e.X, e.Y)
				}

				if player.Flags&FLAG_BULLET_PIERCER == 0 {
					GlobalProjectiles = append(GlobalProjectiles[:j], GlobalProjectiles[j+1:]...)
					j--
				}
			}

			if player.HasHit(p) {
				if player.Flags&FLAG_SHIELD != FLAG_SHIELD {
					if player.Health > 0 {
						player.Health -= 10
					}

					if GlobalScore > 0 {
						GlobalScore -= 10
					}

				} else {
					player.Flags &= ^FLAG_SHIELD
				}
				GlobalProjectiles = append(GlobalProjectiles[:j], GlobalProjectiles[j+1:]...)
				j--
			}
		}

	}
}

func DrawEnemies() {
	for _, e := range GlobalEnemies {
		e.Draw()
		DrawShootDelayCoolDown(&e)
	}
}
