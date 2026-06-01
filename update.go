package main

import (
	"math"
)

func (g *Game) Update() error {
	g.input.Update()
	if g.input.ToggleDebug {
		g.showDebug = !g.showDebug
	}
	switch g.mode {
	case modeMenu:
		g.updateMenu()
	case modeLevelSelect:
		g.updateLevelSelect()
	case modePlaying:
		g.updatePlayingMode()
	case modePaused:
		g.updatePaused()
	case modeLevelClear:
		g.updateLevelClear()
	case modeGameClear:
		g.updateGameClear()
	}
	return nil
}

func (g *Game) updateMenu() {
	if g.startDebugLevelByNumber() {
		return
	}
	if g.input.Start {
		g.startLevel(g.selectedLevel)
	}
	if g.input.LevelDown {
		g.mode = modeLevelSelect
	}
	if g.input.LevelUp {
		g.mode = modeLevelSelect
	}
}

func (g *Game) updateLevelSelect() {
	if g.startDebugLevelByNumber() {
		return
	}
	if g.input.Back {
		g.mode = modeMenu
		return
	}
	if g.input.LevelUp {
		g.selectedLevel = maxInt(0, g.selectedLevel-1)
	}
	if g.input.LevelDown {
		g.selectedLevel = minInt(g.unlockedLevel, g.selectedLevel+1)
	}
	if g.input.Start {
		g.startLevel(g.selectedLevel)
	}
}

func (g *Game) startDebugLevelByNumber() bool {
	if g.input.DebugLevel == 0 {
		return false
	}
	index := g.input.DebugLevel - 1
	if index < len(g.levels) {
		g.unlockedLevel = maxInt(g.unlockedLevel, index)
		g.startLevel(index)
		return true
	}
	return false
}

func (g *Game) updatePlayingMode() {
	if g.input.Pause {
		g.mode = modePaused
		return
	}
	if g.input.Restart {
		g.restartLevel()
		return
	}
	g.updatePlaying()
}

func (g *Game) updatePaused() {
	if g.input.Pause {
		g.mode = modePlaying
	}
	if g.input.Restart {
		g.restartLevel()
	}
	if g.input.Menu {
		g.mode = modeMenu
	}
}

func (g *Game) updateLevelClear() {
	if g.input.Start {
		next := g.currentLevel + 1
		if next >= len(g.levels) {
			g.mode = modeGameClear
			return
		}
		g.startLevel(next)
	}
	if g.input.Menu {
		g.mode = modeMenu
	}
}

func (g *Game) updateGameClear() {
	if g.input.Start {
		g.score = 0
		g.falls = 0
		g.startLevel(0)
	}
	if g.input.Menu || g.input.Back {
		g.mode = modeMenu
	}
}

func (g *Game) startLevel(index int) {
	g.currentLevel = clampInt(index, 0, len(g.levels)-1)
	g.selectedLevel = g.currentLevel
	g.loadLevel(g.currentLevel)
	g.mode = modePlaying
}

func (g *Game) restartLevel() {
	g.loadLevel(g.currentLevel)
	g.mode = modePlaying
}

func (g *Game) loadLevel(index int) {
	g.level = cloneLevel(g.levels[index])
	g.levelScore = 0
	g.hasWeapon = false
	g.player.Big = false
	g.invulnerable = 0
	g.resetPlayer()
}

// cloneLevel 将关卡模板复制成当前局的运行时关卡，避免收集物、敌人位置、Boss 血量等变化污染原始关卡数据。
func cloneLevel(level Level) Level {
	level.Platforms = append([]Rect(nil), level.Platforms...)
	level.Collect = append([]Collectible(nil), level.Collect...)
	level.PowerUps = append([]PowerUp(nil), level.PowerUps...)
	level.Oranges = append([]OrangePower(nil), level.Oranges...)
	level.Enemies = append([]Enemy(nil), level.Enemies...)
	if level.Boss != nil {
		boss := *level.Boss
		level.Boss = &boss
	}
	return level
}

