package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestNewLevelPackFallsBackToEmbeddedJSON(t *testing.T) {
	t.Chdir(t.TempDir())

	levelPack, err := newLevelPack()
	if err != nil {
		t.Fatal(err)
	}
	if levelPack.Source != levelSourceEmbedded {
		t.Fatalf("expected embedded levels, got %v", levelPack.Source)
	}
	if len(levelPack.Levels) == 0 {
		t.Fatal("expected at least one level")
	}
}

func TestNewLevelPackPrefersExternalJSON(t *testing.T) {
	t.Chdir(t.TempDir())

	if err := os.MkdirAll("levels", 0755); err != nil {
		t.Fatal(err)
	}
	data, err := marshalLevelsJSON([]Level{{Name: "custom", Width: 1234}})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.FromSlash(levelsJSONPath), data, 0644); err != nil {
		t.Fatal(err)
	}

	levelPack, err := newLevelPack()
	if err != nil {
		t.Fatal(err)
	}
	if levelPack.Source != levelSourceExternal {
		t.Fatalf("expected external levels, got %v", levelPack.Source)
	}
	if got := levelPack.Levels[0].Name; got != "custom" {
		t.Fatalf("expected custom level, got %q", got)
	}
}

func TestMarshalLevelsJSONUsesIndent(t *testing.T) {
	data, err := marshalLevelsJSON([]Level{{Name: "test"}})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte("\n  {")) {
		t.Fatalf("expected indented JSON, got %q", data)
	}
}
