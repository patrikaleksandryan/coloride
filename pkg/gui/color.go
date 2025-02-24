package gui

import "github.com/veandco/go-sdl2/sdl"

type Color sdl.Color

func MakeColor(r, g, b uint8) Color {
	return Color{R: r, G: g, B: b, A: 255}
}

func SetColor(color Color) {
	Renderer.SetDrawColor(color.R, color.G, color.B, color.A)
}
