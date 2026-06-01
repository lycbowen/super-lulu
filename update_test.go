package main

import (
	"image"
	"testing"
)

func TestResetPlayerClearsBossArenaLock(t *testing.T) {
	g := &Game{
		player: &Player{},
		level: Level{
			Spawn: image.Point{X: 80, Y: 385},
			Boss: &Boss{
				Active:      true,
				Aggro:       true,
				ArenaMinX:   1000,
				ArenaMaxX:   2000,
				ArenaLocked: true,
				State:       bossCharge,
				Timer:       30,
			},
		},
	}

	g.resetPlayer()

	if g.level.Boss.ArenaLocked {
		t.Fatal("expected boss arena lock to be cleared")
	}
	if g.level.Boss.Aggro {
		t.Fatal("expected boss aggro to be cleared")
	}
	if g.camera != 0 {
		t.Fatalf("expected camera reset to 0, got %v", g.camera)
	}
}

func TestResolveBossVerticalNeedsEnoughPlatformSupport(t *testing.T) {
	g := &Game{
		level: Level{
			Platforms: []Rect{{X: 0, Y: 110, W: 100, H: 20}},
		},
	}
	boss := &Boss{
		Rect: Rect{X: 95, Y: 91, W: 100, H: 20},
		VY:   1,
	}

	landed := g.resolveBossVertical(boss)

	if landed {
		t.Fatal("expected boss to keep falling when only touching the platform edge")
	}
	if boss.OnGround {
		t.Fatal("expected boss to be off ground without enough platform support")
	}
}

func TestResolveBossVerticalLandsWithEnoughPlatformSupport(t *testing.T) {
	g := &Game{
		level: Level{
			Platforms: []Rect{{X: 0, Y: 110, W: 100, H: 20}},
		},
	}
	boss := &Boss{
		Rect: Rect{X: 60, Y: 91, W: 100, H: 20},
		VY:   1,
	}

	landed := g.resolveBossVertical(boss)

	if !landed {
		t.Fatal("expected boss to land when enough of its body is supported")
	}
	if !boss.OnGround {
		t.Fatal("expected boss to be on ground after landing")
	}
}

func TestLevelSelectPaginationPageCount(t *testing.T) {
	g := newTestGameWithLevels(6)

	if got := g.levelSelectPageCount(); got != 2 {
		t.Fatalf("expected 2 pages, got %d", got)
	}
	if got := g.lastLevelOnPage(1); got != 5 {
		t.Fatalf("expected last level on page 2 to be 5, got %d", got)
	}
}

func TestLevelSelectDownMovesToNextPage(t *testing.T) {
	g := newTestGameWithLevels(6)
	g.setSelectedLevel(3)

	g.moveSelectedLevelDown()

	if g.selectedLevel != 4 {
		t.Fatalf("expected selected level 4, got %d", g.selectedLevel)
	}
	if g.levelSelectPage != 1 {
		t.Fatalf("expected page 1, got %d", g.levelSelectPage)
	}
}

func TestLevelSelectUpMovesToPreviousPage(t *testing.T) {
	g := newTestGameWithLevels(6)
	g.setSelectedLevel(4)

	g.moveSelectedLevelUp()

	if g.selectedLevel != 3 {
		t.Fatalf("expected selected level 3, got %d", g.selectedLevel)
	}
	if g.levelSelectPage != 0 {
		t.Fatalf("expected page 0, got %d", g.levelSelectPage)
	}
}

func TestMoveLevelSelectPagePreservesClosestRow(t *testing.T) {
	g := newTestGameWithLevels(6)
	g.setSelectedLevel(2)

	g.moveLevelSelectPage(1)

	if g.selectedLevel != 5 {
		t.Fatalf("expected selected level 5, got %d", g.selectedLevel)
	}
	if g.levelSelectPage != 1 {
		t.Fatalf("expected page 1, got %d", g.levelSelectPage)
	}
}

func TestSetSelectedLevelClampsToUnlockedLevel(t *testing.T) {
	g := newTestGameWithLevels(6)
	g.unlockedLevel = 0

	g.setSelectedLevel(4)

	if g.selectedLevel != 0 {
		t.Fatalf("expected selected level to clamp to 0, got %d", g.selectedLevel)
	}
	if g.levelSelectPage != 0 {
		t.Fatalf("expected page 0, got %d", g.levelSelectPage)
	}
}

func TestLevelSelectCannotMoveDownIntoLockedLevels(t *testing.T) {
	g := newTestGameWithLevels(6)
	g.unlockedLevel = 0
	g.setSelectedLevel(0)

	g.moveSelectedLevelDown()

	if g.selectedLevel != 0 {
		t.Fatalf("expected selected level to stay 0, got %d", g.selectedLevel)
	}
}

func TestLevelSelectCannotPageIntoLockedLevels(t *testing.T) {
	g := newTestGameWithLevels(6)
	g.unlockedLevel = 0
	g.setSelectedLevel(0)

	g.moveLevelSelectPage(1)

	if g.selectedLevel != 0 {
		t.Fatalf("expected selected level to stay 0, got %d", g.selectedLevel)
	}
	if g.levelSelectPage != 0 {
		t.Fatalf("expected page to stay 0, got %d", g.levelSelectPage)
	}
}

func newTestGameWithLevels(count int) *Game {
	levels := make([]Level, count)
	for i := range levels {
		levels[i] = Level{Name: "test"}
	}
	return &Game{
		player:        &Player{},
		levels:        levels,
		unlockedLevel: count - 1,
	}
}
