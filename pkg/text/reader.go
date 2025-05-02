package text

import (
	"github.com/patrikaleksandryan/coloride/pkg/color"
	"github.com/patrikaleksandryan/coloride/pkg/syntax"
)

type Reader struct {
	text *TextImpl

	curLine    *Line
	curLineNum int
	column     int // 0-based character number in curLine

	// Syntax highlighting
	symbolEnd    int // Column, where current symbol (lexem, token) ends
	symbolClass  int
	symbolColor  color.Color // Cache of symbolClass converted to Color
	nestingLevel int         // Can be derived from previous lines
}

type ColoredChar struct {
	Char    rune
	Color   color.Color
	BgColor color.Color
}

// TopLine resets the internal state of the reader and returns line number of the first line visible on the screen.
func (r *Reader) TopLine() int {
	r.curLine = r.text.topLine
	r.curLineNum = r.text.topLineNum
	r.column = 0
	r.symbolEnd = 0
	r.symbolClass = 0
	r.nestingLevel = 0
	return r.curLineNum
}

func (r *Reader) NextLine() int {
	r.curLine = r.curLine.next
	if r.curLine == nil {
		return -1
	}
	r.curLineNum++
	return r.curLineNum
}

func SymbolClassToColor(symbolClass int) color.Color {
	switch symbolClass {
	case syntax.CNone:
		return color.White
	case syntax.CComment:
		return color.MakeColor(120, 120, 120)
	case syntax.CIdent:
		return color.MakeColor(200, 200, 200)
	case syntax.CKeyword:
		return color.MakeColor(210, 150, 50)
	case syntax.CString:
		return color.MakeColor(70, 210, 50)
	case syntax.CNumber:
		return color.MakeColor(40, 235, 235)
	case syntax.CProcCall:
		return color.MakeColor(200, 180, 100)
	default:
		return color.White
	}
}

func (r *Reader) HighlightSyntax(char *ColoredChar) {
	if r.column >= r.symbolEnd {
		var length int
		r.symbolClass, length, r.nestingLevel = syntax.Scan(r.curLine.chars[r.column:], r.nestingLevel, r.symbolClass)
		r.symbolEnd = r.column + length
		r.symbolColor = SymbolClassToColor(r.symbolClass)
	}
	char.Color = r.symbolColor
}

func (r *Reader) Colorize(char *ColoredChar) {
	run, _ := r.curLine.FindRun(r.column)
	colorInfo := palette[run.color]
	if colorInfo.overrideColor {
		char.Color = colorInfo.Color
	}
	if colorInfo.overrideBgColor {
		char.BgColor = colorInfo.BgColor
	}
}

func (r *Reader) FirstChar() (char ColoredChar, ok bool) {
	r.column = 0 // Important to always do for ShouldPaintFullLine
	if len(r.curLine.chars) != 0 {
		r.symbolEnd = 0
		char.Char = r.curLine.chars[0]
		r.HighlightSyntax(&char)
		r.Colorize(&char)
		ok = true
	}
	return
}

func (r *Reader) NextChar() (char ColoredChar, ok bool) {
	r.column++
	if r.column != len(r.curLine.chars) {
		char.Char = r.curLine.chars[r.column]
		r.HighlightSyntax(&char)
		r.Colorize(&char)
		ok = true
	}
	return
}

// ShouldPaintFullLine reports whether the current line ends with a colorized new line character
// and places its colot into bgColor. Must be called right after NextChar returned false.
func (r *Reader) ShouldPaintFullLine(bgColor *color.Color) bool {
	run, _ := r.curLine.FindRun(r.column)
	//fmt.Println("lineNum=", r.curLineNum, "  X=", X, "  column=", r.column, "  runNil?=", run == nil)
	//fmt.Printf("  \"%s\"\n", string(r.curLine.chars))
	colorInfo := palette[run.color]
	if colorInfo.overrideBgColor {
		*bgColor = colorInfo.BgColor
		return true
	}
	return false
}

func NewReader(text *TextImpl) *Reader {
	r := &Reader{
		text: text,
	}
	return r
}
