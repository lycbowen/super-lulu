package main

import (
	"embed"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*.png
var embeddedAssets embed.FS

func loadAssets() (Assets, error) {
	playerImage, err := loadImage("assets/lulu.png")
	if err != nil {
		return Assets{}, fmt.Errorf("load player image: %w", err)
	}
	iceCreamImage, err := loadImage("assets/icecream.png")
	if err != nil {
		return Assets{}, fmt.Errorf("load ice cream image: %w", err)
	}
	orangeImage, err := loadImage("assets/orange.png")
	if err != nil {
		return Assets{}, fmt.Errorf("load orange image: %w", err)
	}
	bossImage, err := loadImage("assets/niuniu.png")
	if err != nil {
		return Assets{}, fmt.Errorf("load boss image: %w", err)
	}
	return Assets{
		Player:   playerImage,
		IceCream: iceCreamImage,
		Orange:   orangeImage,
		Boss:     bossImage,
	}, nil
}

func loadImage(path string) (*ebiten.Image, error) {
	f, err := embeddedAssets.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}
