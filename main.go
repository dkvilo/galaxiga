package main

import (
	"fmt"
	"galaxiga/game"
	"galaxiga/pkg/base"
	"galaxiga/pkg/text"
	"image/color"
	"log"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	STAR_COUNT = 1000
)

var stars [][]int32 = make([][]int32, STAR_COUNT)

func init() {
	for i := 0; i < STAR_COUNT; i++ {
		stars[i] = []int32{int32(base.RandomFloat(0, float32(game.WindowWidth))), int32(base.RandomFloat(0, float32(game.WindowHeight)))}
	}
}

func DrawBackorundStarfield(player *game.Entity, deltaTime float32) {

	worldWidth, worldHeight := base.Unproject(float32(game.WindowWidth), float32(game.WindowHeight))

	gl.Begin(gl.QUADS)
	gl.Color3ub(0, 0, 0)
	gl.Vertex2f(0, 0)
	gl.Vertex2f(worldWidth, 0)
	gl.Vertex2f(worldWidth, worldHeight)
	gl.Vertex2f(0, worldHeight)
	gl.End()

	gl.Begin(gl.POINTS)

	for i := 0; i < STAR_COUNT; i++ {
		var alpha uint8 = uint8(base.RandomInt(0, 180))
		if base.RandomInt(0, 3) == 0 {
			alpha = 255
		}
		gl.Color4ub(255, 255, 255, alpha)
		gl.Vertex2i(stars[i][0], stars[i][1])
	}

	gl.End()
}

func main() {
	var winTitle string = "Space Shooter?"
	var window *sdl.Window
	var context sdl.GLContext
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

	window.SetFullscreen(sdl.WINDOW_FULLSCREEN)

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

	player := game.Entity{
		Rect:          game.Rect{X: game.WindowWidth/2 - 16, Y: game.WindowHeight - 64, W: 32, H: 32},
		Flags:         game.FLAG_PLAYER,
		Vel:           game.Velocity{},
		Color:         base.Color{R: 255, G: 255, B: 255},
		ShootDelay:    10,
		ShootDelayMax: 10,
		Health:        100,
		MoveDelay:     0,
		Texture:       tm.LoadTexture("player", "res/player.png"),
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

	tm.LoadTexture("enemy_1", "res/enemy_1.png")
	tm.LoadTexture("enemy_2", "res/enemy_2.png")
	tm.LoadTexture("enemy_3", "res/enemy_3.png")

	tm.LoadTexture("shiled", "res/shiled.png")
	tm.LoadTexture("pierce", "res/pierce.png")

	tm.LoadTexture("double_shoot", "res/double_shoot.png")

	game.GlobalAudioManager.LoadAudio("shoot", "res/shoot.wav")
	game.GlobalAudioManager.LoadAudio("explosion", "res/boom.wav")
	game.GlobalAudioManager.LoadAudio("hit", "res/hurt.wav")
	game.GlobalAudioManager.LoadAudio("wave", "res/wave.wav")
	game.GlobalAudioManager.LoadAudio("lost", "res/lost.wav")

	game.GlobalAudioManager.SetVolume("shoot", 30)

	mix.VolumeMusic(12)
	bgMusic, err := game.LoadMusic("res/raining_bits.ogg")
	if err != nil {
		panic(err)
	}

	// bgMusic.Play()

	bgMusic.Music.Play(3)

	enemyCount := 5
	game.SpawnEnemy(tm, enemyCount)

	pe.IsRunning = true

	var currentTime, deltaTime, frameTime, lastTime, lastFrame uint64
	const MAX_FPS = 60
	var event sdl.Event

	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(0.2, 0.2, 0.3, 1.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)

	running := true
	for running {

		currentTime = sdl.GetTicks64()
		deltaTime = currentTime - lastTime
		frameTime += deltaTime
		lastTime = currentTime

		if frameTime >= 1000/MAX_FPS {
			frameTime = 0
		}

		target_frame_time := uint64(1000 / MAX_FPS)
		frame_time := sdl.GetTicks64() - lastFrame

		if frame_time < target_frame_time {
			sdl.Delay(uint32(target_frame_time - frame_time))
		}
		lastFrame = currentTime

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

		DrawBackorundStarfield(&player, float32(deltaTime))

		if len(game.GlobalEnemies) == 0 {

			game.GlobalAudioManager.Play("wave")

			pe.Reset()
			game.GlobalProjectiles = game.GlobalProjectiles[:0]

			player.Health += 10
			player.ShootDelayMax -= 5
			player.ShootDelay -= 5

			enemyCount += 2
			game.SpawnEnemy(tm, enemyCount)
			for i := 0; i < enemyCount; i++ {
				game.GlobalEnemies[i].ShootDelayMax -= 5
				game.GlobalEnemies[i].ShootDelay -= 5
			}

		}

		game.DrawProjectiles(&pe)
		game.UpdateProjectiles(&pe)

		player.Draw()
		game.DrawEnemies()

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

		game.PlayerController(&player)

		game.UpdateEnemies(&player, &pe, tm)
		game.DrawShootDelayCoolDown(&player)

		pe.Update(float32(deltaTime))
		pe.Draw(float32(deltaTime))

		txt.Draw()
		txt.UpdateText(fmt.Sprintf("Score: %d", game.GlobalScore))
		scoreLabel.Draw()

		game.DrawHealthBar(&player)

		game.ClearUpParticles(&pe)

		if player.Health <= 0 {

			game.GlobalAudioManager.Play("lost")
			game.GlobalProjectiles = game.GlobalProjectiles[:0]
			game.GlobalEnemies = game.GlobalEnemies[:0]
			game.GlobalUpgrades = game.GlobalUpgrades[:0]

			// reset player
			player.Health = 100
			player.Flags = game.FLAG_PLAYER | game.FLAG_SHIELD
			player.Rect.X = game.WindowWidth/2 - 16
			player.Rect.Y = game.WindowHeight - 64
			player.Vel.X = 0
			player.Vel.Y = 0

			game.GlobalPrevScore = game.GlobalScore
			prevScoreLabel = fmt.Sprintf("Prev Score: %d", game.GlobalPrevScore)
			scoreLabel.UpdateText(prevScoreLabel)

			// reset score and enemy count
			game.GlobalScore = 0
			enemyCount = 5
			game.SpawnEnemy(tm, enemyCount)

			// reset particle emitter
			pe.Reset()

			// reset shoot delay for player
			player.ShootDelayMax = game.PROJECTILE_SHOOT_DELAY
			player.ShootDelay = game.PROJECTILE_SHOOT_DELAY
		}

		fmt.Println("Enemise: ", len(game.GlobalEnemies))
		fmt.Println("Projectiles: ", len(game.GlobalProjectiles))
		fmt.Println("Upgrades: ", len(game.GlobalUpgrades))
		fmt.Println("Particles: ", len(pe.Particles))

		window.GLSwap()
	}
}
