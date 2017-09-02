#!/bin/bash

# Generate Gresource files
cd "${MESON_SOURCE_ROOT}"/data
python gen-gresource-go.py
bash compile-gresource.sh

export GOPATH="${MESON_SOURCE_ROOT}"
export GOROOT="${MESON_SOURCE_ROOT}"/go
export PATH=$GOROOT/bin:$PATH

cd "${MESON_SOURCE_ROOT}"
if [ ! -f go1.9.linux-amd64.tar.gz ]; then
	curl -O https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz
fi
tar -C $GOPATH -xzf go1.9.linux-amd64.tar.gz

go_project_home=$GOPATH/src/github.com/mubitosh/qrshare

cd $GOPATH && mkdir -p src pkg bin $go_project_home

mv vendor $go_project_home/
mv *.go *.c *.h $go_project_home/

which go
go version

cd $go_project_home
go build -i -o com.github.mubitosh.qrshare -tags gtk_3_18
mv com.github.mubitosh.qrshare $GOPATH/bin/com.github.mubitosh.qrshare