package main

import (
	"net/http"
	"regexp"

	"github.com/go-chi/chi"
)

func itemHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	// Validate item ID
	match, _ := regexp.MatchString("^[A-Z0-9]{10}$", id)
	if !match {
		returnTemplate(w, "error", errorVars{HTTPCode: 404, Message: "Invalid Item ID"})
		return
	}

	item := item{}
	item.ID = id
	item.get()

	// importItems()

	if item.Link != "" {
		returnTemplate(w, "error", errorVars{HTTPCode: 404, Message: "Can't find item"})
		return
	}

	// Return template
	itemVars := itemVars{}
	itemVars.Item = item

	returnTemplate(w, "item", itemVars)
	return
}

type itemVars struct {
	Item item
}
