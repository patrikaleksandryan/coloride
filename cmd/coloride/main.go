package main

import (
	"fmt"
	"github.com/patrikaleksandryan/coloride/pkg/gui"
	"os"
)

const (
	windowWidth  = 1000
	windowHeight = 720
)

func initInterface() {
	btn1 := gui.NewFrame(gui.NewButton("Hello"), 100, 50, 120, 32)
	btn2 := gui.NewFrame(gui.NewButton("World"), 100, 200, 600, 400)
	btn3 := gui.NewFrame(gui.NewButton("Station"), 50, 30, 200, 200)
	btn4 := gui.NewFrame(gui.NewButton("Mir"), 300, 100, 300, 80)

	btn4.View.SetOnClick(func() {
		fmt.Println("THIS IS BUTTON 4")
	})

	gui.Append(btn1)
	gui.Append(btn2)
	btn2.Append(btn3)
	btn3.Append(btn4)
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
