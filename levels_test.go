package main

import (
	"bytes"
	"testing"
)

func TestNewLevelsLoadsEmbeddedJSON(t *testing.T) {
	levels, err := newLevels()
	if err != nil {
		t.Fatal(err)
	}
	if len(levels) == 0 {
		t.Fatal("expected at least one level")
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
