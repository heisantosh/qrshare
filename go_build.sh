#!/bin/bash

# Generate Gresource files
cd "${MESON_SOURCE_ROOT}"/data
python gen-gresource-go.py
bash compile-gresource.sh

# Set up GOPATH and move all the necessary files inside the GOPATH
export GOPATH="${MESON_SOURCE_ROOT}"

go_project_home=$GOPATH/src/github.com/mubitosh/qrshare

cd $GOPATH && mkdir -p src pkg bin $go_project_home

mv vendor $go_project_home/
mv *.go *.c *.h $go_project_home/
# So that we can run dep ensure to update the project dependencies
mv Gopkg.* $go_project_home/

go env

cd $go_project_home

# go get -u github.com/golang/dep/cmd/dep

ls -l $GOPATH/bin

$GOPATH/bin/dep ensure

go build -ldflags="-s -w" -i -o com.github.mubitosh.qrshare -tags gtk_3_18
mv com.github.mubitosh.qrshare $GOPATH/bin/com.github.mubitosh.qrshare
