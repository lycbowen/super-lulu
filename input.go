package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Update 把 Ebiten 的键盘输入转换成游戏语义，后续接手柄、改键位或做回放时只需要改这一层。
func (i *InputState) Update() {
	*i = InputState{}

	i.MoveLeft = ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft)
	i.MoveRight = ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight)
	i.Jump = ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp)

	i.Shoot = inpututil.IsKeyJustPressed(ebiten.KeyJ)
	i.Start = inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	i.Pause = inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyP)
	i.Back = inpututil.IsKeyJustPressed(ebiten.KeyEscape)
	i.Restart = inpututil.IsKeyJustPressed(ebiten.KeyR)
	i.Menu = inpututil.IsKeyJustPressed(ebiten.KeyM)
	i.LevelUp = inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW)
	i.LevelDown = inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS)
	i.PageLeft = inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA)
	i.PageRight = inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD)
	i.ToggleDebug = inpututil.IsKeyJustPressed(ebiten.KeyF3)
	i.ToggleLang = inpututil.IsKeyJustPressed(ebiten.KeyL)

	debugKeys := []ebiten.Key{
		ebiten.KeyDigit1,
		ebiten.KeyDigit2,
		ebiten.KeyDigit3,
		ebiten.KeyDigit4,
		ebiten.KeyDigit5,
		ebiten.KeyDigit6,
	}
	for index, key := range debugKeys {
		if inpututil.IsKeyJustPressed(key) {
			i.DebugLevel = index + 1
			return
		}
	}
}
