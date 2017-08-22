package main

import "net/http"

func productHandler(w http.ResponseWriter, r *http.Request) {
	returnTemplate(w, "soon", nil)
}