func (g *Game) resetPlayer() {
	g.resetBossEncounterOnRespawn()
	g.player.X = float64(g.level.Spawn.X)
	g.player.Y = float64(g.level.Spawn.Y)
	g.player.VX = 0
	g.player.VY = 0
	g.player.OnGround = false
	if g.player.Facing == 0 {
		g.player.Facing = 1
	}
	g.camera = 0
	g.projectiles = nil
	g.shotCooldown = 0
}

// resetBossEncounterOnRespawn 清理 Boss 战的临时锁场状态，避免玩家回出生点后相机仍被 Boss 房间限制。
func (g *Game) resetBossEncounterOnRespawn() {
	boss := g.level.Boss
	if boss == nil || !boss.Active {
		return
	}
	boss.ArenaLocked = false
	boss.Aggro = false
	boss.State = bossPatrol
	boss.Timer = 0
	boss.ChargeDir = 0
	boss.VY = 0
	boss.AttackCooldown = maxInt(boss.AttackCooldown, 90)
}

func (g *Game) updatePlaying() {
	if g.invulnerable > 0 {
		g.invulnerable--
	}
	g.updateEnemies()
	g.updateBoss()
	g.updatePlayer()
	g.updateCollectibles()
	g.updatePowerUps()
	g.updateOranges()
	g.updateProjectiles()

	if g.player.Rect().Intersects(g.level.Goal) && !g.bossBlocksGoal() {
		g.unlockNextLevel()
		g.mode = modeLevelClear
		return
	}

	if g.player.Y > screenHeight+180 {
		g.falls++
		g.player.Big = false
		g.invulnerable = 0
		g.resetPlayer()
	}

	target := g.player.X - screenWidth*0.38
	target = math.Max(0, math.Min(target, g.level.Width-screenWidth))
	if boss, ok := g.lockedBossArena(); ok {
		target = clampCameraToArena(target, boss.ArenaMinX, boss.ArenaMaxX, g.level.Width)
	}
	g.camera += (target - g.camera) * 0.12
}

func (g *Game) bossBlocksGoal() bool {
	return g.level.Boss != nil && g.level.Boss.Active
}

