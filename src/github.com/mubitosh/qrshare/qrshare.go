package main

import (
	"github.com/gotk3/gotk3/gtk"

	"os"
)

const (
	appID = "com.github.mubitosh.qrshare"
)

// QrShare represents the state of the QR Share application.
type QrShare struct {
	// Name of file being shared.
	file *string
	// Sharing will stop if no sharing happends during inActive seconds.
	inActive *int
	// Location of QR image.
	image string
	// true if the QR image is displayed using contractor option from right click context menu.
	isContractor bool
	*gtk.Application
}

func (a *QrShare) activate(g *gtk.Application) {
	settings, _ := gtk.SettingsGetDefault()
	settings.Set("gtk-application-prefer-dark-theme", true)
	window := mainWindowNew(a)
	window.ShowAll()
}

func (a *QrShare) commandLine(g *gtk.Application) {
	if len(os.Args) != 3 {
		a.activate(g)
		return
	}
	a.isContractor = true
	window := qrWindowNew(a)
	window.ShowAll()
}
