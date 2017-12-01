package main

import (
	"github.com/gotk3/gotk3/gtk"
)

const (
	appID = "com.github.mubitosh.qrshare"
)

// QrShare represents the state of the QR Share application.
type QrShare struct {
	files        []string // Absolute paths of files being shared.
	image        string   // Path of QR image.
	isContractor bool     // True if app was used with arguments on command line.
	*gtk.Application
}

func (a *QrShare) activate(g *gtk.Application) {
	settings, _ := gtk.SettingsGetDefault()
	settings.Set("gtk-application-prefer-dark-theme", true)
	window := mainWindowNew(a)
	window.ShowAll()
}

func (a *QrShare) commandLine(g *gtk.Application) {
	// If no files are provided, let guide the user to select files or folders.
	if len(a.files) == 0 {
		a.activate(g)
		return
	}

	a.isContractor = true

	settings, _ := gtk.SettingsGetDefault()
	settings.Set("gtk-application-prefer-dark-theme", true)
	window := qrWindowNew(a)
	window.ShowAll()
}
