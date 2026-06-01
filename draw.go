package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	textv2 "github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"
)

var uiFace = textv2.NewGoXFace(basicfont.Face7x13)

func (g *Game) Draw(screen *ebiten.Image) {
	drawBackground(screen, g.camera, g.level)
	switch g.mode {
	case modeMenu:
		g.drawWorld(screen)
		g.drawMenu(screen)
	case modeLevelSelect:
		g.drawWorld(screen)
		g.drawLevelSelect(screen)
	case modePlaying:
		g.drawWorld(screen)
		g.drawHUD(screen)
	case modePaused:
		g.drawWorld(screen)
		g.drawHUD(screen)
		g.drawPause(screen)
	case modeLevelClear:
		g.drawWorld(screen)
		g.drawHUD(screen)
		g.drawLevelClear(screen)
	case modeGameClear:
		g.drawWorld(screen)
		g.drawGameClear(screen)
	}
	if g.showDebug {
		g.drawDebugOverlay(screen)
	}
}

func drawBackground(screen *ebiten.Image, camera float64, level Level) {
	base := level.Theme.Sky
	screen.Fill(color.RGBA{base[0], base[1], base[2], 255})
	for y := 0; y < screenHeight; y += 6 {
		t := float64(y) / screenHeight
		c := color.RGBA{
			R: uint8(clamp(float64(base[0])+20*t, 0, 255)),
			G: uint8(clamp(float64(base[1])-35*t, 0, 255)),
			B: uint8(clamp(float64(base[2])-45*t, 0, 255)),
			A: 255,
		}
		ebitenutil.DrawRect(screen, 0, float64(y), screenWidth, 6, c)
	}

	for i := 0; i < 8; i++ {
		x := math.Mod(float64(i*430)-camera*0.28, level.Width)
		if x < -180 {
			x += level.Width
		}
		y := float64(58 + (i%3)*35)
		drawCloud(screen, x, y)
	}

	hill := level.Theme.Hill
	for x := -math.Mod(camera*0.5, 180) - 90; x < screenWidth+180; x += 180 {
		ebitenutil.DrawRect(screen, x, 438, 120, 102, color.RGBA{hill[0], hill[1], hill[2], 95})
	}
}

func drawCloud(screen *ebiten.Image, x, y float64) {
	ebitenutil.DrawRect(screen, x, y+18, 116, 26, color.RGBA{255, 245, 190, 180})
	ebitenutil.DrawRect(screen, x+18, y+4, 40, 40, color.RGBA{255, 247, 201, 190})
	ebitenutil.DrawRect(screen, x+52, y, 48, 48, color.RGBA{255, 247, 201, 190})
}

func (g *Game) drawWorld(screen *ebiten.Image) {
	for _, platform := range g.level.Platforms {
		drawPlatform(screen, platform, g.camera, g.level.Theme)
	}
	for _, c := range g.level.Collect {
		if !c.Collected {
			drawCollectible(screen, c.Rect, g.camera)
		}
	}
	for _, p := range g.level.PowerUps {
		if !p.Collected {
			g.drawIceCream(screen, p.Rect, g.camera, 0.86)
		}
	}
	for _, o := range g.level.Oranges {
		if !o.Collected {
			g.drawOrange(screen, o.Rect, g.camera, 0.86)
		}
	}
	for _, p := range g.projectiles {
		g.drawIceCream(screen, p.Rect, g.camera, 0.64)
	}
	for _, e := range g.level.Enemies {
		drawEnemy(screen, e.Rect, g.camera)
	}
	if g.level.Boss != nil && g.level.Boss.Active {
		g.drawBoss(screen, g.level.Boss)
	}
	drawGoal(screen, g.level.Goal, g.camera)
	if g.bossBlocksGoal() {
		drawText(screen, "Defeat Niuniu first!", int(g.level.Goal.X-g.camera)-36, int(g.level.Goal.Y)-14, color.RGBA{104, 56, 19, 255})
	}
	if boss, ok := g.lockedBossArena(); ok {
		g.drawBossArenaGates(screen, boss)
	}
	g.drawPlayer(screen)
}

