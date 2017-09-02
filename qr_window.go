package main

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os/exec"
	"path"
	"strings"
	"net/url"
	"strconv"
)

// getIPAddress returns the IP address assigned to default interface.
func getIPAddress() string {
	cmdStr := `ip route get "$(ip route show to 0/0 | grep -oP '(?<=via )\S+')" | grep -oP '(?<=src )\S+'`
	out, err := exec.Command("bash", "-c", cmdStr).Output()
	if err != nil {
		log.Println(out)
		log.Fatal("Failed to get default IP address of this machine:", err)
		return ""
	}
	ipAddr := strings.Trim(string(out), " \n")
	return ipAddr
}

func genQRCode(baseName, qrImage, port string) {
	url := "http://" + getIPAddress() + ":" + port + "/" + url.PathEscape(baseName)
	log.Println("URL:", url)
	cmdStr := `qrencode -o ` + qrImage + ` -m 0 -s 10 "` + url + `"`
	log.Println("qrencode command:", cmdStr)
	out, err := exec.Command("bash", "-c", cmdStr).Output()
	if err != nil {
		log.Println(out)
		log.Fatal("Failed to generate QR code image for given URL:", err)
	}
}

func qrWindowNew(app *App) *gtk.ApplicationWindow {
	fileServer, _ := FileServerNew()
	baseName := path.Base(*app.file)

	genQRCode(baseName, app.image, strconv.Itoa(fileServer.port))

	window, _ := gtk.ApplicationWindowNew(app.gtkApp)
	window.Connect("delete-event", func() bool {
		window.Iconify()
		return true
	})
	window.SetTitle("QR Share - " + baseName)
	window.SetSizeRequest(400, 400)
	window.SetResizable(false)

	image, _ := gtk.ImageNewFromFile(app.image)
	image.SetMarginStart(12)
	image.SetMarginEnd(12)
	image.SetMarginTop(12)
	image.SetMarginBottom(12)

	button, _ := gtk.ButtonNewWithLabel("Stop Share")
	button.SetMarginBottom(12)
	image.SetMarginTop(12)
	button.SetMarginStart(12)
	button.SetMarginEnd(12)
	button.Connect("clicked", func() {
		log.Println("Stopping sharing")
		if app.isCmdLine {
			log.Println("Closing QR window, stopping app")
			app.gtkApp.Quit()
		} else {
			log.Println("Closing QR window, back to main window")
			window.Destroy()
		}
	})

	grid, _ := gtk.GridNew()
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	grid.Add(image)
	grid.Add(button)

	window.Add(grid)

	go fileServer.start(app, window)

	return window
}
	