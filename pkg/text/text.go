package text

//const (
//	left = iota
//	right
//	up
//	bottom
//)

type Text interface {
	HandleDelete()
	HandleBackspace()
	HandleEnter()
	HandleLeft(shift bool)
	HandleRight(shift bool)
	HandleUp(shift bool)
	HandleDown(shift bool)
	HandleChar(r rune)

	CurLine() *Line
	SetCurLine(lineNum int)
	TopLine() (*Line, int)
	CursorX() int
	SetCursorX(cursorX int)

	InSelection(lineNum, charNum int) bool
	IsLeftSelectionEdge() bool
	ClearSelection()
	SetSelection(lineFrom, lineTo, charFrom, charTo int)

	InsertLineAfter(l *Line) *Line
	MergeLines(l *Line)
	DeleteLine(l *Line)

	UpdateCursorMem()
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

	selected  bool
	selection Selection

	first, last *Line // First and last line of document
	topLine     *Line // First line visible on the screen
	topLineNum  int   // 1-based
}

type Line struct {
	chars []rune
	//!TODO add color data
	prev, next *Line
}

func NewLine() *Line {
	return &Line{}
}

func (l *Line) Chars() []rune {
	return l.chars
}

func (l *Line) Prev() *Line {
	return l.prev
}

func (l *Line) Next() *Line {
	return l.next
}

func NewText() Text {
	line := NewLine()
	text := &TextImpl{
		first:      line,
		last:       line,
		curLine:    line,
		curLineNum: 1,
		topLine:    line,
		topLineNum: 1,
	}
	text.setDummyText()
	return text
}

func (t *TextImpl) setDummyText() {
	s := `MODULE TestReader;
IMPORT Out, Texts;
VAR T: Texts.Text; R: Texts.Reader; S: Texts.Scanner;
  ch: CHAR;
BEGIN
  NEW(T); Texts.Open(T, 'Data/TEXT.DAT');

  Texts.OpenScanner(S, T, 1); Texts.Scan(S);
  Out.String(S.s); Out.Char(';'); Out.Ln;

  Texts.OpenReader(R, T, 0);
  Texts.Read(R, ch);
  WHILE ~R.eot DO
    Out.Int(ORD(R.eot), 5);
    Out.Int(ORD(ch), 5); Out.String('   ');
    Out.Char(ch); Out.Ln;
    Texts.Read(R, ch)
  END
END TestReader.`

	for _, r := range s {
		if r == 0xA {
			t.HandleEnter()
		} else {
			t.HandleChar(r)
		}
	}
}

func (t *TextImpl) UpdateCursorMem() {
	t.cursorMem = t.cursorX
}

func (t *TextImpl) HandleDelete() {
	if len(t.curLine.chars) > t.cursorX {
		t.curLine.chars = append(t.curLine.chars[:t.cursorX], t.curLine.chars[t.cursorX+1:]...)
	} else if t.curLine.next != nil {
		dx := len(t.curLine.chars)
		t.HandleRight(false)
		t.MergeLines(t.curLine.prev)
		t.cursorX = dx
	}
	t.UpdateCursorMem()
}

func (t *TextImpl) HandleBackspace() {
	if t.cursorX > 0 {
		t.curLine.chars = append(t.curLine.chars[:t.cursorX-1], t.curLine.chars[t.cursorX:]...)
		t.cursorX--
	} else if t.curLine.prev != nil {
		dx := len(t.curLine.prev.chars)
		t.MergeLines(t.curLine.prev)
		t.cursorX = dx
	}
	t.UpdateCursorMem()
}

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

func (t *TextImpl) IsLeftSelectionEdge() bool {
	return t.curLineNum == t.selection.LineFrom && t.cursorX == t.selection.CharFrom
}

func (t *TextImpl) ClearSelection() {
	t.selected = false
}

func (t *TextImpl) SetSelection(lineFrom, lineTo, charFrom, charTo int) {
	t.selected = true
	t.selection.LineFrom = lineFrom
	t.selection.LineTo = lineTo
	t.selection.CharFrom = charFrom
	t.selection.CharTo = charTo
}

func (t *TextImpl) HandleEnter() {
	newLine := t.InsertLineAfter(t.curLine)
	newLine.chars = append(newLine.chars, t.curLine.chars[t.cursorX:]...)
	t.curLine.chars = t.curLine.chars[:t.cursorX]
	t.curLine = newLine
	t.curLineNum++
	t.cursorX = 0
	t.UpdateCursorMem()
}

func (t *TextImpl) handleSelectLeft(shift bool) {
	if !shift {
		t.ClearSelection()
	} else if t.cursorX > 0 {
		if !t.selected {
			t.SetSelection(t.curLineNum, t.curLineNum, t.cursorX-1, t.cursorX)
		} else if t.IsLeftSelectionEdge() {
			t.selection.CharFrom--
		} else { // On the right edge of selection
			t.selection.CharTo--
			if t.selection.LineFrom == t.selection.LineTo &&
				t.selection.CharFrom == t.selection.CharTo {
				t.ClearSelection()
			}
		}
	} else if t.curLine.prev != nil {
		lineLen := len(t.curLine.prev.chars)
		if !t.selected {
			t.SetSelection(t.curLineNum-1, t.curLineNum, lineLen, 0)
		} else {
			t.selection.CharFrom = lineLen
			t.selection.LineFrom--
		}
	}
}

