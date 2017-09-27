package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"flag"
	"log"
	"os"
)

var qrshare *QrShare

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	qrshare = new(QrShare)

	h := os.Getenv("XDG_DATA_HOME")
	if h == "" {
		h = os.Getenv("HOME") + "/.local/share"
	}
	dir := h + "/" + appID
	os.MkdirAll(dir, 0775)
	qrshare.image = dir + "/qrimage.png"

	qrshare.file = flag.String("file", "", "Path of the file to be shared")
	qrshare.inActive = flag.Int("inactive", 30,
		"Sharing is stopped if no sharing activity happens within a period of inactive seconds")
}

func main() {
	flag.Parse()
	qrshare.Application, _ = gtk.ApplicationNew(appID, glib.APPLICATION_HANDLES_COMMAND_LINE)
	qrshare.Connect("activate", qrshare.activate)
	qrshare.Connect("command-line", qrshare.commandLine)
	qrshare.Run(os.Args)
}
