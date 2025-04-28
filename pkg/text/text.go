package text

import (
	"bufio"
	"fmt"
	"os"

	"github.com/patrikaleksandryan/coloride/pkg/scanner"
)

type Text interface {
	HandleDelete()
	HandleBackspace()
	HandleEnter()
	HandleEscape()
	HandleHome(shift bool)
	HandleEnd(shift bool)
	HandleLeft(shift bool)
	HandleRight(shift bool)
	HandleUp(shift bool)
	HandleDown(shift bool)
	HandlePageUp(shift bool)
	HandlePageDown(shift bool)
	HandleChar(r rune)

	Reader() *Reader
	CurLine() *Line
	CurLineNum() int
	LineByNum(lineNum int) (line *Line, correctedNum int)
	SetCurLine(line *Line, lineNum int)
	TopLine() (*Line, int)
	CursorX() int
	SetCursorX(cursorX int)

	InSelection(lineNum, charNum int) bool
	ClearSelection()
	SetSelection(lineFrom, charFrom, lineTo, charTo int)
	StartMouseSelection()
	ContinueMouseSelection()

	Resize(w, h int)
	SetFontSize(charW, charH int)
	SetTabSize(tabSize int)
	TabSize() int
	ScrollValues() (scrollX, scrollY int)

	VisualToCursorX(l *Line, x int) int
	CursorXToVisual(l *Line, x int) int

	MergeLines(l *Line)
	DeleteLine(l *Line)

	LoadFromFile(fname string) error
	SaveToFile(fname string) error
	ColorizeSelection(color int)
}

// Selection represents a portion of text being selected
type Selection struct {
	LineFrom int // Including
	LineTo   int // Including
	CharFrom int // Including
	CharTo   int // Excluding
}

type TextImpl struct {
	cursorX    int // 0-based, x
	cursorMem  int // 0-based, x
	curLineNum int // 1-based, y
	curLine    *Line

	selected      bool
	selection     Selection
	oldCurLine    *Line // Value of curLine, saved in SelectionBefore
	oldCurLineNum int   // Value of curLineNum, saved in SelectionBefore
	oldCursorX    int   // Value of cursorX, saved in SelectionBefore

	first, last *Line // First and last line of document
	topLine     *Line // First line visible on the screen
	topLineNum  int   // 1-based

	w, h             int // Size of editor frame in pixels
	charW, charH     int // Size of character in pixels
	scrollX, scrollY int // Text scroll relative to frame in pixels, positive
	tabSize          int // Number of spaces in a tab

	reader *Reader
}

func (s Selection) isEmpty() bool {
	return s.LineFrom == s.LineTo && s.CharFrom == s.CharTo
}

func (s Selection) isInverted() bool {
	return s.LineFrom > s.LineTo || s.LineFrom == s.LineTo && s.CharFrom > s.CharTo
}

func NewText(w, h, charW, charH int) Text {
	line := NewLine()
	text := &TextImpl{
		first:      line,
		last:       line,
		curLine:    line,
		curLineNum: 1,
		topLine:    line,
		topLineNum: 1,
	}
	text.Resize(w, h)
	text.SetFontSize(charW, charH)
	text.SetTabSize(4)
	text.reader = NewReader(text)
	return text
}

func (t *TextImpl) LoadFromFile(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	s := scanner.NewScanner(buf)

	err = t.scanFile(s)
	if err != nil {
		return fmt.Errorf("scan file: %w", err)
	}
	return nil
}

func (t *TextImpl) scanFile(s *scanner.Scanner) error {
	s.Scan()
	for s.Sym != scanner.EOT {
		toAppend := make([]rune, 0)
		if s.Sym == scanner.String {
			toAppend = s.String
			s.Scan()
		}

		if s.Sym == scanner.ColorMarker {
			toAppend, t.curLine.spaces = splitTrailingWhitespace(toAppend)
			s.Scan()
			code := make([]rune, 0, 20)
			for s.Sym != scanner.EOT && s.Sym != scanner.NewLine {
				if s.Sym == scanner.String {
					code = append(code, s.String...)
				} else /* s.Sym == scanner.ColorMarker */ {
					code = append(code, '/', '/', '/')
				}
				s.Scan()
			}
			t.curLine.colorCode = code
		}

		for _, r := range toAppend {
			t.HandleChar(r)
		}

		t.curLine.ApplyColorCode()

		if s.Sym == scanner.NewLine {
			t.HandleEnter()
			s.Scan()
		}
	}

	t.MoveToBeginning()

	return nil
}

