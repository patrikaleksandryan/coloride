package gui

import "github.com/veandco/go-sdl2/sdl"

type Color sdl.Color

var (
	Transparent = Color{R: 0, G: 0, B: 0, A: 0}
	Black       = Color{R: 0, G: 0, B: 0, A: 255}
	White       = Color{R: 255, G: 255, B: 255, A: 255}
)

func MakeColor(r, g, b uint8) Color {
	return Color{R: r, G: g, B: b, A: 255}
}

func SetColor(color Color) {
	Renderer.SetDrawColor(color.R, color.G, color.B, color.A)
}
