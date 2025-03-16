package gui

import "github.com/veandco/go-sdl2/sdl"

func PrintChar(r rune, x, y int, color Color) {
	mainFont.PrintChar(r, x, y, sdl.Color(color))
}

func FontSize() (charW, charH int) {
	return mainFont.Size()
}
