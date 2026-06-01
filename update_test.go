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
