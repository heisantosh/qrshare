package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"flag"
	"log"
	"os"
	// "os/exec"
	// "strings"
)

var app *App

// Not sure if a dark background makes it harder to scan QR code??!!
/*
func setDarkTheme() {
	cmdStr := `gsettings get org.gnome.desktop.interface gtk-theme | tr -d "'"`
	out, err := exec.Command("bash", "-c", cmdStr).Output()
	if err != nil {
		log.Println("Command:", cmdStr)
		log.Println("Output:", out)
		log.Println("Failed to set dark theme:", err)
	}
	darkTheme := strings.Trim(string(out), " \n") + ":dark"
	os.Setenv("GTK_THEME", darkTheme)
}
*/

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
	// setDarkTheme()
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
