package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"flag"
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

	qrshare.inActive = flag.Int("inactive", 30,
		"Sharing is stopped if no sharing activity happens within a period of inactive seconds")

	flag.Parse()
	qrshare.files = flag.Args()

	// Feature: Currently the app accepts absolute pathnames only.
	// This is because the app is supposed to be used with contractor.
	// Make it work with relative pathnames also.

	qrshare.Application, _ = gtk.ApplicationNew(appID, glib.APPLICATION_HANDLES_COMMAND_LINE)
	qrshare.Connect("activate", qrshare.activate)
	qrshare.Connect("command-line", qrshare.commandLine)

	initI18n()

	if len(qrshare.files) > 0 {
		log.Println("Provided", len(qrshare.files), "files through command line or contractor:")
		for _, s := range qrshare.files {
			log.Println(s)
		}
	}

	qrshare.Run(os.Args)
}
