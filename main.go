package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"flag"
	"log"
	"os"
)

var app *App

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	h := os.Getenv("XDG_DATA_HOME")
	if h == "" {
		h = os.Getenv("HOME") + "/.local/share"
	}
	app = new(App)
	app.dataDir = h + "/" + APP_ID
	os.MkdirAll(app.dataDir, 0775)
	app.image = app.dataDir + "/qrimage.png"
	app.file = flag.String("file", "", "Path of the file to be shared")
	app.inactive = flag.Int("inactive", 30,
		"Sharing is stopped if no sharing activity happens within a period of inactive seconds")
}

func main() {
	flag.Parse()
	log.Println("File to share:", *app.file, " , inactive time:", *app.inactive)
	app.gtkApp, _ = gtk.ApplicationNew(APP_ID,
		glib.APPLICATION_HANDLES_COMMAND_LINE)
	app.gtkApp.Connect("activate", app.activate)
	app.gtkApp.Connect("command-line", app.cmdLine)
	app.gtkApp.Run(os.Args)
}
