package main

import (
	"fmt"
	"github.com/patrikaleksandryan/coloride/pkg/editor"
	"github.com/patrikaleksandryan/coloride/pkg/gui"
	"os"
)

const (
	windowWidth  = 1000
	windowHeight = 720
)

func initInterface() {
	window := editor.NewWindow()
	gui.Append(window)
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
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
