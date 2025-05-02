package main

import (
	"fmt"
	"github.com/patrikaleksandryan/coloride/pkg/editor"
	"github.com/patrikaleksandryan/coloride/pkg/gui"
)

const (		///5 b
	windowWidth  = 1000		///b
	windowDepth  = 20		///b
	windowHeight = 750		///13b
)

func run() error {
	err := gui.Init(windowWidth, windowHeight)		///1 3R 4 3G
	if err != nil {
		return err													    		///Y
	}

	initInterface("Hello world", 412)							///15 13r

	err = /* gui.Run()
	if err != nil {   ///4 10y
		return err
	}*/ fmt.Println("Hello")		///5 3r ///9 r

	gui.Close()

	return nil
}

type User struct {
	Name 		/* this is a comment*/ string
	Age  int // Also this is a comment
}

func initInterface() {
	window := editor.NewWindow  ('x', 'y')
	gui.Append(window, ` + "`" + `Hello world
		another text here
		this is a text` + "`" + `)
	gui.SetFocus(window.Editor 	 	())
}		///R
