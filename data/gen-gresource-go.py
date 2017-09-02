#!/usr/bin/env python

# GResource for the iconsin pokicons directory.
# Generate gresource xml file and a go file containing paths to the resources.

import os

gresource_xml = """<?xml version="1.0" encoding="UTF-8"?>
<gresources>

GRESOURCES
</gresources>"""

pokicons_go = """package main

var pokicons = []string{PATHS
}
"""

pokicons = os.listdir("pokicons")

gresources = ""
paths = ""

for pokicon in pokicons:
	gresource = "\t<gresource prefix=\"/com/github/mubitosh/qrshare\">\n"
	gresource += "\t\t<file>pokicons/" + pokicon + "</file>\n"
	gresource += "\t</gresource>\n"
	gresources += gresource

	paths += "\n\t\"/com/github/mubitosh/qrshare/pokicons/" + pokicon + "\","

gresource_xml = gresource_xml.replace("GRESOURCES", gresources)

f = open("pokicons.gresource.xml", "w")
f.write(gresource_xml)
f.close()

pokicons_go = pokicons_go.replace("PATHS", paths)

f = open("../pokicons.go", "w")
f.write(pokicons_go)
f.close()