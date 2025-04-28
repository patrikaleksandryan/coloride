package main

import (
	"fmt"
	"os"

	"github.com/patrikaleksandryan/coloride/pkg/editor"
	"github.com/patrikaleksandryan/coloride/pkg/gui"
)

const (		///5 B
	windowWidth  = 1000		///B
	windowDepth  = 20		///B
	windowHeight = 750		///13B
)

func run() error {
	err := gui.Init(windowWidth, windowHeight)		///2 3R! 4 3G!
	if err != nil {
		return err													    		///Y!
	}

	initInterface("Hello world", 412)							///16 13R

	err = /* gui.Run()
	if err != nil {
		return err
	}*/ fmt.Println("Hello")		///6 R3 ///9 R

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
}
