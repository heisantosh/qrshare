#!/bin/bash

glib-compile-resources --generate-source *.gresource.xml
glib-compile-resources --generate-header *.gresource.xml
mv pokicons.h pokicons.c ../

