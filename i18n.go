package main

type language int

const (
	languageEnglish language = iota
	languageChinese
)

func (g *Game) text(en, zh string) string {
	if g.language == languageChinese {
		return zh
	}
	return en
}

func (g *Game) levelSourceLabel() string {
	switch g.levelSource {
	case levelSourceExternal:
		return g.text("External custom levels", "外部自定义关卡")
	default:
		return g.text("Built-in levels", "内置关卡")
	}
}

func (g *Game) levelName(level Level) string {
	if g.language == languageChinese && level.NameZH != "" {
		return level.NameZH
	}
	return level.Name
}

func (g *Game) levelSubtitle(level Level) string {
	if g.language == languageChinese && level.SubtitleZH != "" {
		return level.SubtitleZH
	}
	return level.Subtitle
}

func (g *Game) toggleLanguage() {
	if g.language == languageChinese {
		g.language = languageEnglish
		return
	}
	g.language = languageChinese
}
