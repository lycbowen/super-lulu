package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	playerImage, err := loadImage("assets/lulu.png")
	if err != nil {
		log.Fatalf("failed to load assets/lulu.png: %v", err)
	}
	iceCreamImage, err := loadImage("assets/icecream.png")
	if err != nil {
		log.Fatalf("failed to load assets/icecream.png: %v", err)
	}
	orangeImage, err := loadImage("assets/orange.png")
	if err != nil {
		log.Fatalf("failed to load assets/orange.png: %v", err)
	}
	bossImage, err := loadImage("assets/niuniu.png")
	if err != nil {
		log.Fatalf("failed to load assets/niuniu.png: %v", err)
	}

	g := newGame(playerImage, iceCreamImage, orangeImage, bossImage)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Super Lulu")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func newGame(playerImage, iceCreamImage, orangeImage, bossImage *ebiten.Image) *Game {
	g := &Game{
		player:        &Player{Image: playerImage, Facing: 1},
		iceCream:      iceCreamImage,
		orange:        orangeImage,
		bossImage:     bossImage,
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
		levels:        newLevels(),
		mode:          modeMenu,
		unlockedLevel: 0,
	}
	g.loadLevel(0)
	return g
}