func (t *TextImpl) HandleLeft(shift bool) {
	t.handleSelectLeft(shift)
	if t.cursorX > 0 {
		t.cursorX--
	} else if t.curLine.prev != nil {
		t.curLineNum--
		t.curLine = t.curLine.prev
		t.cursorX = len(t.curLine.chars)
	}
	t.UpdateCursorMem()
}

func (t *TextImpl) handleSelectRight(shift bool) {
	if !shift {
		t.ClearSelection()
	} else if t.cursorX < len(t.curLine.chars) {
		if !t.selected {
			t.SetSelection(t.curLineNum, t.curLineNum, t.cursorX, t.cursorX+1)
		} else if t.IsLeftSelectionEdge() {
			t.selection.CharFrom++
			if t.selection.LineFrom == t.selection.LineTo &&
				t.selection.CharFrom == t.selection.CharTo {
				t.ClearSelection()
			}
		} else { // On the right edge of selection
			t.selection.CharTo++
		}
	} else if t.curLine.next != nil {
		lineLen := len(t.curLine.chars)
		if !t.selected {
			t.SetSelection(t.curLineNum, t.curLineNum+1, lineLen, 0)
		} else {
			t.selection.CharTo = 0
			t.selection.LineTo++
		}
	}
}

func (t *TextImpl) HandleRight(shift bool) {
	t.handleSelectRight(shift)
	if t.cursorX < len(t.curLine.chars) {
		t.cursorX++
	} else if t.curLine.next != nil {
		t.curLineNum++
		t.curLine = t.curLine.next
		t.cursorX = 0
	}
	t.UpdateCursorMem()
}

func (t *TextImpl) handleSelectUp(shift bool) {
	if !shift {
		t.ClearSelection()
	} else if t.curLine.prev != nil {
		if !t.selected {
			charFrom := len(t.curLine.prev.chars)
			if t.cursorX < charFrom {
				charFrom = t.cursorX
			}
			t.SetSelection(t.curLineNum-1, t.curLineNum, charFrom, t.cursorX)
		} else if t.IsLeftSelectionEdge() {
			t.selection.LineFrom--
			length := len(t.curLine.prev.chars)
			if t.cursorMem > length {
				t.selection.CharFrom = length
			} else {
				t.selection.CharFrom = t.cursorMem
			}
		} else { // On the right edge of selection
			t.selection.LineTo--
			length := len(t.curLine.prev.chars)
			if t.selection.CharTo > length {
				t.selection.CharTo = length
			}

			if t.selection.LineFrom == t.selection.LineTo &&
				t.selection.CharFrom == t.selection.CharTo {
				t.ClearSelection()
			} else if t.selection.LineFrom > t.selection.LineTo || t.selection.CharFrom > t.selection.CharTo {
				t.selection.CharFrom, t.selection.CharTo = t.selection.CharTo, t.selection.CharFrom
				t.selection.LineFrom, t.selection.LineTo = t.selection.LineTo, t.selection.LineFrom
			}
		}
	}
}

func (t *TextImpl) HandleUp(shift bool) {
	t.handleSelectUp(shift)
	if t.curLine.prev != nil {
		t.curLineNum--
		t.curLine = t.curLine.prev

		t.cursorX = t.cursorMem
		if t.cursorX > len(t.curLine.chars) {
			t.cursorX = len(t.curLine.chars)
		}
	}
}

func (t *TextImpl) handleSelectDown(shift bool) {
	//!TODO
}

func (t *TextImpl) HandleDown(shift bool) {
	t.handleSelectDown(shift)
	if t.curLine.next != nil {
		t.curLineNum++
		t.curLine = t.curLine.next

		t.cursorX = t.cursorMem
		if t.cursorX > len(t.curLine.chars) {
			t.cursorX = len(t.curLine.chars)
		}
	}
}

func (t *TextImpl) HandleChar(r rune) {
	l := t.curLine
	l.chars = append(l.chars, 0)
	copy(l.chars[t.cursorX+1:], l.chars[t.cursorX:])
	l.chars[t.cursorX] = r
	t.cursorX++
	t.UpdateCursorMem()
}

func (t *TextImpl) CurLine() *Line {
	return t.curLine
}

func (t *TextImpl) SetCurLine(lineNum int) {
	if lineNum > 0 {
		l := t.first
		n := 1
		for l.next != nil && n != lineNum {
			l = l.next
			n++
		}
		t.curLine = l
		t.curLineNum = n
		t.cursorX = 0
		t.UpdateCursorMem()
	}
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

func (t *TextImpl) InsertLineAfter(l *Line) *Line {
	l2 := NewLine()
	if l == t.last {
		t.last = l2
	}

	l2.next = l.next
	l2.prev = l
	if l.next != nil {
		l.next.prev = l2
	}
	l.next = l2

	return l2
}

// MergeLines merges line l with the next line moving up, appending its contents to l.
func (t *TextImpl) MergeLines(l *Line) {
	if l.next != nil {
		l.chars = append(l.chars, l.next.chars...)
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