// updateBoss 是 Boss 的主状态机：先处理冷却和仇恨，再按当前状态执行冲撞、泰山压顶或巡逻逻辑。
func (g *Game) updateBoss() {
	if g.level.Boss == nil || !g.level.Boss.Active {
		return
	}
	b := g.level.Boss
	if b.HitCooldown > 0 {
		b.HitCooldown--
	}
	if b.JumpCooldown > 0 {
		b.JumpCooldown--
	}
	if b.RespawnGrace > 0 {
		b.RespawnGrace--
	}
	if b.BaseY == 0 {
		b.BaseY = b.Rect.Y
	}
	g.updateBossAggro(b)
	g.lockBossArenaIfNeeded(b)
	if !b.Aggro && !g.rectOnScreen(b.Rect, 180) {
		g.patrolBoss(b)
		g.updateBossGravity(b)
		return
	}

	switch b.State {
	case bossWarning:
		b.Timer--
		g.updateBossGravity(b)
		if b.Timer <= 0 {
			b.State = bossCharge
			b.Timer = 90
			b.ChargeStartX = b.Rect.X
			b.ChargeDir = 1
			if g.player.Rect().X+g.player.Rect().W/2 < b.Rect.X+b.Rect.W/2 {
				b.ChargeDir = -1
			}
		}
	case bossCharge:
		speed := b.ChargeSpeed
		if speed <= 0 {
			speed = 8.5
		}
		distance := b.ChargeDistance
		if distance <= 0 {
			distance = 430
		}
		b.Rect.X += b.ChargeDir * speed
		b.Facing = int(b.ChargeDir)
		if b.Rect.X < b.MinX {
			b.Rect.X = b.MinX
		}
		if b.Rect.X+b.Rect.W > b.MaxX {
			b.Rect.X = b.MaxX - b.Rect.W
		}
		g.updateBossGravity(b)
		if math.Abs(b.Rect.X-b.ChargeStartX) >= distance || b.Rect.X <= b.MinX || b.Rect.X+b.Rect.W >= b.MaxX {
			b.State = bossRecover
			b.Timer = 45
		}
	case bossStompWarning:
		b.Timer--
		if b.Timer <= 0 {
			b.State = bossStompRise
			b.Timer = 18
			playerRect := g.player.Rect()
			b.SlamTargetX = clamp(playerRect.X+playerRect.W/2-b.Rect.W/2, 0, g.level.Width-b.Rect.W)
			b.SlamTargetY = math.Max(32, playerRect.Y-300)
		}
	case bossStompRise:
		b.Timer--
		oldX := b.Rect.X
		b.Rect.X += (b.SlamTargetX - b.Rect.X) * 0.28
		if math.Abs(b.Rect.X-oldX) > 0.1 {
			b.Facing = int(signFloat(b.Rect.X - oldX))
		}
		b.Rect.Y += (b.SlamTargetY - b.Rect.Y) * 0.28
		if b.Timer <= 0 {
			b.State = bossStompHang
			b.Rect.X = b.SlamTargetX
			b.Rect.Y = b.SlamTargetY
			b.VY = 0
			b.Timer = b.SlamHangTime
			if b.Timer <= 0 {
				b.Timer = 28
			}
		}
	case bossStompHang:
		b.Rect.X = b.SlamTargetX
		b.Rect.Y = b.SlamTargetY
		b.Timer--
		if b.Timer <= 0 {
			b.State = bossStompFall
		}
	case bossStompFall:
		b.VY += 0.72
		b.VY = math.Min(b.VY, 12.5)
		b.Rect.Y += b.VY
		if g.resolveBossSlamLanding(b) {
			b.State = bossRecover
			b.Timer = 60
			b.VY = 0
		}
	case bossRecover:
		b.Timer--
		g.updateBossGravity(b)
		if b.Timer <= 0 {
			b.State = bossPatrol
			b.AttackCooldown = 140 + g.rng.Intn(80)
		}
	default:
		if b.Aggro {
			g.moveBossAroundPlayer(b)
		} else {
			g.patrolBoss(b)
		}
		g.updateBossGravity(b)
		g.turnBossAwayFromLedge(b)
		if b.AttackCooldown > 0 {
			b.AttackCooldown--
			return
		}
		if g.rng.Intn(100) < 42 {
			b.State = bossWarning
			b.Timer = 60
		} else {
			b.State = bossStompWarning
			b.Timer = 60
			playerCenter := g.player.Rect().X + g.player.Rect().W/2
			b.SlamTargetX = clamp(playerCenter-b.Rect.W/2, 0, g.level.Width-b.Rect.W)
		}
	}
}

func (g *Game) updateBossAggro(b *Boss) {
	if b.Aggro {
		return
	}
	playerCenter := g.player.Rect().X + g.player.Rect().W/2
	bossCenter := b.Rect.X + b.Rect.W/2
	inArena := bossHasArena(b) && playerCenter >= b.ArenaMinX && playerCenter <= b.ArenaMaxX
	if inArena || math.Abs(playerCenter-bossCenter) < 560 || g.rectOnScreen(b.Rect, 80) {
		b.Aggro = true
	}
}

func bossHasArena(b *Boss) bool {
	return b != nil && b.ArenaMaxX > b.ArenaMinX
}

// lockBossArenaIfNeeded 在 Boss 被激活后关门锁场；Boss 死亡后 Active 为 false，锁场自然失效。
func (g *Game) lockBossArenaIfNeeded(b *Boss) {
	if !b.Aggro || !bossHasArena(b) {
		return
	}
	b.ArenaLocked = true
}

