package main

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
)

// optionButtonNew creates a button similar elementary OS granite welcome button.
func optionButtonNew(title, description, iconName string) *gtk.Button {
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

	icon, _ := gtk.ImageNewFromIconName(iconName, gtk.ICON_SIZE_DIALOG)
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
