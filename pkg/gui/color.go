package gui

import (
	"github.com/patrikaleksandryan/coloride/pkg/color"
)

func SetColor(color color.Color) {
	Renderer.SetDrawColor(color.R, color.G, color.B, color.A)
}

func SetRGB(r, g, b uint8) {
	Renderer.SetDrawColor(r, g, b, 255)
}

func SetRGBA(r, g, b, a uint8) {
	Renderer.SetDrawColor(r, g, b, a)
}
