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
		log.Fatalf("加载素材失败：%v", err)
	}

	g, err := newGame(assets)
	if err != nil {
		log.Fatalf("创建游戏失败：%v", err)
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("超级噜噜")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func newGame(assets Assets) (*Game, error) {
	levelPack, err := newLevelPack()
	if err != nil {
		return nil, err
	}
	g := &Game{
		player:        &Player{Facing: 1},
		assets:        assets,
		sound:         newSoundManager(),
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
		levels:        levelPack.Levels,
		levelSource:   levelPack.Source,
		language:      languageEnglish,
		mode:          modeMenu,
		unlockedLevel: 0,
	}
	g.loadLevel(0)
	return g, nil
}
