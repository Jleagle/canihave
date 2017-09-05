package main

import (
	"net/http"
	"github.com/Jleagle/canihave/location"
)

func infoHandler(w http.ResponseWriter, r *http.Request) {

	location.ChangeLanguage(w, r)

	vars := infoVars{}
	vars.Javascript = []string{"//platform.twitter.com/widgets.js"}
	vars.Flag = location.GetAmazonRegion(w, r)
	vars.Flags = regions
	vars.Path = r.URL.Path

	returnTemplate(w, "info", vars)
	return
}

type infoVars struct {
	Path       string
	Javascript []string
	Flag       string
	Flags      map[string]string
}