func (g *Game) lockedBossArena() (*Boss, bool) {
	b := g.level.Boss
	if b == nil || !b.Active || !b.ArenaLocked || !bossHasArena(b) {
		return nil, false
	}
	return b, true
}

func clampCameraToArena(target, minX, maxX, levelWidth float64) float64 {
	maxLevelCamera := math.Max(0, levelWidth-screenWidth)
	if maxX-minX <= screenWidth {
		centered := minX + (maxX-minX-screenWidth)/2
		return clamp(centered, 0, maxLevelCamera)
	}
	return clamp(target, minX, math.Min(maxX-screenWidth, maxLevelCamera))
}

func (g *Game) moveBossAroundPlayer(b *Boss) {
	playerCenter := g.player.Rect().X + g.player.Rect().W/2
	bossCenter := b.Rect.X + b.Rect.W/2
	distance := playerCenter - bossCenter
	absDistance := math.Abs(distance)

	speed := math.Abs(b.Speed)
	if speed < 0.75 {
		speed = 0.75
	}

	var dir float64
	switch {
	case absDistance > 520:
		dir = signFloat(distance)
		speed *= 3.4
	case absDistance > 190:
		dir = signFloat(distance)
		speed *= 2.0
	case absDistance < 95:
		dir = -signFloat(distance)
		speed *= 1.25
	default:
		dir = signFloat(b.Speed)
		speed *= 0.65
	}
	if dir == 0 {
		dir = signFloat(b.Speed)
	}
	if dir == 0 {
		dir = 1
	}
	if b.OnGround && !bossHasGroundAhead(g.level.Platforms, b, dir) {
		if b.JumpCooldown == 0 {
			b.VY = -11.4
			b.OnGround = false
			b.JumpCooldown = 55
		} else {
			return
		}
	}

	b.Rect.X += dir * speed
	b.Facing = int(dir)
	if b.Rect.X < 0 {
		b.Rect.X = 0
	}
	if b.Rect.X+b.Rect.W > g.level.Width {
		b.Rect.X = g.level.Width - b.Rect.W
	}
	b.Speed = math.Abs(b.Speed) * dir
}

func (g *Game) resolveBossVertical(b *Boss) bool {
	br := b.Rect
	b.OnGround = false
	for _, platform := range g.level.Platforms {
		if !br.Intersects(platform) {
			continue
		}
		if b.VY >= 0 && bossHasPlatformSupport(br, platform) && br.Y+br.H-b.VY <= platform.Y+8 {
			b.Rect.Y += platform.Y - (br.Y + br.H)
			b.OnGround = true
			return true
		}
	}
	if b.Rect.Y > screenHeight+260 {
		g.respawnBossOnNearbyGround(b)
		return true
	}
	return false
}

func (g *Game) resolveBossSlamLanding(b *Boss) bool {
	br := b.Rect
	for _, platform := range g.level.Platforms {
		if !br.Intersects(platform) {
			continue
		}
		nearTarget := math.Abs((br.X+br.W/2)-(b.SlamTargetX+b.Rect.W/2)) <= b.Rect.W*0.75
		if nearTarget && b.VY >= 0 && bossHasPlatformSupport(br, platform) && br.Y+br.H-b.VY <= platform.Y+8 {
			b.Rect.Y += platform.Y - (br.Y + br.H)
			b.OnGround = true
			return true
		}
	}
	if b.Rect.Y > screenHeight+260 {
		g.respawnBossOnNearbyGround(b)
		return true
	}
	return false
}

// bossHasPlatformSupport 要求 Boss 和平台有足够的水平重叠，避免只蹭到平台边缘时也被判定为站住。
func bossHasPlatformSupport(bossRect, platform Rect) bool {
	overlap := math.Min(bossRect.X+bossRect.W, platform.X+platform.W) - math.Max(bossRect.X, platform.X)
	minSupport := bossRect.W * 0.35
	return overlap >= minSupport
}

