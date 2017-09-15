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
	// Time in seconds to wait before share stops due to inactivity
	inactive  *int
	dataDir   string
	image     string
	// true if the app was called using contractor
	isContractor bool
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
	a.isContractor = true
	window := qrWindowNew(a)
	window.ShowAll()
}
