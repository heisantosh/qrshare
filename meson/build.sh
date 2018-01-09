#!/bin/bash

appid=com.github.mubitosh.qrshare

# GTK+ version
version=$(pkg-config --modversion gtk+-3.0 | cut -d '.' -f 1,2 | tr . _)

# Build the app
gb build -tags gtk_${version} all
mv bin/qrshare-gtk_${version} bin/${appid}

# Generate translations for the app
for po in $(ls po/*.po)
do
    name=$(basename $po)
    lang=${name%.*}
    mo=po/locale/${lang}/LC_MESSAGES/${appid}.mo
    mkdir -p po/locale/${lang}/LC_MESSAGES/
    msgfmt -c -v -o $mo $po
done
