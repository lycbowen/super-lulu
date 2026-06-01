package main

const (
	screenWidth  = 960
	screenHeight = 540

	levelSelectPageSize = 4

	gravity      = 0.55
	moveAccel    = 0.6
	airAccel     = 0.36
	maxRunSpeed  = 5.2
	friction     = 0.82
	jumpVelocity = -12.8

	invulnerableFrames = 120
)

type gameMode int

const (
	modeMenu gameMode = iota
	modeLevelSelect
	modePlaying
	modePaused
	modeLevelClear
	modeGameClear
)

type bossState int

const (
	bossPatrol bossState = iota
	bossWarning
	bossCharge
	bossStompWarning
	bossStompRise
	bossStompHang
	bossStompFall
	bossRecover
)
