package text

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/patrikaleksandryan/coloride/pkg/colorcode"
	"github.com/patrikaleksandryan/coloride/pkg/scanner"
	"github.com/veandco/go-sdl2/sdl"
)

// Color legend
// AI generated code which passed human check		///g
// AI code whish was not checked		///b

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
	HandleCut()
	HandleCopy()
	HandlePaste()
	HandleSelectAll()

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
	SetUpdaters(editedUpdater EditedUpdater, posUpdater PosUpdater)
	SetFontSize(charW, charH int)
	SetTabSize(tabSize int)
	TabSize() int
	ScrollValues() (scrollX, scrollY int)
	ScrollDelta(dy int)

	VisualToCursorX(l *Line, x int) int
	CursorXToVisual(l *Line, x int) int

	MergeLines(l *Line)
	DeleteLine(l *Line)

	Clear()
	LoadFromFile(fname string) error
	SaveToFile(fname string) error
	ColorizeSelection(color int)
}

type EditedUpdater interface {
	UpdateEdited(edited bool)
}

type PosUpdater interface {
	UpdatePos(line, col int)
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
	lineCount   int

	w, h             int // Size of editor frame in pixels
	charW, charH     int // Size of character in pixels
	scrollX, scrollY int // Text scroll relative to frame in pixels, positive
	tabSize          int // Number of spaces in a tab

	reader        *Reader
	edited        bool // If file was edited after it was opened
	editedUpdater EditedUpdater
	posUpdater    PosUpdater
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
		lineCount:  1,
	}
	text.Resize(w, h)
	text.SetFontSize(charW, charH)
	text.SetTabSize(4)
	text.reader = NewReader(text)
	return text
}

func (t *TextImpl) SetUpdaters(editedUpdater EditedUpdater, posUpdater PosUpdater) {
	t.editedUpdater = editedUpdater
	t.posUpdater = posUpdater
}

func (t *TextImpl) UpdatePos() {
	if t.posUpdater != nil {
		t.posUpdater.UpdatePos(t.curLineNum, t.cursorX+1)
	}
}

func (t *TextImpl) setEdited(edited bool) {
	if t.edited != edited {
		t.edited = edited
		if t.editedUpdater != nil {
			t.editedUpdater.UpdateEdited(t.edited)
		}
	}
}

func (t *TextImpl) Clear() {
	t.cursorX = 0
	t.cursorMem = 0
	t.curLineNum = 1
	t.topLineNum = 1
	t.selected = false
	t.oldCurLine = nil
	t.oldCurLineNum = 1
	t.oldCursorX = 0
	t.scrollX = 0
	t.scrollY = 0
	t.first = NewLine()
	t.last = t.first
	t.topLine = t.first
	t.curLine = t.first
	t.lineCount = 1
	t.setEdited(false)
}

func (t *TextImpl) LoadFromFile(fname string) error {
	t.Clear()
	f, err := os.Open(fname)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	s := scanner.NewScanner(buf)

	err = t.load(s)
	if err != nil {
		return fmt.Errorf("load file: %w", err)
	}
	t.setEdited(false)
	return nil
}

