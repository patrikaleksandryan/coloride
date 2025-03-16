package gui

import (
	"fmt"

	"github.com/patrikaleksandryan/coloride/pkg/font"
	"github.com/patrikaleksandryan/coloride/pkg/text"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	// !TODO rename to "defaultFontSize" etc. and allow to set up from outside the gui
	fontSize = 22
	charW    = 16
	charH    = 26
)

var (
	window   *sdl.Window
	Renderer *sdl.Renderer

	mainFont  font.Font
	mainFrame *FrameImpl //!FIXME change to Frame, add Append

	focusFrame Frame
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

	mainFrame = &FrameImpl{}
	InitFrame(mainFrame, 0, 0, 0, 0)

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
			SetGeometry(mainFrame.body.Next(), 0, 0, w, h)
		}
	}
}

func handleMouseDown(x, y int) {
	mainFrame.HandleClick(x, y)
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
		case *sdl.MouseButtonEvent:
			if e.State == sdl.PRESSED && e.Button == 1 {
				handleMouseDown(int(e.X), int(e.Y))
			}
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
	_ = Renderer.SetDrawColor(0, 0, 0, 255)
	_ = Renderer.Clear()

	//texts := []string{
	//	"Hello, SDL2!",
	//	"This is monospaced text.",
	//	"Rendered using SDL_ttf.",
	//	"Press ESC to exit.",
	//}
	//
	//y := int32(50)
	//for _, line := range texts {
	//	err := renderText(Renderer, mainFont, line, 50, y)
	//	if err != nil {
	//		panic(err)
	//	}
	//	y += 40
	//}

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
