package game

import (
	"fmt"

	"github.com/veandco/go-sdl2/mix"
)

type Audio struct {
	Chunk *mix.Chunk
}

func loadAudio(filename string) (*Audio, error) {
	chunk, err := mix.LoadWAV(filename)
	if err != nil {
		return nil, fmt.Errorf("load WAV: %w", err)
	}

	return &Audio{Chunk: chunk}, nil
}

type AudioManager struct {
	Audios map[string]*Audio
}

func NewAudioManager() *AudioManager {
	return &AudioManager{
		Audios: make(map[string]*Audio),
	}
}

func (am *AudioManager) GetAudio(alias string) *mix.Chunk {
	if audio, ok := am.Audios[alias]; ok {
		return audio.Chunk
	}

	panic(fmt.Sprintf("audio %s not loaded", alias))
}

func (am *AudioManager) LoadAudio(alias, filename string) (*Audio, error) {
	audio, err := loadAudio(filename)
	if err != nil {
		return nil, err
	}
	am.Audios[alias] = audio

	return audio, nil
}

func (am *AudioManager) SetVolume(alias string, volume int) {
	audio := am.GetAudio(alias)
	audio.Volume(volume)
}

func (am *AudioManager) Play(alias string) {
	audio := am.GetAudio(alias)
	channel, err := audio.Play(-1, 0)
	if err != nil {
		fmt.Println("Error playing audio:", err)
		return
	}
	fmt.Println("Playing on channel", channel)
}
