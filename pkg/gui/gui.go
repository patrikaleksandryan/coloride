package gui

import (
	"fmt"
	"runtime"

	"github.com/patrikaleksandryan/coloride/pkg/font"
	"github.com/patrikaleksandryan/coloride/pkg/text"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	zoom     = 1
	fontSize = 16 * zoom
	charW    = 8 * zoom
	charH    = 16 * zoom
)

var (
	window   *sdl.Window
	Renderer *sdl.Renderer

	mainFont  font.Font
	mainFrame Frame

	focusFrame             Frame
	mouseDownFrame         Frame
	mouseDownX, mouseDownY int

	lastMouseX, lastMouseY int // For mouse wheel event, because MouseX and MouseY are not available

	IsMacOS bool
)

// Append appends the given frame to the main frame.
func Append(frame Frame) {
	mainFrame.Append(frame)
}

func SetFocus(frame Frame) {
	if focusFrame != nil {
		focusFrame.LostFocus()
	}
	focusFrame = frame
	if focusFrame != nil {
		focusFrame.GetFocus()
	}
}

func SetWindowTitle(title string) {
	window.SetTitle(title)
}

func Init(windowWidth, windowHeight int) error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("could not initialize SDL: %v", err)
	}

	err = ttf.Init()
	if err != nil {
		return fmt.Errorf("could not initialize TTF: %v", err)
	}

	window, err = sdl.CreateWindow("ColorIDE", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(windowWidth), int32(windowHeight), sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}

	Renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return fmt.Errorf("could not create Renderer: %v", err)
	}

	mainFont, err = font.Open(Renderer, "data/fonts/main.ttf", fontSize, charW, charH)
	if err != nil {
		return fmt.Errorf("could not open font: %v", err)
	}

	mf := &FrameImpl{}
	InitFrame(mf, 0, 0, 0, 0)
	mainFrame = mf

	return nil
}

func Close() {
	if mainFont != nil {
		mainFont.Close()
	}
	if Renderer != nil {
		_ = Renderer.Destroy()
	}
	if window != nil {
		_ = window.Destroy()
	}
	ttf.Quit()
	sdl.Quit()
}

func ResizeMainFrame(w, h int) {
	oldW, oldH := mainFrame.Size()
	if w != oldW || h != oldH {
		mainFrame.Resize(w, h)
		if mainFrame.HasChildren() {
			child := mainFrame.FirstChild()
			if child != nil {
				SetGeometry(child, 0, 0, w, h)
			}
		}
	}
}

func handleMouseWheel(e *sdl.MouseWheelEvent) {
	mainFrame.HandleMouseWheel(lastMouseX, lastMouseY, e.PreciseX, e.PreciseY, e.Direction == sdl.MOUSEWHEEL_FLIPPED)
}

func handleMouseMove(x, y int, buttons uint32) {
	lastMouseX, lastMouseY = x, y
	if mouseDownFrame != nil {
		mouseDownFrame.MouseMove(x-mouseDownX, y-mouseDownY, buttons)
	} else {
		mainFrame.HandleMouseMove(x, y, buttons)
	}
}

func handleMouseDown(x, y, button int) {
	mouseDownFrame, mouseDownX, mouseDownY =
		mainFrame.HandleMouseDown(x, y, button)
}

func handleMouseUp(x, y, button int) {
	if mouseDownFrame != nil {
		x, y = x-mouseDownX, y-mouseDownY
		mousePressed := mouseDownFrame.MousePressed()
		mouseDownFrame.MouseUp(x, y, button)
		if mousePressed && button == 1 {
			mouseDownFrame.Click()
		}
		mouseDownFrame = nil
	}
}

func handleCharInput(r rune) {
	if focusFrame != nil {
		focusFrame.OnCharInput(r)
	}
}

func handleTextInput(e *sdl.TextInputEvent) {
	t := e.GetText()
	for _, r := range t {
		handleCharInput(r)
	}
}

func handleKeyDown(e *sdl.KeyboardEvent) {
	if focusFrame != nil {
		focusFrame.OnKeyDown(int(e.Keysym.Sym), e.Keysym.Mod)
	}

	switch e.Keysym.Sym {
	case sdl.K_BACKSPACE:
		handleCharInput(text.KeyBackspace)
	case sdl.K_TAB:
		handleCharInput(text.KeyTab)
	case sdl.K_DELETE:
		handleCharInput(text.KeyDelete)
	case sdl.K_RETURN, sdl.K_KP_ENTER:
		handleCharInput(text.KeyEnter)
	}
}

func handleKeyUp(e *sdl.KeyboardEvent) {
	if focusFrame != nil {
		focusFrame.OnKeyUp(int(e.Keysym.Sym), e.Keysym.Mod)
	}
}

func handleKeyboard(e *sdl.KeyboardEvent) {
	if e.State == sdl.PRESSED {
		handleKeyDown(e)
	} else {
		handleKeyUp(e)
	}
}

func handleWindowEvent(event *sdl.WindowEvent) {
	switch event.Event {
	case sdl.WINDOWEVENT_RESIZED:
		w, h := int(event.Data1), int(event.Data2)
		ResizeMainFrame(w, h)
	}
}

func handleEvents(running *bool) {
	event := sdl.PollEvent()
	for event != nil {
		switch e := event.(type) {
		case *sdl.MouseMotionEvent:
			handleMouseMove(int(e.X), int(e.Y), e.State)
		case *sdl.MouseButtonEvent:
			if e.State == sdl.PRESSED {
				handleMouseDown(int(e.X), int(e.Y), int(e.Button))
			} else if e.State == sdl.RELEASED {
				handleMouseUp(int(e.X), int(e.Y), int(e.Button))
			}
		case *sdl.MouseWheelEvent:
			handleMouseWheel(e)
		case *sdl.TextInputEvent:
			handleTextInput(e)
		case *sdl.KeyboardEvent:
			handleKeyboard(e)
		case *sdl.WindowEvent:
			handleWindowEvent(e)
		case *sdl.QuitEvent:
			*running = false
		}
		event = sdl.PollEvent()
	}
}

func render() {
	Renderer.SetDrawColor(0, 0, 0, 255)
	Renderer.Clear()
	mainFrame.Render(0, 0)
	Renderer.Present()
}

func Run() error {
	w, h := window.GetSize()
	ResizeMainFrame(int(w), int(h))

	running := true
	for running {
		handleEvents(&running)
		render()
	}
	return nil
}

// IsCtrlCmdPressed reports whether CMD key is pressed on macOS or CTRL is pressed on other OSes.
func IsCtrlCmdPressed(mod uint16) bool {
	return IsMacOS && mod&sdl.KMOD_GUI != 0 ||
		!IsMacOS && mod&sdl.KMOD_CTRL != 0
}

func init() {
	IsMacOS = runtime.GOOS == "darwin"
}
