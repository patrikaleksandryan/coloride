package main

import ( ///r
	"fmt"                                               ///r
	"github.com/patrikaleksandryan/coloride/pkg/editor" ///r
	"github.com/patrikaleksandryan/coloride/pkg/gui"    ///r
	"os"                                                ///r
) ///r

const ( ///g
	windowWidth  = 1400 ///g
	windowHeight = 1150 ///g
) ///g

func initInterface() { ///b
	window := editor.NewWindow()  ///b
	gui.Append(window)            ///b
	gui.SetFocus(window.Editor()) ///b
} ///b

func run() error {
	err := gui.Init(windowWidth, windowHeight) ///y
	if err != nil {                            ///y
		return err ///y
	} ///y

	initInterface() ///1 15R

	err = gui.Run() ///1 15G
	if err != nil {
		return err ///2 10B
	}

	gui.Close() ///1 11Y

	return nil
}

func main() {
	fmt.Println(`
		COLOR IDE
	`)

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
