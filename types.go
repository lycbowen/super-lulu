package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Rect struct {
	X float64
	Y float64
	W float64
	H float64
}

func (r Rect) Intersects(o Rect) bool {
	return r.X < o.X+o.W && r.X+r.W > o.X && r.Y < o.Y+o.H && r.Y+r.H > o.Y
}

type Player struct {
	X        float64
	Y        float64
	VX       float64
	VY       float64
	OnGround bool
	Facing   int
	Big      bool
}

func (p *Player) Rect() Rect {
	if p.Big {
		return Rect{X: p.X + 2, Y: p.Y - 12, W: 44, H: 70}
	}
	return Rect{X: p.X + 7, Y: p.Y + 4, W: 34, H: 54}
}

type Collectible struct {
	Rect      Rect
	Collected bool
}

type PowerUp struct {
	Rect      Rect
	Collected bool
}

type OrangePower struct {
	Rect      Rect
	Collected bool
}

type Projectile struct {
	Rect        Rect
	VX          float64
	StartX      float64
	MaxDistance float64
	Active      bool
}

type Boss struct {
	Rect           Rect
	Hitbox         Rect
	BaseY          float64
	VY             float64
	OnGround       bool
	MinX           float64
	MaxX           float64
	Speed          float64
	ChargeSpeed    float64
	ChargeDistance float64
	ChargeDir      float64
	Facing         int
	ChargeStartX   float64
	SlamTargetX    float64
	SlamTargetY    float64
	SlamHangTime   int
	AttackCooldown int
	AttackPattern  int
	HitCooldown    int
	JumpCooldown   int
	RespawnGrace   int
	Timer          int
	State          bossState
	Aggro          bool
	HP             int
	MaxHP          int
	Active         bool
}

type Enemy struct {
	Rect  Rect
	MinX  float64
	MaxX  float64
	Speed float64
}

type Theme struct {
	Name string
	Sky  [3]uint8
	Hill [3]uint8
	Base [3]uint8
	Top  [3]uint8
	Trim [3]uint8
}

type Level struct {
	Name      string
	Subtitle  string
	Width     float64
	Spawn     image.Point
	Platforms []Rect
	Collect   []Collectible
	PowerUps  []PowerUp
	Oranges   []OrangePower
	Enemies   []Enemy
	Boss      *Boss
	Goal      Rect
	Theme     Theme
}

type Assets struct {
	Player   *ebiten.Image
	IceCream *ebiten.Image
	Orange   *ebiten.Image
	Boss     *ebiten.Image
}

type InputState struct {
	MoveLeft    bool
	MoveRight   bool
	Jump        bool
	Shoot       bool
	Start       bool
	Pause       bool
	Back        bool
	Restart     bool
	Menu        bool
	LevelUp     bool
	LevelDown   bool
	ToggleDebug bool
	DebugLevel  int
}

type Game struct {
	player        *Player
	assets        Assets
	input         InputState
	rng           *rand.Rand
	levels        []Level
	level         Level
	mode          gameMode
	selectedLevel int
	currentLevel  int
	unlockedLevel int
	camera        float64
	score         int
	levelScore    int
	falls         int
	hasWeapon     bool
	showDebug     bool
	projectiles   []Projectile
	shotCooldown  int
	invulnerable  int
}
