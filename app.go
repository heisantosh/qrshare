package main

import (
	"github.com/gotk3/gotk3/gtk"
	"os"
)

const (
	APP_ID = "com.github.mubitosh.qrshare"
)

type App struct {
	file      *string
	inactive  *int // time in seconds
	dataDir   string
	image     string
	isCmdLine bool
	gtkApp    *gtk.Application
}

func (a *App) activate(g *gtk.Application) {
	window := mainWindowNew(a)
	window.ShowAll()
}

func (a *App) cmdLine(g *gtk.Application) {
	if len(os.Args) != 3 {
		a.activate(g)
		return
	}
	a.isCmdLine = true
	window := qrWindowNew(a)
	window.ShowAll()
}
