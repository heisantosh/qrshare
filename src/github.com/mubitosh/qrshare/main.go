package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"log"
	"os"
	"path/filepath"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	qrshare := new(QrShare)

	h := os.Getenv("XDG_DATA_HOME")
	if h == "" {
		h = filepath.Join(os.Getenv("HOME"), ".local/share")
	}
	dir := filepath.Join(h, appID)
	os.MkdirAll(dir, 0775)
	qrshare.image = filepath.Join(dir, "qrimage.png")

	qrshare.files = os.Args[1:]

	qrshare.Application, _ = gtk.ApplicationNew(appID,
		glib.APPLICATION_HANDLES_COMMAND_LINE|glib.APPLICATION_NON_UNIQUE)
	qrshare.Connect("activate", qrshare.activate)
	qrshare.Connect("command-line", qrshare.commandLine)

	initI18n()

	qrshare.Run(os.Args)
}
