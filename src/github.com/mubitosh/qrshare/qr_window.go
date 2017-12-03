package main

import (
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"

	"image/png"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// getIPAddress returns the IP address assigned to default interface.
func getIPAddress() (string, error) {
	cmdStr := `ip route get "$(ip route show to 0/0 | grep -oP '(?<=via )\S+')" | grep -oP '(?<=src )\S+'`
	out, err := exec.Command("bash", "-c", cmdStr).Output()
	if err != nil {
		return "", err
	}
	ipAddr := strings.Trim(string(out), " \n")
	return ipAddr, nil
}

// genQRCode generates a QR image of URL of the file to be served. It returns the URL.
func genQRCode(qrImage, port string) (string, error) {
	ipAddr, err := getIPAddress()
	if err != nil {
		return "", err
	}

	url := "http://" + ipAddr + ":" + port + sharedPath

	log.Println("URL to share:", url)

	qrCode, _ := qr.Encode(url, qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 300, 300)
	file, _ := os.Create(qrImage)
	defer file.Close()
	png.Encode(file, qrCode)

	return url, nil
}

// alertViewNew returns a gtk Grid similar to AlertView widget from elementary granite library.
func alertViewNew() *gtk.Grid {
	titleLabel, _ := gtk.LabelNew(T("Network is not available"))
	titleLabel.SetHExpand(true)
	styleCtx, _ := titleLabel.GetStyleContext()
	styleCtx.AddClass("h2")
	titleLabel.SetMaxWidthChars(45)
	titleLabel.SetLineWrap(true)
	titleLabel.SetLineWrapMode(pango.WRAP_CHAR)
	titleLabel.SetXAlign(0)

	descriptionLabel, _ := gtk.LabelNew(T("Connect to the same network as the device\nyou will be using to scan the QR code."))
	descriptionLabel.SetHExpand(true)
	descriptionLabel.SetLineWrap(true)
	descriptionLabel.SetUseMarkup(true)
	descriptionLabel.SetXAlign(0)
	descriptionLabel.SetVAlign(gtk.ALIGN_START)
	actionButton, _ := gtk.LinkButtonNewWithLabel("settings://settings/network",
		T("Network Settings..."))
	actionButton.SetMarginTop(24)
	actionButton.SetHAlign(gtk.ALIGN_END)

	image, _ := gtk.ImageNewFromIconName("network-error", gtk.ICON_SIZE_DIALOG)
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
func qrWindowNew(app *QrShare) *gtk.ApplicationWindow {
	fileServer, _ := fileServerNew()

	url, err := genQRCode(app.image, strconv.Itoa(fileServer.port))

	window, _ := gtk.ApplicationWindowNew(app.Application)
	window.SetTitle("QR Share")
	window.SetSizeRequest(400, 400)
	window.SetResizable(false)

	grid, _ := gtk.GridNew()
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	if err != nil {
		grid = alertViewNew()
		window.Add(grid)
		return window
	}

	// Let only close using the Stop Sharing button
	window.Connect("delete-event", func() bool {
		window.Iconify()
		return true
	})

	// Dummy label to get focus instead of urlLabel when opened
	// for the first time
	focusLabel, _ := gtk.LabelNew("")
	focusLabel.SetSelectable(true)

	urlLabel, _ := gtk.LabelNew(url)
	urlLabel.SetMarginStart(12)
	urlLabel.SetMarginEnd(12)
	// urlLabel.SetMarginTop(12)
	urlLabel.SetMarginBottom(12)
	urlLabel.SetSelectable(true)
	styleCtx, _ := urlLabel.GetStyleContext()
	styleCtx.AddClass("h3")

	image, _ := gtk.ImageNewFromFile(app.image)
	image.SetMarginStart(12)
	image.SetMarginEnd(12)
	image.SetMarginTop(12)
	image.SetMarginBottom(12)

	button, _ := gtk.ButtonNewWithLabel(T("Stop Sharing"))
	button.SetMarginBottom(12)
	image.SetMarginTop(12)
	button.SetMarginStart(12)
	button.SetMarginEnd(12)
	button.Connect("clicked", func() {
		if app.isContractor {
			app.Application.Quit()
		} else {
			window.Destroy()
		}
	})

	grid.Add(focusLabel)
	grid.Add(urlLabel)
	grid.Add(image)
	grid.Add(button)
	window.Add(grid)

	go fileServer.start(app, window)

	return window
}
