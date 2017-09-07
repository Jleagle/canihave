package main

import (
	"net/http"
	"github.com/Jleagle/canihave/location"
)

func infoHandler(w http.ResponseWriter, r *http.Request) {

	location.DetectLanguageChange(w, r)

	vars := infoVars{}
	vars.Javascript = []string{"//platform.twitter.com/widgets.js"}
	vars.Flag = location.GetAmazonRegion(w, r)
	vars.Flags = regions
	vars.Path = r.URL.Path
	vars.WebPage = INFO

	returnTemplate(w, "info", vars)
	return
}

type infoVars struct {
	WebPage    string
	Path       string
	Javascript []string
	Flag       string
	Flags      map[string]string
}
