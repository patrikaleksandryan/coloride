package gui

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var (
	window   *sdl.Window
	renderer *sdl.Renderer
	mainFont *ttf.Font
)

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
		int32(windowWidth), int32(windowHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return fmt.Errorf("could not create renderer: %v", err)
	}

	mainFont, err = ttf.OpenFont("data/fonts/main.ttf", 24)
	if err != nil {
		return fmt.Errorf("could not open font: %v", err)
	}

	return nil
}

func Close() {
	if mainFont != nil {
		mainFont.Close()
	}
	if renderer != nil {
		_ = renderer.Destroy()
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

func handleEvents(running *bool) {
	event := sdl.PollEvent()
	for event != nil {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			*running = false
		case *sdl.KeyboardEvent:
			if e.Keysym.Sym == sdl.K_ESCAPE {
				*running = false
			}
		}
		event = sdl.PollEvent()
	}
}

func render() {
	_ = renderer.SetDrawColor(0, 0, 0, 255)
	_ = renderer.Clear()

	texts := []string{
		"Hello, SDL2!",
		"This is monospaced text.",
		"Rendered using SDL_ttf.",
		"Press ESC to exit.",
	}

	y := int32(50)
	for _, line := range texts {
		err := renderText(renderer, mainFont, line, 50, y)
		if err != nil {
			panic(err)
		}
		y += 40
	}

	renderer.Present()
}

func Run() error {
	running := true
	for running {
		handleEvents(&running)
		render()
	}
	return nil
}
