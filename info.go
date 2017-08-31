package main

import (
	"net/http"
	"github.com/Jleagle/canihave/location"
)

func infoHandler(w http.ResponseWriter, r *http.Request) {

	vars := infoVars{}
	vars.Search = ""
	vars.Search64 = ""
	vars.Javascript = []string{"//platform.twitter.com/widgets.js"}
	vars.Flag = location.GetAmazonRegion(w, r)
	vars.Flags = regions

	returnTemplate(w, "info", vars)
	return
}

type infoVars struct {
	Search     string
	Search64   string
	Javascript []string
	Flag       string
	Flags      map[string]string
}
