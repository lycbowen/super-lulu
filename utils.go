package main

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func clampInt(v, min, max int) int {
	return minInt(maxInt(v, min), max)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func signFloat(v float64) float64 {
	if v < 0 {
		return -1
	}
	if v > 0 {
		return 1
	}
	return 0
}

func (b Boss) WorldHitbox() Rect {
	if b.Hitbox.W <= 0 || b.Hitbox.H <= 0 {
		return Rect{
			X: b.Rect.X + b.Rect.W*0.18,
			Y: b.Rect.Y + b.Rect.H*0.24,
			W: b.Rect.W * 0.64,
			H: b.Rect.H * 0.70,
		}
	}
	return Rect{
		X: b.Rect.X + b.Hitbox.X,
		Y: b.Rect.Y + b.Hitbox.Y,
		W: b.Hitbox.W,
		H: b.Hitbox.H,
	}
}
