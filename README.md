# QR Share

QR Share allows you to easily share files by scanning a QR code.

<table cellspacing="0" cellpadding="0" style="border:none">
	<tr>
		<td><img src="https://raw.githubusercontent.com/mubitosh/qrshare/master/data/screenshot-main-window.png"></td>
		<td><img src="https://raw.githubusercontent.com/mubitosh/qrshare/master/data/screenshot-qr-window.png"></td>
	</tr>
</table>

To stop sharing click the "Stop share" button. Sharing stops and the application window is closed automatically if no one is downloading the shared file.

The device running QR Share app and the device scanning the QR code must be on the same network. For example, I have QR Share app running on my laptop and use my smartphone to scan the QR code when both are on the same network.

## Download the app

You can download the latest release of the app from [here](https://github.com/mubitosh/qrshare/releases)

## Building and installing

### Clone the git repo

```bash
$ mkdir -p ~/Devel/debs && cd ~/Devel/debs
$ git clone https://github.com/mubitosh/qrshare.git
```

### Compile & Run the app

[dep](https://github.com/golang/dep) tool is used to make sure the dependencies are in place. Make sure to have it installed. Then we can run go build to build the binary.

```bash
$ cd ~/Devel/debs/qrshare
$ dep ensure
$ go build -ldflags="-s -w" -i -o ./bin/com.github.mubitosh.qrshare -tags gtk_3_18
$ ./bin/com.github.mubitosh.qrshare
```

### Building a deb package

To build a package, copy the binary built in the previous step to the bin directory. Then run dpkg-buildpackage to get a .deb file.

```bash
$ cd ~/Devel/debs/qrshare
$ dpkg-buildpackage
```

### Installing

To install the package dpkg can be used as below. Replace the RELEASE_NUMBER with the correct value.

```bash
$ sudo dpkg -i ~/Devel/debs/com.github.mubitosh.qrshare_RELEASE_NUMBER_amdd64.deb
```

### Uninstalling

```bash
$ sudo dpkg -r com.github.mubitosh.qrshare
```

This project uses dep [https://github.com/golang/dep](https://github.com/golang/dep) for golang dependency management.

Go bindings for GTK3 gotk3 [https://github.com/gotk3/gotk3](https://github.com/gotk3/gotk3)

A barcode creation lib for golang [https://github.com/boombuler/barcode](https://github.com/boombuler/barcode)

## How it works

The components are a file server, a contractor file to have an option in the right click menu and a QR encoder. When the QR Share option is selected from the right click menu, the app starts a file server in the background. A link is generated in the form [http://default-interface-ip-address:random-port-number/shared-file-name](#how-it-works). The QR encoder simply encodes this link and window displays the generated QR code. This code can be scanned by any app which can recognise QR code. Afer clicking the link, the file can be downloaded. To stop the file server from running forever, a timer runs in the background. If no activity happens within a grace of 30 seconds (default value), the timer automatically stops the app. This also shuts down the file server.

## Note

This app was built using [elementary OS](https://elementary.io) Loki. It should hopefully run on other Debian/Ubuntu based distros.

### About downloading files:

The sharing will stop after 30 seconds in the scenario where a video/audio file is being streamed and not being downloaded. During the stream the client may buffer the contents and stay idle. The server will assume there is no activity if no download activity happens within 30 seconds (default value) and it will stop sharing.
If the file is being downloaded, the server will stop only after the download is complete and a grace of 30 seconds.
