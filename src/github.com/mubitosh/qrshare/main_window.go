package main

import (
	"C"

	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"

	"log"
	"os/exec"
	"os/user"
	"unsafe"
)

// mainWindowNew returns Granite Welcome screen style window.
func mainWindowNew(qrshare *QrShare) *gtk.ApplicationWindow {
	titleLabel, _ := gtk.LabelNew(T("Share files and folders"))
	styleCtx, _ := titleLabel.GetStyleContext()
	styleCtx.AddClass("h1")
	titleLabel.SetJustify(gtk.JUSTIFY_CENTER)
	titleLabel.SetHExpand(true)

	subtitleLabel, _ := gtk.LabelNew(T("Use any of the options below to share\nScan the QR code to download"))
	styleCtx, _ = subtitleLabel.GetStyleContext()
	styleCtx.AddClass("h2")
	styleCtx.AddClass("dim-label")
	subtitleLabel.SetJustify(gtk.JUSTIFY_CENTER)
	subtitleLabel.SetLineWrap(true)
	subtitleLabel.SetLineWrapMode(pango.WRAP_WORD)
	subtitleLabel.SetHExpand(true)

	browseButton := optionButtonNew(T("Select files or folders"),
		T("Click here to select files or folders for sharing"),
		"text-x-preview")

	rightClickButton := optionButtonNew(T("Right click in Files application"),
		T("Select files and folders and select the QR Share option"),
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
		qrshare.files = selectFiles(&window.Window)
		// No file was selected
		if len(qrshare.files) == 0 {
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

// selectFiles displays a file chooser dialog to select multiple files.
// It returns a list of names of the selected files.
// Folders cannot be selected.
func selectFiles(window *gtk.Window) []string {
	files := []string{}
	chooser, _ := gtk.FileChooserDialogNewWith2Buttons(T("Select files or folders to share"),
		window,
		gtk.FILE_CHOOSER_ACTION_OPEN&gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER,
		// I hope something like this was allowed to select both files and folders.
		// gtk.FILE_CHOOSER_ACTION_OPEN & gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER,
		T("Cancel"), gtk.RESPONSE_CANCEL,
		T("Select"), gtk.RESPONSE_ACCEPT)
	chooser.SetSelectMultiple(true)
	response := chooser.Run()
	if response == int(gtk.RESPONSE_ACCEPT) {
		list, err := chooser.GetFilenames()
		if err != nil {
			log.Println("Error getting selected filenames:", err)
		} else {
			list.Foreach(func(ptr unsafe.Pointer) {
				files = append(files, C.GoString((*C.char)(ptr)))
			})
		}
		list.Free()
	}
	chooser.Destroy()

	return files
}
