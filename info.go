package main

import (
	"net/http"
)

func infoHandler(w http.ResponseWriter, r *http.Request) {

	vars := infoVars{}
	vars.Javascript = []string{"//platform.twitter.com/widgets.js"}

	returnTemplate(w, "info", vars)
	return
}

type infoVars struct {
	Javascript []string
}
