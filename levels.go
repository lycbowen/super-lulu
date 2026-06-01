package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
)

const levelsJSONPath = "levels/levels.json"

//go:embed levels/*.json
var embeddedLevels embed.FS

type levelSource int

const (
	levelSourceEmbedded levelSource = iota
	levelSourceExternal
)

type LevelPack struct {
	Levels []Level
	Source levelSource
	Path   string
}

func (p LevelPack) SourceLabel() string {
	switch p.Source {
	case levelSourceExternal:
		return "External custom levels"
	default:
		return "Built-in levels"
	}
}

// marshalLevelsJSON 统一用带缩进的格式导出关卡，避免生成压成一行、后续不好手改。
func marshalLevelsJSON(levels []Level) ([]byte, error) {
	data, err := json.MarshalIndent(levels, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

// newLevelPack 优先读取外部 levels/levels.json，缺失时回落到内置关卡，方便发布版和自定义地图共存。
func newLevelPack() (LevelPack, error) {
	if data, err := os.ReadFile(levelsJSONPath); err == nil {
		levels, err := parseLevelsJSON(data)
		if err != nil {
			return LevelPack{}, fmt.Errorf("parse external levels %s: %w", levelsJSONPath, err)
		}
		return LevelPack{Levels: levels, Source: levelSourceExternal, Path: levelsJSONPath}, nil
	} else if !os.IsNotExist(err) {
		return LevelPack{}, fmt.Errorf("read external levels %s: %w", levelsJSONPath, err)
	}

	data, err := embeddedLevels.ReadFile(levelsJSONPath)
	if err != nil {
		return LevelPack{}, fmt.Errorf("read embedded levels json: %w", err)
	}
	levels, err := parseLevelsJSON(data)
	if err != nil {
		return LevelPack{}, fmt.Errorf("parse embedded levels: %w", err)
	}
	return LevelPack{Levels: levels, Source: levelSourceEmbedded, Path: levelsJSONPath}, nil
}

func parseLevelsJSON(data []byte) ([]Level, error) {
	var levels []Level
	if err := json.Unmarshal(data, &levels); err != nil {
		return nil, fmt.Errorf("parse levels json: %w", err)
	}
	if len(levels) == 0 {
		return nil, fmt.Errorf("levels json has no levels")
	}
	return levels, nil
}
