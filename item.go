package main

import (
	"net/http"
	"regexp"

	"github.com/Jleagle/canihave/models"
	"github.com/go-chi/chi"
)

func itemHandler(w http.ResponseWriter, r *http.Request) {

	//Import some items to test
	models.ImportItems()

	// Validate item ID
	id := chi.URLParam(r, "id")
	match, _ := regexp.MatchString("^[A-Z0-9]{10}$", id)
	if !match {
		returnTemplate(w, "error", errorVars{HTTPCode: 404, Message: "Invalid Item ID"})
		return
	}

	// Get item details
	item := models.Item{}
	item.ID = id
	item.Get()

	if item.Link == "" {
		returnTemplate(w, "error", errorVars{HTTPCode: 404, Message: "Can't find item"})
		return
	}

	// Return template
	vars := itemVars{}
	vars.Item = item
	vars.Javascript = []string{"//platform.twitter.com/widgets.js"}

	returnTemplate(w, "item", vars)
	return
}

type itemVars struct {
	Item       models.Item
	Javascript []string
}
