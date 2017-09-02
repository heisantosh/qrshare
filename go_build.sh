#!/bin/bash

# Generate Gresource files
cd "${MESON_SOURCE_ROOT}"/data
python gen-gresource-go.py
bash compile-gresource.sh

export GOPATH="${MESON_SOURCE_ROOT}"
go_project_home=$GOPATH/src/github.com/mubitosh/qrshare

cd $GOPATH && mkdir -p src pkg bin $go_project_home

go get -u github.com/golang/dep/cmd/dep
export PATH=$PATH:$GOPATH/bin

mv vendor $go_project_home/
mv *.go *.c *.h $go_project_home/

which go
go version

cd $go_project_home
dep ensure
go build -i -o com.github.mubitosh.qrshare -tags gtk_3_18
mv com.github.mubitosh.qrshare $GOPATH/bin/com.github.mubitosh.qrshare