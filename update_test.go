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
