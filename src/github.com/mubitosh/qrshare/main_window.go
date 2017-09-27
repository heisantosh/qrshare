package main

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"

	"os/exec"
	"os/user"
)

// mainWindowNew returns Granite Welcome screen style window.
func mainWindowNew(qrshare *QrShare) *gtk.ApplicationWindow {
	titleLabel, _ := gtk.LabelNew("Share a file with QR Share")
	styleCtx, _ := titleLabel.GetStyleContext()
	styleCtx.AddClass("h1")
	titleLabel.SetJustify(gtk.JUSTIFY_CENTER)
	titleLabel.SetHExpand(true)

	subtitleLabel, _ := gtk.LabelNew("Use any of the options below to share\n" +
		"Scan the QR code to download the file")
	styleCtx, _ = subtitleLabel.GetStyleContext()
	styleCtx.AddClass("h2")
	styleCtx.AddClass("dim-label")
	subtitleLabel.SetJustify(gtk.JUSTIFY_CENTER)
	subtitleLabel.SetLineWrap(true)
	subtitleLabel.SetLineWrapMode(pango.WRAP_WORD)
	subtitleLabel.SetHExpand(true)

	browseButton := optionButtonNew("Select a file",
		"Click here to select a file for sharing",
		"text-x-preview")

	rightClickButton := optionButtonNew("Right Click in Files",
		"Right click on any file in Files and select the QR Share option",
		"system-file-manager")

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

	window, _ := gtk.ApplicationWindowNew(qrshare.Application)
	window.SetTitle("QR Share")
	window.SetSizeRequest(500, 500)
	window.SetResizable(false)
	window.Add(grid)

	browseButton.Connect("clicked", func() {
		*qrshare.file = chooseFile(&window.Window)
		// No file was selected
		if *qrshare.file == "" {
			return
		}
		qrWindow := qrWindowNew(qrshare)
		qrWindow.ShowAll()
	})

	rightClickButton.Connect("clicked", func() {
		openFilesApp()
	})

	return window
}

func openFilesApp() {
	youser, _ := user.Current()
	cmd := exec.Command("pantheon-files", youser.HomeDir)
	cmd.Start()
}

func chooseFile(window *gtk.Window) string {
	file := ""
	chooser, _ := gtk.FileChooserDialogNewWith2Buttons("Select a file to share",
		window,
		gtk.FILE_CHOOSER_ACTION_OPEN,
		"Cancel", gtk.RESPONSE_CANCEL,
		"Select", gtk.RESPONSE_ACCEPT)
	response := chooser.Run()
	if response == int(gtk.RESPONSE_ACCEPT) {
		file = chooser.GetFilename()
	}
	chooser.Destroy()

	return file
}