// splitTrailingWhitespace splits s into two parts, the second one being all trailing whitespaces from s.
func splitTrailingWhitespace(s []rune) ([]rune, []rune) {
	i := len(s)
	for i != 0 && s[i-1] <= ' ' {
		i--
	}
	return s[:i], s[i:]
}

func (t *TextImpl) SaveToFile(fname string) error {
	//!TODO
	_ = fname
	return nil
}

func (t *TextImpl) setDummyText() {
	//	s := `MODULE TestReader;
	//IMPORT Out, Texts;
	//VAR T: Texts.Text; R: Texts.Reader; S: Texts.Scanner;
	//  ch: CHAR;
	//BEGIN
	//	NEW(T); Texts.Open(T, 'Data/TEXT.DAT');
	//	Texts.OpenScanner(S, T, 1); Texts.Scan(S);
	//    Out.String(S.s); Out.Char(';'); Out.Ln;
	//    Texts.OpenReader(R, T, 0);
	//	Texts.Read(R, ch);
	//    WHILE ~R.eot DO
	//        Out.Int(ORD(R.eot), 5);
	//		Out.Int(ORD(ch), 5); Out.String('   ');
	//		Out.Char(ch); Out.Ln;
	//        Texts.Read(R, ch)
	//	END
	//END TestReader.`

	s := `package main

import (
	"fmt"
	"os"

	"github.com/patrikaleksandryan/coloride/pkg/editor"
	"github.com/patrikaleksandryan/coloride/pkg/gui"
)

const (
	windowWidth  = 1000
	windowHeight = 750
)

func run() error {
	err := gui.Init(windowWidth, windowHeight)
	if err != nil {
		return err
	}

	initInterface("Hello world", 412)

	err = /* gui.Run()
	if err != nil {
		return err
	}*/ fmt.Println("Hello")

	gui.Close()

	return nil
}

type User struct {
	Name 		/* this is a comment*/ string
	Age  int // Also this is a comment
}

func initInterface() {
	window := editor.NewWindow  ('x', 'y')
	gui.Append(window, ` + "`" + `Hello world
		another text here
		this is a text` + "`" + `)
	gui.SetFocus(window.Editor 	 	())
}`

	for _, r := range s {
		if r == 0xA {
			t.HandleEnter()
		} else {
			t.HandleChar(r)
		}
	}
	t.MoveToBeginning()

	r1 := t.first.next.next.runs
	r2 := &Run{length: 4, color: 7 - 4}
	r1.length -= 4
	r1.next = r2
}

func (t *TextImpl) Reader() *Reader {
	return t.reader
}

func (t *TextImpl) UpdateCursorMem() {
	t.cursorMem = t.CursorXToVisual(t.curLine, t.cursorX)
}

func (t *TextImpl) HandleEscape() {
	t.ClearSelection()
}

func (t *TextImpl) HandleDelete() {
	t.ClearSelection()
	if t.cursorX != len(t.curLine.chars) {
		t.curLine.DeleteChar(t.cursorX)
	} else if t.curLine.next != nil {
		dx := len(t.curLine.chars)
		t.HandleRight(false)
		t.MergeLines(t.curLine.prev)
		t.cursorX = dx
	}
	t.UpdateCursorMem()
}

func (t *TextImpl) HandleBackspace() {
	t.ClearSelection()
	if t.cursorX != 0 {
		t.curLine.DeleteChar(t.cursorX - 1)
		t.cursorX--
	} else if t.curLine.prev != nil {
		dx := len(t.curLine.prev.chars)
		t.MergeLines(t.curLine.prev)
		t.cursorX = dx
	}
	t.UpdateCursorMem()
}

