package main

import (
	"cgo/game"
	"cgo/pkg/base"
	"cgo/pkg/text"
	"fmt"
	"image/color"
	"log"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

var enemies = make([]game.Entity, 0)

func SpawnEnemy(tm *game.TextureManager, count int) {
	for i := 0; i < count; i++ {

		x := float32(base.RandomInt(0, game.WindowWidth/64)) * 64
		y := float32(64)

		randDel := base.RandomInt(0, 100)
		e := game.Entity{
			Rect: game.Rect{
				X: x,
				Y: y,
				W: 42,
				H: 42,
			},
			Health:        100,
			Vel:           game.Velocity{X: 0, Y: 0},
			ShootDelay:    randDel, // randomize shoot delay, so they don't all shoot at once
			ShootDelayMax: 100,
			MoveDelay:     100,
			Flags:         game.FLAG_ENEMY,
			Texture:       tm.GetTexture(fmt.Sprintf("enemy_%d", base.RandomInt(1, 3))),
		}
		enemies = append(enemies, e)
	}
}

var score uint32 = 0
var prevScore uint32 = 0

func UpdateEnemies(player *game.Entity, pe *game.ParticleEmitter, tm *game.TextureManager) {
	for i := 0; i < len(enemies); i++ {
		e := &enemies[i]
		e.X += e.Vel.X
		e.Y += e.Vel.Y

		if e.ShootDelay > 0 {
			e.ShootDelay--
		}

		if e.MoveDelay > 0 {
			e.MoveDelay--
		}

		if e.Health <= 0 {
			enemies = append(enemies[:i], enemies[i+1:]...)
			i--
		}

		if e.ShootDelay == 0 {
			game.ShootProjectile(e, game.Velocity{X: 0, Y: game.PROJECTILE_SPEED / 2}, game.FLAG_ENEMY)
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

			if e.X > float32(game.WindowWidth-e.W) {
				e.X = float32(game.WindowWidth - e.W)
			}

			if e.Y < 0 {
				e.Y = 0
			}

			if e.Y > float32(game.WindowHeight-e.H) {
				e.Y = float32(game.WindowHeight - e.H)
			}

			e.Vel.X *= 2
			e.Vel.Y *= 2
		}

		if e.Y > float32(game.WindowHeight) {
			player.Health -= 20
			e.Health = 0
		}

		if e.Collides(player) {
			e.Health = 0
			player.Health = 0
			pe.EmitExplosion(e.X, e.Y, 10, 10, 40, 100)
			pe.EmitExplosion(player.X, player.Y, 10, 10, 40, 100)
			game.GlobalAudioManager.Play("explosion")
		}

		for j := 0; j < len(game.GlobalProjectiles); j++ {
			p := &game.GlobalProjectiles[j]
			if e.HasHit(p) {
				game.GlobalAudioManager.Play("explosion")
				e.Health = 0
				score += 10
				pe.EmitExplosion(p.X, p.Y, 10, 10, 40, 100)
				// chance to spawn a powerup when an enemy dies (1 in 20)
				if base.RandomInt(0, 20) == 0 {
					game.SpawnRandomUpgrade(tm, e.X, e.Y)
				}

				if player.Flags&game.FLAG_BULLET_PIERCER == 0 {
					game.GlobalProjectiles = append(game.GlobalProjectiles[:j], game.GlobalProjectiles[j+1:]...)
					j--
				}
			}

			if player.HasHit(p) {
				if player.Flags&game.FLAG_SHIELD != game.FLAG_SHIELD {
					if player.Health > 0 {
						player.Health -= 10
					}

					if score > 0 {
						score -= 10
					}
				} else {
					player.Flags &= ^game.FLAG_SHIELD
				}
				game.GlobalProjectiles = append(game.GlobalProjectiles[:j], game.GlobalProjectiles[j+1:]...)
				j--
			}
		}

	}
}

func DrawEnemies() {
	for _, e := range enemies {
		e.Draw()
		DrawShootDelayCoolDown(&e)
	}
}

func PlayerController(e *game.Entity) {
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

	if e.X > game.WindowWidth-e.W {
		e.X = game.WindowWidth - e.W
	}

	if e.Y < 0 {
		e.Y = 0
	}

	if e.Y > game.WindowHeight-e.H {
		e.Y = game.WindowHeight - e.H
	}

	if key[sdl.SCANCODE_SPACE] == 1 {

		if e.ShootDelay > 0 {
			return
		}

		if e.Flags&game.FLAG_DOUBLE_SHOT == game.FLAG_DOUBLE_SHOT {
			game.ShootProjectile(e, game.Velocity{X: -1, Y: -game.PROJECTILE_SPEED}, game.FLAG_PLAYER)
			game.ShootProjectile(e, game.Velocity{X: 1, Y: -game.PROJECTILE_SPEED}, game.FLAG_PLAYER)
		} else {
			game.ShootProjectile(e, game.Velocity{X: 0, Y: -game.PROJECTILE_SPEED}, game.FLAG_PLAYER)
		}

	}

}

func DrawShootDelayCoolDown(e *game.Entity) {

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

func DrawHealthBar(e *game.Entity) {

	// Draw health bar on bottom left of screen (static)
	gl.Begin(gl.QUADS)

	offset := 10

	panelMaxWidth := 400

	gl.Color4ub(255, 0, 0, 255)
	gl.Vertex2f(float32(offset), game.WindowHeight-30)
	gl.Vertex2f(float32(panelMaxWidth), game.WindowHeight-30)
	gl.Vertex2f(float32(panelMaxWidth), game.WindowHeight)
	gl.Vertex2f(float32(offset), game.WindowHeight)

	gl.Color4ub(0, 255, 0, 255)

	healthBarWidth := float32(e.Health) * float32(panelMaxWidth) / float32(100)
	gl.Vertex2f(float32(offset), game.WindowHeight-30)
	gl.Vertex2f(healthBarWidth, game.WindowHeight-30)
	gl.Vertex2f(healthBarWidth, game.WindowHeight)
	gl.Vertex2f(float32(offset), game.WindowHeight+float32(offset))

	gl.End()

}

func main() {
	var winTitle string = "Space Shooter?"
	var window *sdl.Window
	var context sdl.GLContext
	var event sdl.Event
	var running bool
	var err error

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, game.WindowWidth, game.WindowHeight, sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}

	defer window.Destroy()
	context, err = window.GLCreateContext()
	if err != nil {
		panic(err)
	}

	if err := mix.Init(mix.INIT_MP3 | mix.INIT_FLAC | mix.INIT_MOD | mix.INIT_OGG); err != nil {
		log.Fatal(err)
	}

	defer mix.Quit()

	if err := mix.OpenAudio(mix.DEFAULT_FREQUENCY, mix.DEFAULT_FORMAT, 2, 1024); err != nil {
		log.Println(err)
		return
	}

	defer mix.CloseAudio()

	defer sdl.GLDeleteContext(context)
	if err = gl.Init(); err != nil {
		panic(err)
	}

	tm := game.NewTextureManager()

	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(0.2, 0.2, 0.3, 1.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)

	player := game.Entity{
		Rect:          game.Rect{X: game.WindowWidth/2 - 16, Y: game.WindowHeight - 64, W: 32, H: 32},
		Flags:         game.FLAG_PLAYER,
		Vel:           game.Velocity{},
		Color:         base.Color{R: 255, G: 255, B: 255},
		ShootDelay:    10,
		ShootDelayMax: 10,
		Health:        100,
		MoveDelay:     0,
		Texture:       tm.LoadTexture("player", "player.png"),
	}

	txt, err := text.NewText("Score: 0", 32, color.RGBA{255, 255, 255, 255}, color.RGBA{0, 0, 0, 255}, 10, 10)
	if err != nil {
		panic(err)
	}

	prevScoreLabel := "Prev Score: 0"
	scoreLabel, err := text.NewText(prevScoreLabel, 32, color.RGBA{255, 255, 255, 255}, color.RGBA{0, 0, 0, 255}, 10, 50)
	if err != nil {
		panic(err)
	}

	pe := game.ParticleEmitter{}

	tm.LoadTexture("enemy_1", "enemy_1.png")
	tm.LoadTexture("enemy_2", "enemy_2.png")
	tm.LoadTexture("enemy_3", "enemy_3.png")

	tm.LoadTexture("shiled", "shiled.png")
	tm.LoadTexture("pierce", "pierce.png")

	tm.LoadTexture("double_shoot", "double_shoot.png")

	game.GlobalAudioManager.LoadAudio("shoot", "shoot.wav")
	game.GlobalAudioManager.LoadAudio("explosion", "boom.wav")
	game.GlobalAudioManager.LoadAudio("hit", "hurt.wav")
	game.GlobalAudioManager.LoadAudio("wave", "wave.wav")
	game.GlobalAudioManager.LoadAudio("lost", "lost.wav")

	game.GlobalAudioManager.SetVolume("shoot", 10)

	enemyCount := 5
	SpawnEnemy(tm, enemyCount)

	pe.IsRunning = true

	var currentTime, deltaTime, frameTime, lastTime, lastFrame uint64
	const MAX_FPS = 60
	running = true
	for running {

		currentTime = sdl.GetTicks64()
		deltaTime = currentTime - lastTime
		frameTime += deltaTime
		lastTime = currentTime

		if frameTime >= 1000/MAX_FPS {
			frameTime = 0
		}

		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseMotionEvent:

			case *sdl.MouseButtonEvent:

			case *sdl.KeyboardEvent:
				keyEvent := event.(*sdl.KeyboardEvent)
				if keyEvent.Type == sdl.KEYDOWN {
					if keyEvent.Keysym.Sym == sdl.K_ESCAPE {
						running = false
					}

				}
			}
		}

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		gl.Ortho(0, game.WindowWidth, game.WindowHeight, 0, -1, 1)
		gl.Viewport(0, 0, game.WindowWidth, game.WindowHeight)
		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()

		if len(enemies) == 0 {

			game.GlobalAudioManager.Play("wave")

			player.Health += 10
			player.ShootDelayMax -= 5
			player.ShootDelay -= 5

			enemyCount += 2
			SpawnEnemy(tm, enemyCount)
			for i := 0; i < enemyCount; i++ {
				enemies[i].ShootDelayMax -= 5
				enemies[i].ShootDelay -= 5
			}

		}

		game.DrawProjectiles(&pe)
		game.UpdateProjectiles(&pe)

		player.Draw()
		DrawEnemies()

		for i := 0; i < len(game.GlobalUpgrades); i++ {
			game.GlobalUpgrades[i].Draw()
			game.GlobalUpgrades[i].Update()

			// Check collision with player and apply upgrade
			if game.GlobalUpgrades[i].Collides(&player) {
				player.Flags |= game.Flags(game.GlobalUpgrades[i].Flags)

				game.GlobalUpgrades[i] = game.GlobalUpgrades[len(game.GlobalUpgrades)-1]
				game.GlobalUpgrades = game.GlobalUpgrades[:len(game.GlobalUpgrades)-1]
			}
		}

		UpdateEnemies(&player, &pe, tm)
		DrawShootDelayCoolDown(&player)

		PlayerController(&player)

		pe.Update(float32(deltaTime))
		pe.Draw(float32(deltaTime))

		txt.Draw()
		txt.UpdateText(fmt.Sprintf("Score: %d", score))
		scoreLabel.Draw()

		// Draw Player Health Bar down left corner
		DrawHealthBar(&player)

		if player.Health <= 0 {

			game.GlobalAudioManager.Play("lost")
			// remove all projectiles
			game.GlobalProjectiles = game.GlobalProjectiles[:0]
			// remove all enemies
			enemies = enemies[:0]
			// remove all upgrades
			game.GlobalUpgrades = game.GlobalUpgrades[:0]

			// reset player
			player.Health = 100
			player.Flags = game.FLAG_PLAYER | game.FLAG_SHIELD
			player.Rect.X = game.WindowWidth/2 - 16
			player.Rect.Y = game.WindowHeight - 64
			player.Vel.X = 0
			player.Vel.Y = 0

			prevScore = score
			prevScoreLabel = fmt.Sprintf("Prev Score: %d", prevScore)
			scoreLabel.UpdateText(prevScoreLabel)

			// reset score and enemy count
			score = 0
			enemyCount = 5
			SpawnEnemy(tm, enemyCount)

			// reset particle emitter
			pe.IsRunning = true

			// reset shoot delay for player
			player.ShootDelayMax = game.PROJECTILE_SHOOT_DELAY
			player.ShootDelay = game.PROJECTILE_SHOOT_DELAY
		}

		window.GLSwap()

		target_frame_time := uint64(1000 / MAX_FPS)
		frame_time := sdl.GetTicks64() - lastFrame
		if frame_time < target_frame_time {
			sdl.Delay(uint32(target_frame_time - frame_time))
		}
		lastFrame = currentTime

	}
}