func drawPlatform(screen *ebiten.Image, r Rect, camera float64, theme Theme) {
	x := r.X - camera
	if x+r.W < -40 || x > screenWidth+40 {
		return
	}
	base := theme.Base
	top := theme.Top
	trim := theme.Trim
	ebitenutil.DrawRect(screen, x, r.Y, r.W, r.H, color.RGBA{base[0], base[1], base[2], 255})
	ebitenutil.DrawRect(screen, x, r.Y, r.W, 12, color.RGBA{top[0], top[1], top[2], 255})
	ebitenutil.DrawRect(screen, x+8, r.Y+12, r.W-16, 5, color.RGBA{trim[0], trim[1], trim[2], 255})
}

func drawCollectible(screen *ebiten.Image, r Rect, camera float64) {
	x := r.X - camera
	ebitenutil.DrawRect(screen, x+4, r.Y, 16, 24, color.RGBA{255, 231, 83, 255})
	ebitenutil.DrawRect(screen, x, r.Y+6, 24, 12, color.RGBA{255, 249, 156, 255})
}

func (g *Game) drawIceCream(screen *ebiten.Image, r Rect, camera float64, scale float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(r.X-camera, r.Y)
	screen.DrawImage(g.assets.IceCream, op)
}

func (g *Game) drawOrange(screen *ebiten.Image, r Rect, camera float64, scale float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(r.X-camera, r.Y)
	screen.DrawImage(g.assets.Orange, op)
}

func drawEnemy(screen *ebiten.Image, r Rect, camera float64) {
	x := r.X - camera
	ebitenutil.DrawRect(screen, x, r.Y+8, r.W, r.H-8, color.RGBA{255, 126, 65, 255})
	ebitenutil.DrawRect(screen, x+8, r.Y, r.W-16, 14, color.RGBA{255, 177, 66, 255})
	ebitenutil.DrawRect(screen, x+9, r.Y+18, 7, 7, color.RGBA{88, 45, 28, 255})
	ebitenutil.DrawRect(screen, x+r.W-16, r.Y+18, 7, 7, color.RGBA{88, 45, 28, 255})
}

func (g *Game) drawBoss(screen *ebiten.Image, boss *Boss) {
	x := boss.Rect.X - g.camera
	if x+boss.Rect.W < -80 || x > screenWidth+80 {
		return
	}
	op := &ebiten.DrawImageOptions{}
	scaleX := boss.Rect.W / 300
	scaleY := boss.Rect.H / 300
	if boss.Facing < 0 {
		op.GeoM.Scale(-scaleX, scaleY)
		op.GeoM.Translate(x+boss.Rect.W, boss.Rect.Y)
	} else {
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(x, boss.Rect.Y)
	}
	if boss.HitCooldown > 0 && (boss.HitCooldown/4)%2 == 0 {
		op.ColorScale.Scale(1.35, 0.75, 0.75, 1)
	}
	screen.DrawImage(g.assets.Boss, op)

	barW := boss.Rect.W
	barX := x
	barY := boss.Rect.Y - 22
	ebitenutil.DrawRect(screen, barX, barY, barW, 12, color.RGBA{101, 54, 34, 255})
	fill := barW * float64(boss.HP) / float64(boss.MaxHP)
	if fill > 2 {
		ebitenutil.DrawRect(screen, barX+1, barY+1, fill-2, 10, color.RGBA{255, 93, 76, 255})
	}
	if boss.State == bossWarning {
		g.drawBossWarning(screen, boss, x, "Niuniu charge!")
	}
	if boss.State == bossStompWarning || boss.State == bossStompRise || boss.State == bossStompHang || boss.State == bossStompFall {
		g.drawBossStompCue(screen, boss)
	}
	if boss.State == bossCharge {
		ebitenutil.DrawRect(screen, x-10, boss.Rect.Y+boss.Rect.H-18, boss.Rect.W+20, 12, color.RGBA{255, 93, 76, 95})
	}
}

