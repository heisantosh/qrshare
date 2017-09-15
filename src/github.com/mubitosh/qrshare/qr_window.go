package main

import (
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"

	"image/png"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

const (
	NO_NETWORK      = "No Network"
	NO_QR_GENERATED = "No QR Generated"
	NO_ERROR        = ""
)

// getIPAddress returns the IP address assigned to default interface.
func getIPAddress() (string, string) {
	cmdStr := `ip route get "$(ip route show to 0/0 | grep -oP '(?<=via )\S+')" | grep -oP '(?<=src )\S+'`
	out, err := exec.Command("bash", "-c", cmdStr).Output()
	if err != nil {
		return "", NO_NETWORK
	}
	ipAddr := strings.Trim(string(out), " \n")
	return ipAddr, NO_ERROR
}

func genQRCode(baseName, qrImage, port string) string {
	ipAddr, errStr := getIPAddress()
	if errStr != NO_ERROR {
		return errStr
	}

	u, _ := url.Parse(baseName)
	url := "http://" + ipAddr + ":" + port + "/" + u.EscapedPath()

	qrCode, _ := qr.Encode(url, qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 300, 300)
	file, _ := os.Create(qrImage)
	defer file.Close()
	png.Encode(file, qrCode)

	return NO_ERROR
}

// alertViewNew returns a gtk Grid similar to AlertView widget from elementary granite library.
func alertViewNew(errStr string) *gtk.Grid {
	titleLabel, _ := gtk.LabelNew("Unable to create QR code")
	if errStr == NO_NETWORK {
		titleLabel.SetText("Network is not available")
	}
	titleLabel.SetHExpand(true)
	styleCtx, _ := titleLabel.GetStyleContext()
	styleCtx.AddClass("h2")
	titleLabel.SetMaxWidthChars(45)
	titleLabel.SetLineWrap(true)
	titleLabel.SetLineWrapMode(pango.WRAP_CHAR)
	titleLabel.SetXAlign(0)

	descriptionLabel, _ := gtk.LabelNew("Some dependant components might be missing.\n" +
		"Reinstall the app and try again.")
	if errStr == NO_NETWORK {
		descriptionLabel.SetText("Connect to the same network as the device\n" +
			"you will be using to scan the QR code.")
	}
	descriptionLabel.SetHExpand(true)
	descriptionLabel.SetLineWrap(true)
	descriptionLabel.SetUseMarkup(true)
	descriptionLabel.SetXAlign(0)
	descriptionLabel.SetVAlign(gtk.ALIGN_START)

	actionButton, _ := gtk.LinkButtonNewWithLabel("https://appcenter.elementary.io/" + APP_ID,
		"Go to AppCenter...")
	if errStr == NO_NETWORK {
		actionButton, _ = gtk.LinkButtonNewWithLabel("settings://settings/network",
			"Network Settings...")
	}
	actionButton.SetMarginTop(24)
	actionButton.SetHAlign(gtk.ALIGN_END)

	image, _ := gtk.ImageNewFromIconName("system-software-install", gtk.ICON_SIZE_DIALOG)
	if errStr == NO_NETWORK {
		image.SetFromIconName("network-error", gtk.ICON_SIZE_DIALOG)
	}
	image.SetMarginTop(6)
	image.SetMarginBottom(6)
	image.SetMarginEnd(6)

	grid, _ := gtk.GridNew()
	grid.SetColumnSpacing(12)
	grid.SetRowSpacing(6)
	grid.SetHAlign(gtk.ALIGN_CENTER)
	grid.SetVAlign(gtk.ALIGN_CENTER)
	grid.SetVExpand(true)
	grid.SetMarginTop(24)
	grid.SetMarginBottom(24)
	grid.SetMarginStart(24)
	grid.SetMarginEnd(24)

	grid.Attach(image, 1, 1, 1, 2)
	grid.Attach(titleLabel, 2, 1, 1, 1)
	grid.Attach(descriptionLabel, 2, 2, 1, 1)
	grid.Attach(actionButton, 2, 3, 1, 1)

	return grid
}

// qrWindowNew returns a window that displays the QR code image of the URL from
// where the shared file can be downloaded.
func qrWindowNew(app *App) *gtk.ApplicationWindow {
	fileServer, _ := FileServerNew()
	baseName := path.Base(*app.file)

	errStr := genQRCode(baseName, app.image, strconv.Itoa(fileServer.port))

	window, _ := gtk.ApplicationWindowNew(app.gtkApp)
	window.SetTitle("QR Share - " + baseName)
	window.SetSizeRequest(400, 400)
	window.SetResizable(false)

	grid, _ := gtk.GridNew()
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	if errStr != NO_ERROR {
		grid = alertViewNew(errStr)
		window.Add(grid)

		return window
	}

	// Let only close using the Stop Sharing button
	window.Connect("delete-event", func() bool {
		window.Iconify()
		return true
	})

	image, _ := gtk.ImageNewFromFile(app.image)
	image.SetMarginStart(12)
	image.SetMarginEnd(12)
	image.SetMarginTop(12)
	image.SetMarginBottom(12)

	button, _ := gtk.ButtonNewWithLabel("Stop Sharing")
	button.SetMarginBottom(12)
	image.SetMarginTop(12)
	button.SetMarginStart(12)
	button.SetMarginEnd(12)
	button.Connect("clicked", func() {
		if app.isContractor {
			app.gtkApp.Quit()
		} else {
			window.Destroy()
		}
	})

	grid.Add(image)
	grid.Add(button)

	window.Add(grid)

	go fileServer.start(app, window)

	return window
}