func (t *TextImpl) load(s *scanner.Scanner) error {
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

		// Apply color code only after HandleEnter
		line := t.curLine

		if s.Sym == scanner.NewLine {
			t.curLine.NewLineType = s.NewLineType
			t.HandleEnter()
			s.Scan()
		}

		line.ApplyColorCode()
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
	f, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	buf := bufio.NewWriter(f)

	err = t.write(buf)
	if err != nil {
		f.Close()
		return fmt.Errorf("write file: %w", err)
	}

	err = buf.Flush()
	if err != nil {
		f.Close()
		return fmt.Errorf("flush buffer: %w", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("close file: %w", err)
	}
	t.setEdited(false)
	return nil
}

func (t *TextImpl) write(buf *bufio.Writer) error {
	line := t.first
	for line != nil {
		err := t.writeLine(buf, line)
		if err != nil {
			return err
		}
		line = line.next
	}
	return nil
}

func (t *TextImpl) writeLine(buf *bufio.Writer, line *Line) error {
	_, err := buf.WriteString(string(line.chars))
	if err != nil {
		return err
	}

	if line.IsColorized() {
		if len(line.spaces) != 0 {
			_, err = buf.WriteString(string(line.spaces))
		} else {
			_, err = buf.WriteString("\t\t")
		}
		if err != nil {
			return err
		}
		t.openColorComment(buf)
		t.writeColorCode(buf, line.runs)
		t.closeColorComment(buf)
	}

	err = t.writeNewLine(buf, line.NewLineType)
	return err
}

func (t *TextImpl) openColorComment(buf *bufio.Writer) error {
	_, err := buf.WriteString("
	return err
}

func (t *TextImpl) closeColorComment(_ *bufio.Writer) error {
	return nil
}

// writeColorCode function generated by DeepSeek (passed human check)		///g
func (t *TextImpl) writeColorCode(buf *bufio.Writer, runs *Run) error {		///19g 14B g
	run := runs		///g
	first := true		///g
		///g
	//     !A || B    ===    A -> B		///g
	// run.next != nil || run.color != 0		///g
	// (run.next == nil) -> (run.color != 0)		///g
		///g
	// Iterate all runs, except the last one in case the last one is standard (color = 0)		///g
	//for run != nil && !(run.next == nil && run.color == 0) {		///g
	for run != nil && (run.next != nil || run.color != 0) {		///g
		// Output space, but not the first time		///g
		if !first {		///g
			_, err := buf.WriteString(" ")		///g
			if err != nil {		///g
				return err		///g
			}		///g
		}		///g
		first = false		///g
		// Output code for the run		///g
		if run.color == 0 { // Number. Represents a gap. The last one is not written		///g
			_, err := buf.WriteString(strconv.Itoa(run.length))		///g
			if err != nil {		///g
				return err		///g
			}		///g
		} else if run.next == nil { // Letter, because it is the last one		///g
			_, err := buf.WriteRune(colorcode.ToLetter(run.color))		///g
			if err != nil {		///g
				return err		///g
			}		///g
		} else { // Number + Letter		///g
			_, err := buf.WriteString(strconv.Itoa(run.length))		///g
			if err != nil {		///g
				return err		///g
			}		///g
			_, err = buf.WriteRune(colorcode.ToLetter(run.color))		///g
			if err != nil {		///g
				return err		///g
			}		///g
		}		///g
		run = run.next		///g
	}		///g
	return nil		///g
}		///g

func newLineTypeToString(newLineType int) string {
	switch newLineType {
	case scanner.LF:
		return "\n"
	case scanner.CRLF:
		return "\r\n"
	case scanner.CR:
		return "\r"
	default:
		panic("impossible: newLineType")
	}
}

func (t *TextImpl) writeNewLine(buf *bufio.Writer, newLineType int) error {
	_, err := buf.WriteString(newLineTypeToString(newLineType))
	return err
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

// HandleDelete function generated by DeepSeek (passed human check)		///g
func (t *TextImpl) HandleDelete() {		///19g 12B g
	if t.selected {		///g
		t.DeleteSelectedText()		///g
	} else {		///g
		t.ClearSelection()		///g
		if t.cursorX != len(t.curLine.chars) {		///g
			t.curLine.DeleteChar(t.cursorX)		///g
			t.setEdited(true)		///g
		} else if t.curLine.next != nil {		///g
			dx := len(t.curLine.chars)		///g
			t.HandleRight(false)		///g
			t.MergeLines(t.curLine.prev)		///g
			t.cursorX = dx		///g
			t.setEdited(true)		///g
		}		///g
		t.UpdateCursorMem()		///g
	}		///g
}		///g

func (t *TextImpl) HandleBackspace() {
	if t.selected {
		t.DeleteSelectedText()
	} else {
		t.ClearSelection()
		if t.cursorX != 0 {
			t.curLine.DeleteChar(t.cursorX - 1)
			t.cursorX--
			t.setEdited(true)
		} else if t.curLine.prev != nil {
			dx := len(t.curLine.prev.chars)
			t.MergeLines(t.curLine.prev)
			t.cursorX = dx
			t.setEdited(true)
		}
		t.UpdateCursorMem()
	}
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

func (t *TextImpl) DeleteSelectedText() {
	if t.selected {
		sel := t.selection
		// First line of selection
		line, lineNum := t.LineByNum(t.selection.LineFrom)
		if sel.LineFrom == sel.LineTo { // One line selected
			line.DeleteRange(sel.CharFrom, sel.CharTo)
			t.SetCursorX(t.selection.CharFrom)
		} else { // Two or more lines
			next := line.next
			if sel.CharFrom == 0 {
				t.DeleteLine(line) // Delete the first line of selection entierly
			} else {
				line.DeleteRange(sel.CharFrom, len(line.chars)+1) // Delete second part of the first line
				t.SetCurLine(line, lineNum)
				t.SetCursorX(t.selection.CharFrom)
			}
			lineNum++
			line = next
			// Selection inner lines
			for lineNum != sel.LineTo {
				next := line.next
				t.DeleteLine(line)
				lineNum++
				line = next
			}
			// Last line of selection
			line.DeleteRange(0, sel.CharTo)
			if sel.CharFrom == 0 {
				t.SetCurLine(line, t.selection.LineFrom)
				t.SetCursorX(0)
			} else {
				t.MergeLines(line.prev)
			}
		}
		t.selected = false
		t.setEdited(true)
	}
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
	t.UpdatePos()
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

// ScrollDelta is generated by ChatGPT		///b
func (t *TextImpl) ScrollDelta(dy int) {		///19b 11B b
	t.scrollY += dy		///b
		///b
	max := t.lineCount*t.charH - t.h		///b
	if t.scrollY > max {		///b
		t.scrollY = max		///b
	}		///b
	if t.scrollY < 0 {		///b
		t.scrollY = 0		///b
	}		///b
}		///b

func (t *TextImpl) MoveToBeginning() {
	t.curLine = t.first
	t.curLineNum = 1
	t.cursorX = 0
	t.UpdateCursorMem()
	t.MoveToCursor()
	t.UpdatePos()
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
	t.MoveToCursor()
	t.setEdited(true)
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

// VisualToCursorX returns x, recalculated as cursorX for the given line.
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

	// Approximate several-space tab character clicks
	diff := visualX - oldVisualX
	if diff > 1 && x <= oldVisualX+diff/2 {
		i--
	}

	return i
}

// CursorXToVisual returns x, recalculated as visual X for the given line.
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

	t.DeleteSelectedText()
	t.curLine.InsertChar(t.cursorX, r)

	t.cursorX++
	t.UpdateCursorMem()
	t.SelectionAfter(false)
	t.setEdited(true)
}

func (t *TextImpl) HandleCut() {
	t.HandleCopy()
	t.DeleteSelectedText()
}

func (t *TextImpl) HandleCopy() {
	text := t.SelectedText()
	sdl.SetClipboardText(text)
}

func (t *TextImpl) HandlePaste() {
	text, err := sdl.GetClipboardText()
	if err == nil {
		t.DeleteSelectedText()
		t.InsertText(text)
	}
}

func (t *TextImpl) HandleSelectAll() {
	t.SetSelection(1, 0, t.lineCount, len(t.last.chars))
}

func (t *TextImpl) SelectedText() string {
	var b strings.Builder
	if t.selected {
		sel := t.selection
		// First line of selection
		line, lineNum := t.LineByNum(sel.LineFrom)
		if sel.LineFrom == sel.LineTo { // One line selected
			b.WriteString(line.StringRange(sel.CharFrom, sel.CharTo))
		} else { // Two or more lines
			b.WriteString(line.StringRange(sel.CharFrom, len(line.chars)))
			b.WriteString(newLineTypeToString(line.NewLineType))
			lineNum++
			line = line.next
			// Selection inner lines
			for lineNum != sel.LineTo {
				b.WriteString(line.StringRange(0, len(line.chars)))
				b.WriteString(newLineTypeToString(line.NewLineType))
				lineNum++
				line = line.next
			}
			// Last line of selection
			b.WriteString(line.StringRange(0, sel.CharTo))
		}
	}
	return b.String()
}

func (t *TextImpl) InsertText(text string) {
	for _, ch := range text {
		if ch == '\n' {
			t.HandleEnter()
		} else if ch != '\r' {
			t.HandleChar(ch)
		}
	}
	t.setEdited(true)
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
	t.MoveToCursor()
	t.UpdatePos()
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
	t.UpdatePos()
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
	t.lineCount++
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

		if l == t.topLine {
			if t.topLine.prev != nil {
				t.topLine = t.topLine.prev
			} else {
				t.topLine = t.first
			}
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
		t.lineCount--
		t.ScrollDelta(0)
		t.MoveToCursor()
	}
}

func (t *TextImpl) ColorizeSelection(color int) {
	if t.selected {
		sel := t.selection
		// First line of selection
		line, lineNum := t.LineByNum(sel.LineFrom)
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

