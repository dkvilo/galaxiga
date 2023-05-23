package game

const (
	WindowWidth            = 1920
	WindowHeight           = 1080
	PROJECTILE_SHOOT_DELAY = 10
	PROJECTILE_SPEED       = 10
)

var GlobalScore uint32 = 0
var GlobalPrevScore uint32 = 0

var GlobalAudioManager = NewAudioManager()
