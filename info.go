package main

import (
	"net/http"

	"github.com/Jleagle/canihave/links"
	"github.com/Jleagle/canihave/location"
)

func infoHandler(w http.ResponseWriter, r *http.Request) {

	location.DetectLanguageChange(w, r)

	vars := infoVars{}
	vars.Javascript = []string{"//platform.twitter.com/widgets.js"}
	vars.Flag = location.GetAmazonRegion(w, r)
	vars.Flags = location.GetRegions()
	vars.Path = r.URL.Path
	vars.Links = links.GetHeaderLinks(r)

	returnTemplate(w, "info", vars)
	return
}

type infoVars struct {
	commonTemplateVars
}