// InSelection reports whether character charNum on line lineNum is currently being selected.
// Also returns true if charNum = -1 and selection spans until lineNum (i.e. new line character is being selected).
func (t *TextImpl) InSelection(lineNum, charNum int) bool {
	if !t.selected {
		return false
	}

	if lineNum < t.selection.LineFrom {
		return false
	}

	if lineNum > t.selection.LineTo {
		return false
	}

	if t.selection.LineFrom < lineNum && lineNum < t.selection.LineTo {
		return true
	}

	// The whole selection is within a single line
	if lineNum == t.selection.LineFrom && lineNum == t.selection.LineTo {
		return charNum >= t.selection.CharFrom && charNum < t.selection.CharTo
	}

	if lineNum == t.selection.LineFrom {
		return charNum >= t.selection.CharFrom
	}

	// lineNum == t.selection.LineTo
	return charNum < t.selection.CharTo
}

func (t *TextImpl) wasLeftSelectionEdge() bool {
	return t.oldCurLineNum == t.selection.LineFrom && t.oldCursorX == t.selection.CharFrom
}

func (t *TextImpl) IndentLength(line *Line) int {
	i := 0
	for i != len(line.chars) && line.chars[i] <= ' ' {
		i++
	}
	return i
}

func (t *TextImpl) ClearSelection() {
	t.selected = false
}

func (t *TextImpl) SetSelection(lineFrom, charFrom, lineTo, charTo int) {
	t.selected = true
	t.selection.LineFrom = lineFrom
	t.selection.LineTo = lineTo
	t.selection.CharFrom = charFrom
	t.selection.CharTo = charTo

	t.normalizeSelection()
}

func (t *TextImpl) normalizeSelection() {
	if t.selection.isEmpty() {
		t.ClearSelection()
	} else if t.selection.isInverted() {
		t.selection.CharFrom, t.selection.CharTo = t.selection.CharTo, t.selection.CharFrom
		t.selection.LineFrom, t.selection.LineTo = t.selection.LineTo, t.selection.LineFrom
	}
}

func (t *TextImpl) SelectionBefore() {
	t.oldCurLine = t.curLine
	t.oldCurLineNum = t.curLineNum
	t.oldCursorX = t.cursorX
}

func (t *TextImpl) SelectionAfter(shift bool) {
	if !shift {
		t.ClearSelection()
	} else if !t.selected { // Shift is pressed, there is no selection yet
		t.SetSelection(t.oldCurLineNum, t.oldCursorX, t.curLineNum, t.cursorX)
	} else if t.wasLeftSelectionEdge() { // Shift is pressed, selection exists
		t.SetSelection(t.curLineNum, t.cursorX, t.selection.LineTo, t.selection.CharTo)
	} else {
		t.SetSelection(t.selection.LineFrom, t.selection.CharFrom, t.curLineNum, t.cursorX)
	}
	t.MoveToCursor()
}

func (t *TextImpl) CursorTooLow() bool {
	return t.curLineNum*t.charH > t.scrollY+t.h
}

func (t *TextImpl) CursorTooHigh() bool {
	return (t.curLineNum-1)*t.charH < t.scrollY
}

func (t *TextImpl) MoveToCursor() {
	if t.CursorTooLow() {
		y := t.curLineNum * t.charH
		t.scrollY = y - t.h
	} else if t.CursorTooHigh() {
		y := (t.curLineNum - 1) * t.charH
		t.scrollY = y
	}
}

func (t *TextImpl) MoveToBeginning() {
	t.curLine = t.first
	t.curLineNum = 1
	t.cursorX = 0
	t.UpdateCursorMem()
	t.MoveToCursor()
}

func (t *TextImpl) StartMouseSelection() {
	t.ClearSelection()
	t.SelectionBefore()
}

func (t *TextImpl) ContinueMouseSelection() {
	t.SelectionAfter(true)
	t.SelectionBefore()
}

func (t *TextImpl) HandleEnter() {
	t.ClearSelection()

	t.SplitLine(t.curLine, t.cursorX)

	t.curLine = t.curLine.next
	t.curLineNum++
	t.cursorX = 0
	t.UpdateCursorMem()
}

func (t *TextImpl) HandleHome(shift bool) {
	t.SelectionBefore()
	ident := t.IndentLength(t.curLine)
	if t.cursorX > ident {
		t.cursorX = ident
	} else {
		t.cursorX = 0
	}
	t.UpdateCursorMem()
	t.SelectionAfter(shift)
}

