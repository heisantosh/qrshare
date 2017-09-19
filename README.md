# QR Share

[![Get it on AppCenter](https://appcenter.elementary.io/badge.svg)](https://appcenter.elementary.io/com.github.mubitosh.qrshare)

QR Share allows you to easily share files by scanning a QR code.

<table cellspacing="0" cellpadding="0" style="border:none">
	<tr>
		<td><img src="https://raw.githubusercontent.com/mubitosh/qrshare/master/data/screenshot-app.png"></td>
		<td><img src="https://raw.githubusercontent.com/mubitosh/qrshare/master/data/screenshot-qr-window.png"></td>
	</tr>
</table>

To stop sharing click the "Stop share" button. Sharing stops and the application window is closed automatically if no one is downloading the shared file.

The device running QR Share app and the device scanning the QR code must be on the same network. For example, I have QR Share app running on my laptop and use my smartphone to scan the QR code when both are on the same network.

## Download the app

You can download the latest release of the app from [here](https://github.com/mubitosh/qrshare/releases)

## Building and installing

### Install dependencies

```bash
$ sudo apt install libgtk-3-dev meson golang-go gb
```

### Building & Installing

All commands below are executed in the same directory.

To build the project, check out from github and run `gb build`

```bash
$ git clone https://github.com/mubitosh/qrshare.git
$ cd qrshare
$ gb build -tags gtk_3_18 all
```

To run the app

```bash
$ bin/qrshare-gtk_3_18
```

To install properly, we need to build a `.deb` package.

```bash
$ dpkg-buildpackage
```

The `deb` file should be available in the parent directory of the current directory. Below is an example of building the 0.7.0 release.

```bash
$ ls ..
com.github.mubitosh.qrshare_0.7.0_amd64.changes
com.github.mubitosh.qrshare_0.7.0_amd64.deb
com.github.mubitosh.qrshare_0.7.0.dsc
com.github.mubitosh.qrshare_0.7.0.tar.xz
qrshare
```

To install use the `dpkg` tool. The release number should vary depending upon the current release.

```bash
$ sudo dpkg -i ../com.github.mubitosh.qrshare_0.7.0_amd64.deb
```

To unistall, just the below command. This one can be executed from anywhere.

```bash
$ sudo dpkg -r com.github.mubitosh.qrshare
```

This project uses [gb](https://getgb.io/) for golang project management.

Go bindings for GTK3 gotk3 [https://github.com/gotk3/gotk3](https://github.com/gotk3/gotk3)

A barcode creation lib for golang [https://github.com/boombuler/barcode](https://github.com/boombuler/barcode)

## How it works

The components are a file server, a contractor file to have an option in the right click menu and a QR encoder. When the QR Share option is selected from the right click menu, the app starts a file server in the background. A link is generated in the form [http://default-interface-ip-address:random-port-number/shared-file-name](#how-it-works). The QR encoder simply encodes this link and window displays the generated QR code. This code can be scanned by any app which can recognise QR code. Afer clicking the link, the file can be downloaded. To stop the file server from running forever, a timer runs in the background. If no activity happens within a grace of 30 seconds (default value), the timer automatically stops the app. This also shuts down the file server.

## Note

This app was built using [elementary OS](https://elementary.io).

### About downloading files:

The sharing will stop after 30 seconds in the scenario where a video/audio file is being streamed and not being downloaded. During the stream the client may buffer the contents and stay idle. The server will assume there is no activity if no download activity happens within 30 seconds (default value) and it will stop sharing.
If the file is being downloaded, the server will stop only after the download is complete and a grace of 30 seconds.
