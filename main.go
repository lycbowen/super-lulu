package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	assets, err := loadAssets()
	if err != nil {
		log.Fatalf("failed to load assets: %v", err)
	}

	g := newGame(assets)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Super Lulu")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func newGame(assets Assets) *Game {
	g := &Game{
		player:        &Player{Facing: 1},
		assets:        assets,
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
		levels:        newLevels(),
		mode:          modeMenu,
		unlockedLevel: 0,
	}
	g.loadLevel(0)
	return g
}
