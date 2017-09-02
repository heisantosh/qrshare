package main

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"

	"log"
	"os/exec"
	"os/user"
)

const (
	FILE_SELECTED = -3
)

// mainWindowNew creates an granite Welcome screen style window.
func mainWindowNew(app *App) *gtk.ApplicationWindow {
	titleLabel, _ := gtk.LabelNew("Share a file with QR Share")
	styleCtx, _ := titleLabel.GetStyleContext()
	styleCtx.AddClass("h1")
	titleLabel.SetJustify(gtk.JUSTIFY_CENTER)
	titleLabel.SetHExpand(true)

	subtitleLabel, _ := gtk.LabelNew("Use any of the options below to share\nScan the QR code to download the file")
	styleCtx, _ = subtitleLabel.GetStyleContext()
	styleCtx.AddClass("h2")
	styleCtx.AddClass("dim-label")
	subtitleLabel.SetJustify(gtk.JUSTIFY_CENTER)
	subtitleLabel.SetLineWrap(true)
	subtitleLabel.SetLineWrapMode(pango.WRAP_WORD)
	subtitleLabel.SetHExpand(true)

	browseButton := optionButtonNew("Select a file",
		"Click here to select a file for sharing")

	rightClickButton := optionButtonNew("Right Click in Files",
		"Right click on any file in Files and select the QR Share option")

	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 9)
	box.SetHAlign(gtk.ALIGN_CENTER)
	box.SetMarginTop(12)
	box.SetMarginBottom(12)
	box.SetMarginStart(12)
	box.SetMarginEnd(12)
	box.PackStart(browseButton, false, false, 0)
	box.PackStart(rightClickButton, false, false, 0)

	grid, _ := gtk.GridNew()
	styleCtx, _ = grid.GetStyleContext()
	styleCtx.AddClass("welcome")
	grid.SetHExpand(true)
	grid.SetVExpand(true)
	grid.SetMarginTop(12)
	grid.SetMarginBottom(24)
	grid.SetVAlign(gtk.ALIGN_CENTER)
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	grid.Add(titleLabel)
	grid.Add(subtitleLabel)
	grid.Add(box)

	window, _ := gtk.ApplicationWindowNew(app.gtkApp)
	window.SetTitle("QR Share")
	window.SetSizeRequest(500, 500)
	window.SetResizable(false)
	window.Add(grid)

	browseButton.Connect("clicked", func() {
		log.Println("Browse button clicked: This should open a file choose dialog")
		*app.file = chooseFile(&window.Window)
		// No file was selected
		if *app.file == "" {
			return
		}
		qrWindow := qrWindowNew(app)
		qrWindow.ShowAll()
	})

	rightClickButton.Connect("clicked", func() {
		log.Println("Right click button clicked: This should open Files app")
		openFilesApp()
	})

	return window
}

func openFilesApp() {
	youser, _ := user.Current()
	cmd := exec.Command("pantheon-files", youser.HomeDir)
	cmd.Start()
	log.Println("Started Files app")
}

func chooseFile(window *gtk.Window) string {
	file := ""
	chooser, _ := gtk.FileChooserDialogNewWith2Buttons("Select a file to share",
		window, gtk.FILE_CHOOSER_ACTION_OPEN, "Select", gtk.RESPONSE_ACCEPT,
		"Cancel", gtk.RESPONSE_CANCEL)
	response := chooser.Run()
	if response == FILE_SELECTED {
		file = chooser.GetFilename()
		log.Println("Selected file:", file)
	}
	chooser.Destroy()

	return file
}