func (g *Game) drawBossWarning(screen *ebiten.Image, boss *Boss, screenX float64, label string) {
	pulse := 1 + float64((60-boss.Timer)%12)/24
	iconSize := 26 * pulse
	iconX := screenX + boss.Rect.W/2 - iconSize/2
	iconY := boss.Rect.Y - 58 - iconSize*0.15
	ebitenutil.DrawRect(screen, iconX, iconY, iconSize, iconSize, color.RGBA{255, 72, 58, 230})
	drawText(screen, "!", int(iconX+iconSize/2-3), int(iconY+iconSize/2+5), color.RGBA{255, 245, 190, 255})
	drawCenteredTextAt(screen, label, int(screenX+boss.Rect.W/2), int(iconY-8), color.RGBA{120, 48, 25, 255})
}

func (g *Game) drawBossStompCue(screen *ebiten.Image, boss *Boss) {
	targetX := boss.SlamTargetX - g.camera
	if boss.State == bossStompWarning {
		playerRect := g.player.Rect()
		targetX = playerRect.X + playerRect.W/2 - boss.Rect.W/2 - g.camera
	}
	if boss.State == bossStompWarning {
		g.drawBossWarning(screen, boss, boss.Rect.X-g.camera, "Taishan stomp!")
	}
	alpha := uint8(115)
	if boss.State == bossStompHang || boss.State == bossStompFall {
		alpha = 185
	}
	ebitenutil.DrawRect(screen, targetX, boss.BaseY+boss.Rect.H-12, boss.Rect.W, 12, color.RGBA{255, 64, 54, alpha})
	ebitenutil.DrawRect(screen, targetX+boss.Rect.W*0.2, boss.BaseY+boss.Rect.H-24, boss.Rect.W*0.6, 8, color.RGBA{255, 220, 90, alpha})
}

func drawGoal(screen *ebiten.Image, r Rect, camera float64) {
	x := r.X - camera
	ebitenutil.DrawRect(screen, x, r.Y, 8, r.H, color.RGBA{126, 84, 45, 255})
	ebitenutil.DrawRect(screen, x+8, r.Y+8, 74, 42, color.RGBA{255, 223, 74, 255})
	ebitenutil.DrawRect(screen, x+18, r.Y+16, 48, 18, color.RGBA{255, 151, 53, 255})
}

func (g *Game) drawBossArenaGates(screen *ebiten.Image, boss *Boss) {
	drawArenaGate(screen, boss.ArenaMinX-g.camera)
	drawArenaGate(screen, boss.ArenaMaxX-g.camera)
	drawText(screen, "Boss arena locked", int(boss.ArenaMinX-g.camera)+18, 92, color.RGBA{104, 56, 19, 255})
}

func drawArenaGate(screen *ebiten.Image, x float64) {
	if x < -24 || x > screenWidth+24 {
		return
	}
	ebitenutil.DrawRect(screen, x-5, 76, 10, screenHeight-76, color.RGBA{104, 56, 19, 185})
	ebitenutil.DrawRect(screen, x-12, 76, 24, 8, color.RGBA{255, 205, 78, 230})
	for y := 104; y < screenHeight; y += 46 {
		ebitenutil.DrawRect(screen, x-10, float64(y), 20, 8, color.RGBA{255, 205, 78, 185})
	}
}

func (g *Game) drawPlayer(screen *ebiten.Image) {
	if g.invulnerable > 0 && (g.invulnerable/6)%2 == 0 {
		return
	}
	op := &ebiten.DrawImageOptions{}
	scale := 0.92
	drawY := g.player.Y
	if g.player.Big {
		scale = 1.18
		drawY -= 18
	}
	if g.player.Facing < 0 {
		op.GeoM.Scale(-scale, scale)
		op.GeoM.Translate(g.player.X-g.camera+50*scale, drawY)
	} else {
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(g.player.X-g.camera, drawY)
	}
	op.ColorScale.Scale(1.04, 1.04, 1.02, 1)
	screen.DrawImage(g.assets.Player, op)
}

