#!/bin/bash

export GTK_VERSION=$(pkg-config --modversion gtk+-3.0 | tr . _ | cut -d '_' -f 1-2)
export GTK_BUILD_TAG="gtk_$GTK_VERSION"

export GOPATH="${MESON_SOURCE_ROOT}"
go_project_home=$GOPATH/src/github.com/mubitosh/qrshare
cd $GOPATH && mkdir -p src pkg bin $go_project_home

go version
go env

mv vendor $go_project_home/
mv *.go $go_project_home/

cd $go_project_home
go build -ldflags="-s -w" -i -o com.github.mubitosh.qrshare -tags $GTK_BUILD_TAG
mv com.github.mubitosh.qrshare $GOPATH/bin/com.github.mubitosh.qrshare
