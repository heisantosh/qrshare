package main

// #cgo pkg-config: gio-2.0
// #include "pokicons.h"
import "C"

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"

	"math/rand"
	"time"
)

var iCount int
var lastIcon int

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	iCount = len(pokicons)
	lastIcon = -1
}

// Select a random icon. Use the paths in pokicons.
func getRandomIcon() *gtk.Image {
	i := rand.Intn(iCount)
	// Very naive way to get a different icon path
	if i == lastIcon {
		i = rand.Intn(iCount)
	}
	lastIcon = i
	icon, _ := gtk.ImageNewFromResource(pokicons[i])
	return icon
}

// optionButtonNew creates a button similar elementary OS granite welcome button.
func optionButtonNew(title string, description string) *gtk.Button {
	titleLabel, _ := gtk.LabelNew(title)
	styleCtx, _ := titleLabel.GetStyleContext()
	styleCtx.AddClass("h3")
	titleLabel.SetHAlign(gtk.ALIGN_START)
	titleLabel.SetVAlign(gtk.ALIGN_END)

	descLabel, _ := gtk.LabelNew(description)
	styleCtx, _ = descLabel.GetStyleContext()
	styleCtx.AddClass("dim-label")
	descLabel.SetHAlign(gtk.ALIGN_START)
	descLabel.SetVAlign(gtk.ALIGN_START)
	descLabel.SetLineWrap(true)
	descLabel.SetLineWrapMode(pango.WRAP_WORD)

	icon := getRandomIcon()
	icon.SetPixelSize(48)
	icon.SetHAlign(gtk.ALIGN_CENTER)
	icon.SetVAlign(gtk.ALIGN_CENTER)

	grid, _ := gtk.GridNew()
	grid.SetColumnSpacing(12)
	grid.SetMarginTop(6)
	grid.SetMarginBottom(6)
	grid.SetMarginStart(6)
	grid.SetMarginEnd(6)
	grid.Attach(icon, 0, 0, 1, 2)
	grid.Attach(titleLabel, 1, 0, 1, 1)
	grid.Attach(descLabel, 1, 1, 1, 1)

	button, _ := gtk.ButtonNew()
	styleCtx, _ = button.GetStyleContext()
	styleCtx.AddClass("flat")
	button.Add(grid)

	return button
}
