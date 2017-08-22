package main

import "net/http"

func comingSoonHandler(w http.ResponseWriter, r *http.Request) {
	returnTemplate(w, "soon", nil)
}
