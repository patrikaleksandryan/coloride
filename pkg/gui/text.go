package gui

import (
	"github.com/patrikaleksandryan/coloride/pkg/color"
	"github.com/veandco/go-sdl2/sdl"
)

func PrintChar(r rune, x, y int, color, bgColor color.Color) {
	mainFont.PrintChar(r, x, y, sdl.Color(color), sdl.Color(bgColor))
}

func Print(s string, x, y int, color, bgColor color.Color) {
	charW, _ := mainFont.Size()
	for _, r := range s {
		PrintChar(r, x, y, color, bgColor)
		x += charW
	}
}

func PrintAlign(s string, x, y, w int, color, bgColor color.Color, align int) {
	charW, _ := mainFont.Size()

	switch align {
	case AlignCenter:
		x += (w - charW*len(s)) / 2
	case AlignRight:
		x += w - charW*len(s)
	}

	Print(s, x, y, color, bgColor)
}

func PrintCentered(s string, x, y, w, h int, color, bgColor color.Color) {
	charW, charH := mainFont.Size()
	strW := charW * len(s)
	Print(s, x+(w-strW)/2, y+(h-charH)/2, color, bgColor)
}

func FontSize() (charW, charH int) {
	return mainFont.Size()
}
