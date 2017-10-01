package main

import (
	"path/filepath"

	"github.com/gosexy/gettext"
)

func initI18n() {
	gettext.SetLocale(gettext.LC_ALL, "")
	gettext.BindTextdomain(appID, filepath.Join("/usr/share", "locale"))
	gettext.BindTextdomainCodeset(appID, "UTF-8")
	gettext.Textdomain(appID)
}

// T returns the value of gettext.Gettext.
// It is a shorthand of using gettext.Gettext function.
func T(s string) string {
	return gettext.Gettext(s)
}
