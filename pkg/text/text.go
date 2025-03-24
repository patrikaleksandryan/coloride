package text

type Text interface {
	HandleDelete()
	HandleBackspace()
	HandleEnter()
	HandleLeft()
	HandleRight()
	HandleUp()
	HandleDown()
	HandleChar(r rune)

	CurLine() *Line
	TopLine() *Line
	CursorX() int

	InsertLineAfter(l *Line) *Line
	MergeLines(l *Line)
	DeleteLine(l *Line)

	UpdateCursorMem()
}

type TextImpl struct {
	cursorX    int // 0-based, x
	cursorMem  int // 0-based, x
	curLineNum int // 1-based, y
	curLine    *Line

	first, last *Line
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
	return &TextImpl{
		first:      line,
		last:       line,
		curLine:    line,
		topLine:    line,
		topLineNum: 1,
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
		t.HandleRight()
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

func (t *TextImpl) HandleEnter() {
	newLine := t.InsertLineAfter(t.curLine)
	newLine.chars = append(newLine.chars, t.curLine.chars[t.cursorX:]...)
	t.curLine.chars = t.curLine.chars[:t.cursorX]
	t.curLine = newLine
	t.curLineNum++
	t.cursorX = 0
	t.UpdateCursorMem()
}

func (t *TextImpl) HandleLeft() {
	if t.cursorX > 0 {
		t.cursorX--
	} else if t.curLine.prev != nil {
		t.curLineNum--
		t.curLine = t.curLine.prev
		t.cursorX = len(t.curLine.chars)
	}
	t.UpdateCursorMem()
}

func (t *TextImpl) HandleRight() {
	if t.cursorX < len(t.curLine.chars) {
		t.cursorX++
	} else if t.curLine.next != nil {
		t.curLineNum++
		t.curLine = t.curLine.next
		t.cursorX = 0
	}
	t.UpdateCursorMem()
}

func (t *TextImpl) HandleUp() {
	if t.curLine.prev != nil {
		t.curLineNum--
		t.curLine = t.curLine.prev

		t.cursorX = t.cursorMem
		if t.cursorX > len(t.curLine.chars) {
			t.cursorX = len(t.curLine.chars)
		}
	}
}

func (t *TextImpl) HandleDown() {
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

func (t *TextImpl) TopLine() *Line {
	return t.topLine
}

func (t *TextImpl) CursorX() int {
	return t.cursorX
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