func (t *TextImpl) HandleEnd(shift bool) {
	t.SelectionBefore()
	t.cursorX = len(t.curLine.chars)
	t.UpdateCursorMem()
	t.SelectionAfter(shift)
}

func (t *TextImpl) HandleLeft(shift bool) {
	t.SelectionBefore()
	if t.cursorX > 0 {
		t.cursorX--
	} else if t.curLine.prev != nil {
		t.curLineNum--
		t.curLine = t.curLine.prev
		t.cursorX = len(t.curLine.chars)
	}
	t.UpdateCursorMem()
	t.SelectionAfter(shift)
}

func (t *TextImpl) HandleRight(shift bool) {
	t.SelectionBefore()
	if t.cursorX < len(t.curLine.chars) {
		t.cursorX++
	} else if t.curLine.next != nil {
		t.curLineNum++
		t.curLine = t.curLine.next
		t.cursorX = 0
	}
	t.UpdateCursorMem()
	t.SelectionAfter(shift)
}

func (t *TextImpl) HandleUp(shift bool) {
	t.SelectionBefore()
	if t.curLine.prev != nil {
		t.curLineNum--
		t.curLine = t.curLine.prev

		t.cursorX = t.VisualToCursorX(t.curLine, t.cursorMem)
		if t.cursorX > len(t.curLine.chars) {
			t.cursorX = len(t.curLine.chars)
		}
	} else {
		t.cursorX = 0
	}
	t.SelectionAfter(shift)
}

func (t *TextImpl) HandleDown(shift bool) {
	t.SelectionBefore()
	if t.curLine.next != nil {
		t.curLineNum++
		t.curLine = t.curLine.next

		t.cursorX = t.VisualToCursorX(t.curLine, t.cursorMem)
		if t.cursorX > len(t.curLine.chars) {
			t.cursorX = len(t.curLine.chars)
		}
	} else {
		t.cursorX = len(t.curLine.chars)
	}
	t.SelectionAfter(shift)
}

// VisualToCursorX returns x recalculated as cursorX for the given line.
func (t *TextImpl) VisualToCursorX(l *Line, x int) int {
	tabSize := t.TabSize()
	m := l.Chars()
	oldVisualX := 0
	visualX := 0
	i := 0
	for i != len(m) && visualX < x {
		charCount := 1
		if m[i] == '\t' {
			charCount = tabSize - visualX%tabSize
		}
		oldVisualX = visualX
		visualX += charCount
		i++
	}

	//TODO use oldVisualX to approximate tab
	_ = oldVisualX

	return i
}

// CursorXToVisual returns x recalculated as visual X for the given line.
func (t *TextImpl) CursorXToVisual(l *Line, x int) int {
	tabSize := t.TabSize()
	m := l.Chars()
	if x > len(m) {
		x = len(m)
	}
	visualX := 0
	i := 0
	for i != x {
		charCount := 1
		if m[i] == '\t' {
			charCount = tabSize - visualX%tabSize
		}
		visualX += charCount
		i++
	}
	return visualX
}

func (t *TextImpl) HandlePageUp(shift bool) {
	t.SelectionBefore()

	lines := t.h / t.charH
	if lines == 0 {
		lines = 1
	}

	for t.curLine.prev != nil && lines != 0 {
		t.curLineNum--
		t.curLine = t.curLine.prev
		t.scrollY -= t.charH
		lines--
	}

	if t.scrollY < 0 {
		t.scrollY = 0
	}

	if lines == 0 {
		t.cursorX = t.VisualToCursorX(t.curLine, t.cursorMem)
		if t.cursorX > len(t.curLine.chars) {
			t.cursorX = len(t.curLine.chars)
		}
	} else { // Reached start of text
		t.cursorX = 0
	}

	t.SelectionAfter(shift)
}

func (t *TextImpl) HandlePageDown(shift bool) {
	t.SelectionBefore()

	lines := t.h / t.charH
	if lines == 0 {
		lines = 1
	}

	for t.curLine.next != nil && lines != 0 {
		t.curLineNum++
		t.curLine = t.curLine.next
		t.scrollY += t.charH
		lines--
	}

	if lines == 0 {
		t.cursorX = t.VisualToCursorX(t.curLine, t.cursorMem)
		if t.cursorX > len(t.curLine.chars) {
			t.cursorX = len(t.curLine.chars)
		}
	} else { // Reached end of text
		t.cursorX = len(t.curLine.chars)
	}

	t.SelectionAfter(shift)
}

