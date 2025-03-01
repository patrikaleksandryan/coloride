package text

import "fmt"

type Text interface {
	HandleDelete()
	HandleBackspace()
	HandleEnter()
	HandleLeft()
	HandleRight()
	HandleChar(r rune)
}

type TextImpl struct {
	body   []rune
	cursor int // offset of text cursor from start of body
}

func NewText() Text {
	return &TextImpl{}
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
