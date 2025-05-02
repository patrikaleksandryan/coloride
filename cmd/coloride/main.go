package main

import (
	"fmt"
	"os"

	"github.com/patrikaleksandryan/coloride/pkg/editor"
	"github.com/patrikaleksandryan/coloride/pkg/gui"
)

const (
	windowWidth  = 1400
	windowHeight = 1150
)

func initInterface() {
	window := editor.NewWindow()
	gui.Append(window)
	gui.SetFocus(window.Editor())
}

func run() error {
	err := gui.Init(windowWidth, windowHeight)
	if err != nil {
		return err
	}

	initInterface()

	err = gui.Run()
	if err != nil {
		return err
	}

	gui.Close()

	return nil
}

func main() {
	fmt.Println(`
		 ▗▄▄▖ ▗▄▖ ▗▖    ▗▄▖ ▗▄▄▖     ▗▄▄▄▖▗▄▄▄  ▗▄▄▄▖
		▐▌   ▐▌ ▐▌▐▌   ▐▌ ▐▌▐▌ ▐▌      █  ▐▌  █ ▐▌   
		▐▌   ▐▌ ▐▌▐▌   ▐▌ ▐▌▐▛▀▚▖      █  ▐▌  █ ▐▛▀▀▘
		▝▚▄▄▖▝▚▄▞▘▐▙▄▄▖▝▚▄▞▘▐▌ ▐▌    ▗▄█▄▖▐▙▄▄▀ ▐▙▄▄▖
	`)

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