func (g *Game) drawHUD(screen *ebiten.Image) {
	drawPanel(screen, 12, 12, 655, 52)
	drawText(screen, fmt.Sprintf("Level %d/%d", g.currentLevel+1, len(g.levels)), 28, 38, color.RGBA{96, 55, 23, 255})
	drawText(screen, fmt.Sprintf("Score %04d", g.score), 132, 38, color.RGBA{96, 55, 23, 255})
	drawText(screen, fmt.Sprintf("Resets %d", g.falls), 260, 38, color.RGBA{96, 55, 23, 255})
	weapon := "Ice cream: --"
	if g.hasWeapon {
		weapon = "Ice cream: J"
	}
	drawText(screen, weapon, 365, 38, color.RGBA{96, 55, 23, 255})
	size := "Size: normal"
	if g.player.Big {
		size = "Size: big"
	} else if g.invulnerable > 0 {
		size = "Size: blink"
	}
	drawText(screen, size, 500, 38, color.RGBA{96, 55, 23, 255})
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	drawPanel(screen, 258, 120, 444, 245)
	drawCenteredText(screen, "SUPER LULU", 172, color.RGBA{111, 61, 23, 255})
	drawCenteredText(screen, "A candy platform adventure", 210, color.RGBA{135, 82, 31, 255})
	drawCenteredText(screen, "Enter  Start", 260, color.RGBA{104, 56, 19, 255})
	drawCenteredText(screen, "Up/Down  Level Select", 292, color.RGBA{104, 56, 19, 255})
	drawCenteredText(screen, "A/D move   Space jump   J ice cream   1-6 debug levels", 330, color.RGBA{135, 82, 31, 255})
}

func (g *Game) drawLevelSelect(screen *ebiten.Image) {
	drawPanel(screen, 220, 92, 520, 350)
	drawCenteredText(screen, "CHOOSE A LEVEL", 142, color.RGBA{111, 61, 23, 255})
	for i, level := range g.levels {
		y := 198 + i*58
		locked := i > g.unlockedLevel
		prefix := "  "
		if i == g.selectedLevel {
			prefix = "> "
		}
		name := prefix + level.Name
		c := color.RGBA{104, 56, 19, 255}
		if locked {
			name = "  Locked"
			c = color.RGBA{155, 125, 92, 255}
		}
		drawText(screen, name, 285, y, c)
		if !locked {
			drawText(screen, level.Subtitle, 305, y+22, color.RGBA{135, 82, 31, 255})
		}
	}
	drawCenteredText(screen, "Enter start   Esc menu", 405, color.RGBA{104, 56, 19, 255})
}

func (g *Game) drawPause(screen *ebiten.Image) {
	drawPanel(screen, 310, 160, 340, 170)
	drawCenteredText(screen, "PAUSED", 210, color.RGBA{104, 56, 19, 255})
	drawCenteredText(screen, "P/Esc resume", 248, color.RGBA{104, 56, 19, 255})
	drawCenteredText(screen, "R restart   M menu", 282, color.RGBA{104, 56, 19, 255})
}

func (g *Game) drawLevelClear(screen *ebiten.Image) {
	drawPanel(screen, 270, 145, 420, 205)
	drawCenteredText(screen, "LEVEL CLEAR!", 195, color.RGBA{104, 56, 19, 255})
	drawCenteredText(screen, g.level.Name, 230, color.RGBA{104, 56, 19, 255})
	drawCenteredText(screen, fmt.Sprintf("Level score: %d", g.levelScore), 264, color.RGBA{104, 56, 19, 255})
	drawCenteredText(screen, "Enter next   M menu", 305, color.RGBA{104, 56, 19, 255})
}

