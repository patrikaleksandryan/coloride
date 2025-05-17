package main

import (
	"fmt"
	"github.com/patrikaleksandryan/coloride/pkg/editor"
	"github.com/patrikaleksandryan/coloride/pkg/gui"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

const (
	windowWidth  = 1400 ///16 4R
	windowHeight = 1150 ///16 4G
)

func initInterface() { ///y
	window := editor.NewWindow()  ///y
	gui.Append(window)            ///y
	gui.SetFocus(window.Editor()) ///y
} ///y

func handleMouseWheel(e *sdl.MouseWheelEvent) { ///g
	mainFrame.HandleMouseWheel(lastMouseX, lastMouseY, e.PreciseX, e.PreciseY, e.Direction == sdl.MOUSEWHEEL_FLIPPED) ///11g 16G g
} ///g

func run() error {
	err := gui.Init(windowWidth, windowHeight)
	if err != nil {
		return err
	}

	initInterface() ///1 15G

	err = gui.Run()
	if err != nil {
		return err
	}

	gui.Close()

	return nil
}

func main() { ///b
	fmt.Println(`		///b
		COLOR IDE		///b
	`) ///b
	///b
	if err := run(); err != nil { ///b
		fmt.Fprintln(os.Stderr, err) ///b
		os.Exit(1)                   ///b
	} ///b
} ///b
