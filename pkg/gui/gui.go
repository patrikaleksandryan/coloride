package gui

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var (
	window   *sdl.Window
	Renderer *sdl.Renderer
	mainFont *ttf.Font

	mainFrame *FrameDesc
)

// Append appends the given frame to the main frame.
func Append(frame Frame) {
	mainFrame.Append(frame)
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

	window, err = sdl.CreateWindow("SDL2 Text Example", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(windowWidth), int32(windowHeight), sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}

	Renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return fmt.Errorf("could not create Renderer: %v", err)
	}

	mainFont, err = ttf.OpenFont("data/fonts/main.ttf", 24)
	if err != nil {
		return fmt.Errorf("could not open font: %v", err)
	}

	mainFrame = &FrameDesc{}
	InitFrame(mainFrame, 0, 0, windowWidth, windowHeight)

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

func renderText(renderer *sdl.Renderer, font *ttf.Font, text string, x, y int32) error {
	color := sdl.Color{R: 0, G: 255, B: 0, A: 255}
	surface, err := font.RenderUTF8Blended(text, color)
	if err != nil {
		return fmt.Errorf("could not create text surface: %v", err)
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("could not create text texture: %v", err)
	}
	defer texture.Destroy()

	rect := sdl.Rect{X: x, Y: y, W: surface.W, H: surface.H}
	return renderer.Copy(texture, nil, &rect)
}

func ResizeMainFrame(w, h int) {
	mainFrame.Resize(w, h)
	if mainFrame.HasChildren() {
		mainFrame.body.Next().SetGeometry(0, 0, w, h)
	}
}

func handleMouseDown(x, y int) {
	mainFrame.HandleClick(x, y)
}

func handleWindowEvent(event *sdl.WindowEvent) {
	switch event.Event {
	case sdl.WINDOWEVENT_RESIZED:
		w, h := int(event.Data1), int(event.Data2)
		ResizeMainFrame(int(w), int(h))
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
		case *sdl.KeyboardEvent:
			if e.Keysym.Sym == sdl.K_ESCAPE {
				*running = false
			}
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
