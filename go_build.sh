#!/bin/bash

GTK_VERSION=$(pkg-config --modversion gtk+-3.0 | tr . _ | cut -d '_' -f 1-2)
GTK_BUILD_TAG="gtk_$(GTK_VERSION)"

export GOPATH="${MESON_SOURCE_ROOT}"
go_project_home=$GOPATH/src/github.com/mubitosh/qrshare
cd $GOPATH && mkdir -p src pkg bin $go_project_home

cd $GOPATH && wget https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz
tar -C . -xzf go1.9.linux-amd64.tar.gz
export PATH=$GOPATH/go/bin:$GOPATH/bin:$PATH

go version
go env

go get -v -u github.com/golang/dep/cmd/dep

mv vendor $go_project_home/
mv *.go $go_project_home/
mv Gopkg.* $go_project_home/

cd $go_project_home && $GOPATH/bin/dep ensure
go build -ldflags="-s -w" -i -o com.github.mubitosh.qrshare -tags $GTK_BUILD_TAG
mv com.github.mubitosh.qrshare $GOPATH/bin/com.github.mubitosh.qrshare