func (g *Game) updateBossGravity(b *Boss) {
	b.VY += gravity
	b.VY = math.Min(b.VY, 14)
	b.Rect.Y += b.VY
	if g.resolveBossVertical(b) {
		b.VY = 0
	}
}

func (g *Game) respawnBossOnNearbyGround(b *Boss) {
	spots := g.bossRespawnSpots(b)
	if len(spots) == 0 {
		b.Rect.Y = b.BaseY
		b.VY = 0
		b.OnGround = true
		b.RespawnGrace = 90
		return
	}
	spot := nearestRespawnSpot(spots, b.Rect.X+b.Rect.W/2)
	b.Rect.X = spot.X
	b.Rect.Y = spot.Y
	b.VY = 0
	b.OnGround = true
	b.State = bossRecover
	b.Timer = 50
	b.AttackCooldown = 150
	b.RespawnGrace = 90
}

func (g *Game) bossRespawnSpots(b *Boss) []Rect {
	center := b.Rect.X + b.Rect.W/2
	var spots []Rect
	for _, platform := range g.level.Platforms {
		if platform.W < b.Rect.W+80 || platform.H < 20 {
			continue
		}
		minX := platform.X + 36
		maxX := platform.X + platform.W - b.Rect.W - 36
		if maxX < minX {
			continue
		}
		x := clamp(center-b.Rect.W/2, minX, maxX)
		distance := math.Abs((x + b.Rect.W/2) - center)
		if distance > 760 {
			continue
		}
		spots = append(spots, Rect{X: x, Y: platform.Y - b.Rect.H})
	}
	return spots
}

func nearestRespawnSpot(spots []Rect, centerX float64) Rect {
	best := spots[0]
	bestDistance := math.Abs(best.X + best.W/2 - centerX)
	for _, spot := range spots[1:] {
		distance := math.Abs(spot.X + spot.W/2 - centerX)
		if distance < bestDistance {
			best = spot
			bestDistance = distance
		}
	}
	return best
}

func (g *Game) turnBossAwayFromLedge(b *Boss) {
	if !b.OnGround || b.RespawnGrace > 0 || bossHasGroundAhead(g.level.Platforms, b, signFloat(b.Speed)) {
		return
	}
	b.Speed *= -1
	b.Facing = int(signFloat(b.Speed))
}

// bossHasGroundAhead 用 Boss 脚前方的探测点判断是否有平台，防止巡逻或追击时直接走下悬崖。
func bossHasGroundAhead(platforms []Rect, b *Boss, dir float64) bool {
	if dir == 0 {
		return true
	}
	frontX := b.Rect.X + b.Rect.W + 22
	if dir < 0 {
		frontX = b.Rect.X - 22
	}
	footY := b.Rect.Y + b.Rect.H + 10
	for _, platform := range platforms {
		if footY < platform.Y || footY > platform.Y+platform.H+18 {
			continue
		}
		if frontX >= platform.X && frontX <= platform.X+platform.W {
			return true
		}
	}
	return false
}

func (g *Game) patrolBoss(b *Boss) {
	b.Rect.X += b.Speed
	b.Facing = int(signFloat(b.Speed))
	if b.Rect.X < b.MinX || b.Rect.X+b.Rect.W > b.MaxX {
		b.Speed *= -1
		b.Facing = int(signFloat(b.Speed))
		b.Rect.X += b.Speed
	}
}

func (g *Game) unlockNextLevel() {
	if g.currentLevel+1 > g.unlockedLevel {
		g.unlockedLevel = minInt(g.currentLevel+1, len(g.levels)-1)
	}
}

func (g *Game) updateEnemies() {
	for i := range g.level.Enemies {
		e := &g.level.Enemies[i]
		e.Rect.X += e.Speed
		if e.Rect.X < e.MinX || e.Rect.X+e.Rect.W > e.MaxX {
			e.Speed *= -1
			e.Rect.X += e.Speed
		}
	}
}

