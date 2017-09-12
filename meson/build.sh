#!/bin/bash

mkdir qrshare
cp -R vendor qrshare/
cp -R src qrshare/
cd qrshare
gb build -tags gtk_3_18 all
mv bin/qrshare-gtk_3_18 ../bin/com.github.mubitosh.qrshare