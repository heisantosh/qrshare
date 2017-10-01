package main

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"

	"os/exec"
	"os/user"
)

// mainWindowNew returns Granite Welcome screen style window.
func mainWindowNew(qrshare *QrShare) *gtk.ApplicationWindow {
	titleLabel, _ := gtk.LabelNew(T("Share a file with QR Share"))
	styleCtx, _ := titleLabel.GetStyleContext()
	styleCtx.AddClass("h1")
	titleLabel.SetJustify(gtk.JUSTIFY_CENTER)
	titleLabel.SetHExpand(true)

	subtitleLabel, _ := gtk.LabelNew(T("Use any of the options below to share\nScan the QR code to download the file"))
	styleCtx, _ = subtitleLabel.GetStyleContext()
	styleCtx.AddClass("h2")
	styleCtx.AddClass("dim-label")
	subtitleLabel.SetJustify(gtk.JUSTIFY_CENTER)
	subtitleLabel.SetLineWrap(true)
	subtitleLabel.SetLineWrapMode(pango.WRAP_WORD)
	subtitleLabel.SetHExpand(true)

	browseButton := optionButtonNew(T("Select a file"),
		T("Click here to select a file for sharing"),
		"text-x-preview")

	rightClickButton := optionButtonNew(T("Right Click in Files"),
		T("Right click on any file in Files and select the QR Share option"),
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
	window.SetSizeRequest(700, 600)
	window.SetResizable(false)
	window.Add(grid)

	browseButton.Connect("clicked", func() {
		*qrshare.file = selectFile(&window.Window)
		// No file was selected
		if *qrshare.file == "" {
			return
		}
		qrWindow := qrWindowNew(qrshare)
		qrWindow.ShowAll()
	})

	rightClickButton.Connect("clicked", func() {
		openFiles()
	})

	return window
}

func openFiles() {
	youser, _ := user.Current()
	cmd := exec.Command("pantheon-files", youser.HomeDir)
	cmd.Start()
}

func selectFile(window *gtk.Window) string {
	file := ""
	chooser, _ := gtk.FileChooserDialogNewWith2Buttons(T("Select a file to share"),
		window,
		gtk.FILE_CHOOSER_ACTION_OPEN,
		T("Cancel"), gtk.RESPONSE_CANCEL,
		T("Select"), gtk.RESPONSE_ACCEPT)
	response := chooser.Run()
	if response == int(gtk.RESPONSE_ACCEPT) {
		file = chooser.GetFilename()
	}
	chooser.Destroy()

	return file
}
