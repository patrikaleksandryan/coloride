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
		/* Color #0 */ {}, // don't change any colors
		/* Color #1 */ {BgColor: color.MakeColor(50, 25, 25), overrideBgColor: true},
		/* Color #2 */ {BgColor: color.MakeColor(25, 50, 25), overrideBgColor: true},
		/* Color #3 */ {BgColor: color.MakeColor(25, 25, 50), overrideBgColor: true},
		/* Color #4 */ {BgColor: color.MakeColor(50, 50, 25), overrideBgColor: true},
		/* Color #5 */ {Color: color.White, BgColor: color.MakeColor(200, 0, 0), overrideColor: true, overrideBgColor: true},
		/* Color #6 */ {Color: color.White, BgColor: color.MakeColor(0, 170, 0), overrideColor: true, overrideBgColor: true},
		/* Color #7 */ {Color: color.White, BgColor: color.MakeColor(0, 20, 200), overrideColor: true, overrideBgColor: true},
		/* Color #8 */ {Color: color.Black, BgColor: color.MakeColor(240, 230, 0), overrideColor: true, overrideBgColor: true},
	}
}
