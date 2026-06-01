package main

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed levels/*.json
var embeddedLevels embed.FS

// marshalLevelsJSON 统一用带缩进的格式导出关卡，避免生成压成一行、后续不好手改。
func marshalLevelsJSON(levels []Level) ([]byte, error) {
	data, err := json.MarshalIndent(levels, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

// newLevels 从带缩进的 JSON 文件读取关卡。以后改地图主要编辑 levels/levels.json，不需要改 Go 代码。
func newLevels() ([]Level, error) {
	data, err := embeddedLevels.ReadFile("levels/levels.json")
	if err != nil {
		return nil, fmt.Errorf("read levels json: %w", err)
	}

	var levels []Level
	if err := json.Unmarshal(data, &levels); err != nil {
		return nil, fmt.Errorf("parse levels json: %w", err)
	}
	if len(levels) == 0 {
		return nil, fmt.Errorf("levels json has no levels")
	}
	return levels, nil
}
