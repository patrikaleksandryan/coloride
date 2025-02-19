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
	btn1 := gui.NewButton("Hello", 30, 30, 300, 60)
	btn2 := gui.NewButton("World", 30, 120, 300, 300)
	btn3 := gui.NewButton("Station", 60, 60, 210, 360)
	btn4 := gui.NewButton("Mir", 30, 30, 360, 90)

	btn4.OnClick = func() {
		fmt.Println("THIS IS BUTTON 4")
	}

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
