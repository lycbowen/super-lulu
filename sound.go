package main

import (
	"encoding/binary"
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	audioSampleRate = 44100
	maxSoundPlayers = 24
)

type soundID int

const (
	soundJump soundID = iota
	soundCollect
	soundPowerUp
	soundShoot
	soundHurt
	soundBossHit
	soundBossDefeat
	soundLevelClear
)

type SoundManager struct {
	context *audio.Context
	sounds  map[soundID][]byte
	players []*audio.Player
	volume  float64
}

// newSoundManager 生成一组轻量的程序音效，避免项目暂时依赖外部 wav 文件。
func newSoundManager() *SoundManager {
	context := audio.CurrentContext()
	if context == nil {
		context = audio.NewContext(audioSampleRate)
	}
	return &SoundManager{
		context: context,
		sounds: map[soundID][]byte{
			soundJump:       synthChirp(0.12, 220, 620, 0.42, waveSine),
			soundCollect:    synthChirp(0.09, 820, 1320, 0.36, waveSine),
			soundPowerUp:    synthArpeggio([]float64{523.25, 659.25, 783.99, 1046.50}, 0.06, 0.34),
			soundShoot:      synthChirp(0.10, 760, 420, 0.33, waveSquare),
			soundHurt:       synthNoiseChirp(0.20, 250, 85, 0.40),
			soundBossHit:    synthChirp(0.14, 160, 70, 0.44, waveSquare),
			soundBossDefeat: synthNoiseChirp(0.48, 220, 45, 0.48),
			soundLevelClear: synthArpeggio([]float64{392.00, 523.25, 659.25, 783.99, 1046.50}, 0.08, 0.36),
		},
		volume: 0.55,
	}
}

func (s *SoundManager) Play(id soundID) {
	if s == nil {
		return
	}
	data, ok := s.sounds[id]
	if !ok {
		return
	}
	s.cleanupPlayers()
	if len(s.players) >= maxSoundPlayers {
		_ = s.players[0].Close()
		s.players = s.players[1:]
	}
	player := s.context.NewPlayerFromBytes(data)
	player.SetVolume(s.volume)
	player.Play()
	s.players = append(s.players, player)
}

func (s *SoundManager) Update() {
	if s == nil {
		return
	}
	s.cleanupPlayers()
}

func (s *SoundManager) cleanupPlayers() {
	active := s.players[:0]
	for _, player := range s.players {
		if player.IsPlaying() {
			active = append(active, player)
			continue
		}
		_ = player.Close()
	}
	s.players = active
}

func (g *Game) playSound(id soundID) {
	if g.sound == nil {
		return
	}
	g.sound.Play(id)
}

type waveShape int

const (
	waveSine waveShape = iota
	waveSquare
)

func synthChirp(seconds, startFreq, endFreq, volume float64, shape waveShape) []byte {
	total := int(seconds * audioSampleRate)
	out := make([]byte, total*4)
	for i := 0; i < total; i++ {
		t := float64(i) / float64(total)
		freq := startFreq + (endFreq-startFreq)*t
		phase := 2 * math.Pi * freq * float64(i) / audioSampleRate
		sample := math.Sin(phase)
		if shape == waveSquare {
			if sample >= 0 {
				sample = 1
			} else {
				sample = -1
			}
		}
		writeSample(out, i, sample*volume*envelope(t))
	}
	return out
}

func synthArpeggio(notes []float64, noteSeconds, volume float64) []byte {
	total := int(float64(len(notes)) * noteSeconds * audioSampleRate)
	out := make([]byte, total*4)
	noteSamples := int(noteSeconds * audioSampleRate)
	for i := 0; i < total; i++ {
		note := minInt(len(notes)-1, i/noteSamples)
		t := float64(i%noteSamples) / float64(noteSamples)
		phase := 2 * math.Pi * notes[note] * float64(i) / audioSampleRate
		writeSample(out, i, math.Sin(phase)*volume*envelope(t))
	}
	return out
}

func synthNoiseChirp(seconds, startFreq, endFreq, volume float64) []byte {
	total := int(seconds * audioSampleRate)
	out := make([]byte, total*4)
	var seed uint32 = 0x12345678
	for i := 0; i < total; i++ {
		seed = seed*1664525 + 1013904223
		noise := float64(int(seed>>16)&0xffff)/32768 - 1
		t := float64(i) / float64(total)
		freq := startFreq + (endFreq-startFreq)*t
		phase := 2 * math.Pi * freq * float64(i) / audioSampleRate
		tone := math.Sin(phase)
		writeSample(out, i, (tone*0.65+noise*0.35)*volume*envelope(t))
	}
	return out
}

func envelope(t float64) float64 {
	attack := math.Min(t/0.08, 1)
	release := math.Pow(1-t, 1.4)
	return attack * release
}

func writeSample(out []byte, index int, sample float64) {
	sample = clamp(sample, -1, 1)
	value := int16(sample * 32767)
	offset := index * 4
	binary.LittleEndian.PutUint16(out[offset:], uint16(value))
	binary.LittleEndian.PutUint16(out[offset+2:], uint16(value))
}
