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
	btn2a := gui.NewButton("Small", 260, 10, 100, 30)
	btn3 := gui.NewButton("Station", 60, 60, 210, 360)
	btn4 := gui.NewButton("Mir", 30, 30, 360, 250)

	btn2.SetBgColor(gui.MakeColor(220, 50, 10))
	btn2a.SetBgColor(gui.MakeColor(40, 40, 40))
	btn3.SetBgColor(gui.MakeColor(50, 220, 50))
	btn4.SetBgColor(gui.MakeColor(220, 200, 50))

	btn2.SetColor(gui.MakeColor(255, 255, 255))
	btn2a.SetColor(gui.MakeColor(40, 255, 40))
	btn3.SetColor(gui.MakeColor(0, 255, 0))
	btn4.SetColor(gui.MakeColor(255, 0, 0))

	btn4.OnClick = func() {
		fmt.Println("THIS IS BUTTON 4")
	}

	gui.Append(btn1)
	gui.Append(btn2)
	btn2.Append(btn2a)
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
