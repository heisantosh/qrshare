#!/bin/bash

appid=com.github.mubitosh.qrshare

# Build the app
gb build -tags gtk_3_18 all
mv bin/qrshare-gtk_3_18 bin/${appid}

# Generate translations for the app
for po in $(ls po/*.po)
do
    name=$(basename $po)
    lang=${name%.*}
    mo=po/locale/${lang}/LC_MESSAGES/${appid}.mo
    mkdir -p po/locale/${lang}/LC_MESSAGES/
    msgfmt -c -v -o $mo $po
done
