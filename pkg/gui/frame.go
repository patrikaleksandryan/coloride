package gui

type Frame interface {
	Pos() (x, y int)
	SetPos(x, y int)
	Size() (w, h int)
	SetSize(w, h int)
	SetGeometry(x, y, w, h int)
	Visible() bool
	SetVisible(visible bool)
	Enabled() bool
	SetEnabled(enabled bool)
	Click()
	HandleClick(x, y int) bool
	RenderChildren(x, y int)
	Render(x, y int)
	Prev() Frame
	SetPrev(prev Frame)
	Next() Frame
	SetNext(next Frame)
}

// FrameDesc is a base type for all GUI components.
type FrameDesc struct {
	x, y, w, h int
	visible    bool
	enabled    bool
	body       Frame // Ring with a lock
	prev, next Frame

	OnClick func()
}

func InitFrame(frame *FrameDesc, x, y, w, h int) {
	frame.x = x
	frame.y = y
	frame.w = w
	frame.h = h
	frame.visible = true
	frame.enabled = true

	lock := &FrameDesc{}
	lock.prev = lock
	lock.next = lock
	frame.body = lock
}

func (f *FrameDesc) Prev() Frame {
	return f.prev
}

func (f *FrameDesc) SetPrev(prev Frame) {
	f.prev = prev
}

func (f *FrameDesc) Next() Frame {
	return f.next
}

func (f *FrameDesc) SetNext(next Frame) {
	f.next = next
}

func (f *FrameDesc) Append(child Frame) {
	// Add child to the end of the ring (f.body is the lock)
	child.SetNext(f.body)
	child.SetPrev(f.body.Prev())
	f.body.Prev().SetNext(child)
	f.body.SetPrev(child)
}

func (f *FrameDesc) Pos() (x, y int) {
	return f.x, f.y
}

func (f *FrameDesc) SetPos(x, y int) {
	f.x = x
	f.y = y
}

func (f *FrameDesc) Size() (w, h int) {
	return f.w, f.h
}

func (f *FrameDesc) SetSize(w, h int) {
	f.w = w
	f.h = h
}

func (f *FrameDesc) SetGeometry(x, y, w, h int) {
	f.x = x
	f.y = y
	f.w = w
	f.h = h
}

func (f *FrameDesc) Visible() bool {
	return f.visible
}

func (f *FrameDesc) SetVisible(visible bool) {
	f.visible = visible
}

func (f *FrameDesc) Enabled() bool {
	return f.enabled
}

func (f *FrameDesc) SetEnabled(enabled bool) {
	f.enabled = enabled
}

func (f *FrameDesc) Click() {
	if f.OnClick != nil {
		f.OnClick()
	}
}

// HandleClick returns true if click is successful. (x; y) is relative to f.
func (f *FrameDesc) HandleClick(x, y int) bool {
	clicked := false
	if x >= 0 && x < f.w && y >= 0 && y < f.h {
		for c := f.body.Prev(); !clicked && c != f.body; c = c.Prev() {
			if c.Enabled() {
				cx, cy := c.Pos()
				clicked = c.HandleClick(x-cx, y-cy)
			}
		}
		// No child clicked -> clicking on f
		if !clicked {
			f.Click()
			clicked = true
		}
	}
	return clicked
}

func (f *FrameDesc) RenderChildren(x, y int) {
	for c := f.body.Next(); c != f.body; c = c.Next() {
		if c.Visible() {
			cx, cy := c.Pos()
			c.Render(x+cx, y+cy)
		}
	}
}

func (f *FrameDesc) Render(x, y int) {
	f.RenderChildren(x, y)
}
