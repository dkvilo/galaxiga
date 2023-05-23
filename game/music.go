package game

import "github.com/veandco/go-sdl2/mix"

type Music struct {
	Music *mix.Music
}

func LoadMusic(filename string) (*Music, error) {
	music, err := mix.LoadMUS(filename)
	if err != nil {
		return nil, err
	}

	return &Music{Music: music}, nil
}

func (m *Music) Play() error {
	return m.Music.Play(-1)
}
