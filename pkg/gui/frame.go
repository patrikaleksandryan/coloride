package gui

type View interface {
	Draw(x, y, w, h int)
	Click()
}

type Frame struct {
	x, y, w, h int
	visible    bool
	enabled    bool
	view       View
	body       *Frame // Ring with a lock
	prev, next *Frame
}

func NewFrame(x, y, w, h int) *Frame {
	lock := &Frame{}
	lock.prev = lock
	lock.next = lock

	return &Frame{
		visible: true,
		enabled: true,
		body:    lock,
	}
}

func (f *Frame) Append(child *Frame) {
	child.next = f.body
	child.prev = f.body.prev
	f.body.prev.next = child
	f.body.prev = child
}

func (f *Frame) Pos() (x, y int) {
	return f.x, f.y
}

func (f *Frame) Size() (w, h int) {
	return f.w, f.h
}

func (f *Frame) Visible() bool {
	return f.visible
}

func (f *Frame) SetVisible(visible bool) {
	f.visible = visible
}

func (f *Frame) Enabled() bool {
	return f.enabled
}

func (f *Frame) SetEnabled(enabled bool) {
	f.enabled = enabled
}

func (f *Frame) Click() {
	f.view.Click()
}

// TryClick returns true if click is successful. (x; y) is relative to f.
func (f *Frame) HandleClick(x, y int) bool {
	clicked := false
	if x >= 0 && x < f.w && y >= 0 && y < f.h {
		for c := f.body.prev; !clicked && c != f.body; c = c.prev {
			clicked = c.HandleClick(x-c.x, y-c.y)
		}
		// No child clicked -> clicking on f
		if !clicked {
			f.Click()
			clicked = true
		}
	}
	return clicked
}