func (g *Game) updatePlayer() {
	p := g.player
	prevRect := p.Rect()
	accel := moveAccel
	if !p.OnGround {
		accel = airAccel
	}

	if g.input.MoveLeft {
		p.VX -= accel
		p.Facing = -1
	}
	if g.input.MoveRight {
		p.VX += accel
		p.Facing = 1
	}
	if !g.input.MoveLeft && !g.input.MoveRight {
		if p.OnGround {
			p.VX *= friction
		} else {
			p.VX *= 0.98
		}
	}

	p.VX = clamp(p.VX, -maxRunSpeed, maxRunSpeed)
	if p.OnGround && g.input.Jump {
		p.VY = jumpVelocity
		p.OnGround = false
	}
	if g.input.Shoot {
		g.shootIceCream()
	}

	p.VY += gravity
	p.VY = math.Min(p.VY, 14)

	p.X += p.VX
	g.resolveHorizontal()
	p.Y += p.VY
	g.resolveVertical()

	p.X = clamp(p.X, 0, g.level.Width-80)
	g.clampPlayerToLockedArena()

	for i := 0; i < len(g.level.Enemies); i++ {
		e := g.level.Enemies[i]
		if p.Rect().Intersects(e.Rect) {
			landedOnEnemy := p.VY >= 0 && prevRect.Y+prevRect.H <= e.Rect.Y+12
			if landedOnEnemy {
				g.level.Enemies = append(g.level.Enemies[:i], g.level.Enemies[i+1:]...)
				g.score += 200
				g.levelScore += 200
				p.VY = jumpVelocity * 0.48
				p.OnGround = false
				i--
				continue
			}
			g.hurtPlayer()
			return
		}
	}
	if g.level.Boss != nil && g.level.Boss.Active && p.Rect().Intersects(g.level.Boss.WorldHitbox()) {
		bossHitbox := g.level.Boss.WorldHitbox()
		landedOnBoss := p.VY >= 0 && prevRect.Y+prevRect.H <= bossHitbox.Y+14
		if landedOnBoss {
			g.damageBoss(1)
			p.VY = jumpVelocity * 0.55
			p.OnGround = false
			return
		}
		g.hurtPlayer()
		return
	}
}

// clampPlayerToLockedArena 用玩家碰撞框做边界限制，保证变大后的 Lulu 也不会穿过 Boss 房间门。
func (g *Game) clampPlayerToLockedArena() {
	boss, ok := g.lockedBossArena()
	if !ok {
		return
	}
	p := g.player
	pr := p.Rect()
	if pr.X < boss.ArenaMinX {
		p.X += boss.ArenaMinX - pr.X
		p.VX = 0
	}
	pr = p.Rect()
	if pr.X+pr.W > boss.ArenaMaxX {
		p.X += boss.ArenaMaxX - (pr.X + pr.W)
		p.VX = 0
	}
}

func (g *Game) shootIceCream() {
	if !g.hasWeapon || g.shotCooldown > 0 {
		return
	}
	dir := float64(g.player.Facing)
	if dir == 0 {
		dir = 1
	}
	x := g.player.X + 22
	if dir < 0 {
		x = g.player.X + 4
	}
	g.projectiles = append(g.projectiles, Projectile{
		Rect:        Rect{X: x, Y: g.player.Y + 20, W: 30, H: 30},
		VX:          dir * 7.2,
		StartX:      x,
		MaxDistance: 560,
		Active:      true,
	})
	g.shotCooldown = 18
}

// resolveHorizontal 只解决玩家横向进入平台的问题；纵向碰撞分开处理可以减少斜角卡住的情况。
func (g *Game) resolveHorizontal() {
	p := g.player
	pr := p.Rect()
	for _, platform := range g.level.Platforms {
		if !pr.Intersects(platform) {
			continue
		}
		if p.VX > 0 {
			p.X += platform.X - (pr.X + pr.W)
		} else if p.VX < 0 {
			p.X += platform.X + platform.W - pr.X
		}
		p.VX = 0
		pr = p.Rect()
	}
}

