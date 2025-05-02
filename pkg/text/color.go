package text

import (
	"github.com/patrikaleksandryan/coloride/pkg/color"
)

const (
	colorCount = 9
)

type ColorInfo struct {
	Color           color.Color // Color is applied if overrideColor = true
	BgColor         color.Color // BgColor is applied if overrideBgColor = true
	overrideColor   bool
	overrideBgColor bool
}

var palette [colorCount]ColorInfo

func init() {
	palette = [colorCount]ColorInfo{
		{}, // Color #0 - don't change any colors
		{BgColor: color.MakeColor(170, 0, 0), overrideBgColor: true},
		{BgColor: color.MakeColor(0, 170, 0), overrideBgColor: true},
		{BgColor: color.MakeColor(0, 50, 170), overrideBgColor: true},
		{BgColor: color.MakeColor(170, 150, 0), overrideBgColor: true},
		{Color: color.White, BgColor: color.MakeColor(240, 0, 0), overrideColor: true, overrideBgColor: true},
		{Color: color.White, BgColor: color.MakeColor(0, 230, 0), overrideColor: true, overrideBgColor: true},
		{Color: color.White, BgColor: color.MakeColor(0, 200, 255), overrideColor: true, overrideBgColor: true},
		{Color: color.Black, BgColor: color.MakeColor(240, 230, 0), overrideColor: true, overrideBgColor: true},
	}
}
