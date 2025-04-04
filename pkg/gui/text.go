package gui

import "github.com/veandco/go-sdl2/sdl"

func PrintChar(r rune, x, y int, color, bgColor Color) {
	mainFont.PrintChar(r, x, y, sdl.Color(color), sdl.Color(bgColor))
}

func FontSize() (charW, charH int) {
	return mainFont.Size()
}