func (t *TextImpl) HandleChar(r rune) {
	t.SelectionBefore()

	t.curLine.InsertChar(t.cursorX, r)

	t.cursorX++
	t.UpdateCursorMem()
	t.SelectionAfter(false)
}

func (t *TextImpl) CurLine() *Line {
	return t.curLine
}

func (t *TextImpl) CurLineNum() int {
	return t.curLineNum
}

// LineByNum returns a line based on its number, and the corrected line number.
// It returns the first line if lineNum <= 0.
// It returns the last line if lineNum is bigger than the last line.
func (t *TextImpl) LineByNum(lineNum int) (*Line, int) {
	l := t.first
	correctedNum := 1
	if lineNum > 0 {
		for l.next != nil && correctedNum != lineNum {
			l = l.next
			correctedNum++
		}
	}
	return l, correctedNum
}

func (t *TextImpl) SetCurLine(line *Line, lineNum int) {
	t.curLine = line
	t.curLineNum = lineNum
	t.cursorX = 0
	t.UpdateCursorMem()
}

func (t *TextImpl) TopLine() (*Line, int) {
	return t.topLine, t.topLineNum
}

func (t *TextImpl) CursorX() int {
	return t.cursorX
}

func (t *TextImpl) SetCursorX(cursorX int) {
	n := len(t.curLine.chars)
	if cursorX > n {
		cursorX = n
	}
	t.cursorX = cursorX
	t.UpdateCursorMem()
}

func (t *TextImpl) Resize(w, h int) {
	t.w, t.h = w, h
}

func (t *TextImpl) SetFontSize(charW, charH int) {
	t.charW, t.charH = charW, charH
}

func (t *TextImpl) SetTabSize(tabSize int) {
	t.tabSize = tabSize
}

func (t *TextImpl) TabSize() int {
	return t.tabSize
}

func (t *TextImpl) ScrollValues() (scrollX, scrollY int) {
	return t.scrollX, t.scrollY
}

// SplitLine splits the given line at position x, insereting the new line after the given line.
func (t *TextImpl) SplitLine(l *Line, x int) {
	l.Split(x)
	if l == t.last {
		t.last = l.next
	}
}

// MergeLines merges line l with the next line, appending its contents to l.
func (t *TextImpl) MergeLines(l *Line) {
	if l.next != nil {
		length := len(l.chars)
		l.chars = append(l.chars, l.next.chars...)

		if length == 0 {
			l.runs = l.next.runs
		} else {
			run := l.runs
			for run.next != nil {
				run = run.next
			}
			run.length--
			run.next = l.next.runs
		}

		l.NormalizeRuns()
		t.DeleteLine(l.next)
	}
}

// DeleteLine deletes line if it's not the only line.
func (t *TextImpl) DeleteLine(l *Line) {
	if l.prev != nil || l.next != nil {
		if l.prev != nil {
			l.prev.next = l.next
		} else {
			t.first = l.next
		}

		if l.next != nil {
			l.next.prev = l.prev
		} else {
			t.last = l.prev
		}

		if l == t.curLine {
			if l.prev == nil {
				t.curLine = t.first
				t.curLineNum = 1
			} else {
				t.curLine = l.prev
				t.curLineNum--
			}
			t.cursorX = 0
			t.UpdateCursorMem()
		}
	}
}

func (t *TextImpl) ColorizeSelection(color int) {
	if t.selected {
		sel := t.selection
		// First line of selection
		line, lineNum := t.LineByNum(t.selection.LineFrom)
		if sel.LineFrom == sel.LineTo { // One line selected
			line.Colorize(color, sel.CharFrom, sel.CharTo)
		} else { // Two or more lines
			line.Colorize(color, sel.CharFrom, len(line.chars)+1) // First line
			lineNum++
			line = line.next
			// Selection inner lines
			for lineNum != sel.LineTo {
				line.Colorize(color, 0, len(line.chars)+1)
				lineNum++
				line = line.next
			}
			// Last line of selection
			line.Colorize(color, 0, sel.CharTo)
		}
	}
}
