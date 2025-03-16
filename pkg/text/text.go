package text

import "fmt"

type Text interface {
	HandleDelete()
	HandleBackspace()
	HandleEnter()
	HandleLeft()
	HandleRight()
	HandleChar(r rune)

	//!TODO redo this
	Body() []rune
	Cursor() int
}

type TextImpl struct {
	body   []rune
	cursor int // offset of text cursor from start of body

	first, last  *Line
	lineOnTop    *Line // First line visible on the screen
	lineOnTopNum int
}

type Line struct {
	chars []rune
	//!TODO add color data
	prev, next *Line
}

func NewLine() *Line {
	return &Line{}
}

func NewText() Text {
	line := NewLine()
	return &TextImpl{
		first:        line,
		last:         line,
		lineOnTop:    line,
		lineOnTopNum: 1,
	}
}

func (t *TextImpl) Body() []rune {
	return t.body
}

func (t *TextImpl) Cursor() int {
	return t.cursor
}

func (t *TextImpl) HandleDelete() {
	if len(t.body) > t.cursor {
		t.body = append(t.body[:t.cursor], t.body[t.cursor+1:]...)
	}
	t.show()
}

func (t *TextImpl) HandleBackspace() {
	if t.cursor > 0 {
		t.body = append(t.body[:t.cursor-1], t.body[t.cursor:]...)
		t.cursor--
	}
	t.show()
}

func (t *TextImpl) HandleEnter() {
	t.HandleChar('^')
}

func (t *TextImpl) HandleLeft() {
	if t.cursor > 0 {
		t.cursor--
	}
	t.show()
}

func (t *TextImpl) HandleRight() {
	if t.cursor < len(t.body) {
		t.cursor++
	}
	t.show()
}

func (t *TextImpl) HandleChar(r rune) {
	t.body = append(t.body, 0)
	copy(t.body[t.cursor+1:], t.body[t.cursor:])
	t.body[t.cursor] = r
	t.cursor++
	t.show()
}

func (t *TextImpl) show() {
	fmt.Printf("%d \"%s\"\n", t.cursor, string(t.body))
}