// resolveVertical 处理玩家上下方向的平台碰撞，并在落到平台上时设置 OnGround。
func (g *Game) resolveVertical() {
	p := g.player
	p.OnGround = false
	pr := p.Rect()
	for _, platform := range g.level.Platforms {
		if !pr.Intersects(platform) {
			continue
		}
		if p.VY > 0 {
			p.Y += platform.Y - (pr.Y + pr.H)
			p.OnGround = true
		} else if p.VY < 0 {
			p.Y += platform.Y + platform.H - pr.Y
		}
		p.VY = 0
		pr = p.Rect()
	}
}

func (g *Game) updateCollectibles() {
	pr := g.player.Rect()
	for i := range g.level.Collect {
		c := &g.level.Collect[i]
		if !c.Collected && pr.Intersects(c.Rect) {
			c.Collected = true
			g.score += 100
			g.levelScore += 100
		}
	}
}

func (g *Game) updatePowerUps() {
	pr := g.player.Rect()
	for i := range g.level.PowerUps {
		p := &g.level.PowerUps[i]
		if !p.Collected && pr.Intersects(p.Rect) {
			p.Collected = true
			g.hasWeapon = true
			g.score += 250
			g.levelScore += 250
		}
	}
}

func (g *Game) updateOranges() {
	pr := g.player.Rect()
	for i := range g.level.Oranges {
		o := &g.level.Oranges[i]
		if !o.Collected && pr.Intersects(o.Rect) {
			o.Collected = true
			g.player.Big = true
			g.invulnerable = 0
			g.score += 250
			g.levelScore += 250
		}
	}
}

func (g *Game) hurtPlayer() {
	if g.invulnerable > 0 {
		return
	}
	if g.player.Big {
		g.player.Big = false
		g.invulnerable = invulnerableFrames
		return
	}
	g.falls++
	g.resetPlayer()
}

func (g *Game) updateProjectiles() {
	if g.shotCooldown > 0 {
		g.shotCooldown--
	}
	for i := 0; i < len(g.projectiles); i++ {
		p := &g.projectiles[i]
		if !p.Active {
			continue
		}
		p.Rect.X += p.VX
		if math.Abs(p.Rect.X-p.StartX) > p.MaxDistance || !g.rectOnScreen(p.Rect, 8) {
			p.Active = false
			continue
		}
		for j := 0; j < len(g.level.Enemies); j++ {
			if p.Rect.Intersects(g.level.Enemies[j].Rect) {
				g.level.Enemies = append(g.level.Enemies[:j], g.level.Enemies[j+1:]...)
				g.score += 150
				g.levelScore += 150
				p.Active = false
				break
			}
		}
		if p.Active && g.level.Boss != nil && g.level.Boss.Active && p.Rect.Intersects(g.level.Boss.WorldHitbox()) {
			g.damageBoss(1)
			p.Active = false
		}
	}
	active := g.projectiles[:0]
	for _, p := range g.projectiles {
		if p.Active {
			active = append(active, p)
		}
	}
	g.projectiles = active
}

func (g *Game) rectOnScreen(r Rect, padding float64) bool {
	return r.X+r.W >= g.camera-padding && r.X <= g.camera+screenWidth+padding
}

func (g *Game) damageBoss(amount int) {
	if g.level.Boss == nil || !g.level.Boss.Active {
		return
	}
	if g.level.Boss.HitCooldown > 0 {
		return
	}
	g.level.Boss.HP -= amount
	g.level.Boss.HitCooldown = 24
	if g.level.Boss.HP <= 0 {
		g.level.Boss.Active = false
		g.score += 1000
		g.levelScore += 1000
	}
}