func (g *Game) drawGameClear(screen *ebiten.Image) {
	drawPanel(screen, 250, 128, 460, 245)
	drawCenteredText(screen, "ALL LEVELS CLEAR!", 182, color.RGBA{104, 56, 19, 255})
	drawCenteredText(screen, "Lulu owns the candy kingdom.", 220, color.RGBA{104, 56, 19, 255})
	drawCenteredText(screen, fmt.Sprintf("Final score: %d", g.score), 258, color.RGBA{104, 56, 19, 255})
	drawCenteredText(screen, fmt.Sprintf("Total resets: %d", g.falls), 290, color.RGBA{104, 56, 19, 255})
	drawCenteredText(screen, "Enter replay   M menu", 330, color.RGBA{104, 56, 19, 255})
}

// drawDebugOverlay 显示调地图时最常用的世界坐标；鼠标世界 X 等于屏幕 X 加相机偏移。
func (g *Game) drawDebugOverlay(screen *ebiten.Image) {
	mouseX, mouseY := ebiten.CursorPosition()
	worldMouseX := float64(mouseX) + g.camera
	playerRect := g.player.Rect()
	playerScreenX := g.player.X - g.camera

	drawPanel(screen, 12, 76, 360, 150)
	lines := []string{
		"DEBUG F3",
		fmt.Sprintf("Player pos: X %.0f  Y %.0f", g.player.X, g.player.Y),
		fmt.Sprintf("Player rect: X %.0f Y %.0f W %.0f H %.0f", playerRect.X, playerRect.Y, playerRect.W, playerRect.H),
		fmt.Sprintf("Screen pos: X %.0f  Camera %.0f / %.0f", playerScreenX, g.camera, g.level.Width),
		fmt.Sprintf("Mouse screen: X %d  Y %d", mouseX, mouseY),
		fmt.Sprintf("Mouse world: X %.0f  Y %d", worldMouseX, mouseY),
	}
	for i, line := range lines {
		drawText(screen, line, 24, 102+i*20, color.RGBA{84, 42, 20, 255})
	}

	drawRectOutline(screen, playerRect.X-g.camera, playerRect.Y, playerRect.W, playerRect.H, color.RGBA{50, 240, 120, 230})
	ebitenutil.DrawRect(screen, float64(mouseX), 0, 1, screenHeight, color.RGBA{255, 60, 60, 150})
	ebitenutil.DrawRect(screen, 0, float64(mouseY), screenWidth, 1, color.RGBA{255, 60, 60, 150})
}

func drawRectOutline(screen *ebiten.Image, x, y, w, h float64, c color.Color) {
	ebitenutil.DrawRect(screen, x, y, w, 2, c)
	ebitenutil.DrawRect(screen, x, y+h-2, w, 2, c)
	ebitenutil.DrawRect(screen, x, y, 2, h, c)
	ebitenutil.DrawRect(screen, x+w-2, y, 2, h, c)
}

func drawPanel(screen *ebiten.Image, x, y, w, h float64) {
	ebitenutil.DrawRect(screen, x, y, w, h, color.RGBA{255, 245, 199, 220})
	ebitenutil.DrawRect(screen, x, y, w, 5, color.RGBA{255, 205, 78, 255})
	ebitenutil.DrawRect(screen, x, y+h-5, w, 5, color.RGBA{224, 126, 43, 255})
}

func drawCenteredText(screen *ebiten.Image, value string, y int, c color.Color) {
	x := screenWidth/2 - len(value)*7/2
	drawText(screen, value, x+2, y+2, color.RGBA{255, 246, 184, 255})
	drawText(screen, value, x, y, c)
}

func drawCenteredTextAt(screen *ebiten.Image, value string, centerX, y int, c color.Color) {
	x := centerX - len(value)*7/2
	drawText(screen, value, x+2, y+2, color.RGBA{255, 246, 184, 255})
	drawText(screen, value, x, y, c)
}

func drawText(screen *ebiten.Image, value string, x, y int, c color.Color) {
	op := &textv2.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y)-uiFace.Metrics().HAscent)
	op.ColorScale.ScaleWithColor(c)
	textv2.Draw(screen, value, uiFace, op)
}
