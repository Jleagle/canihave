package main

import (
	"net/http"
	"regexp"

	"github.com/Jleagle/canihave/bots"
	"github.com/Jleagle/canihave/location"
	"github.com/Jleagle/canihave/models"
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

	// Get item details
	item := models.Item{}
	item.ID = id
	item.Get()

	if item.Link == "" {
		returnTemplate(w, "error", errorVars{HTTPCode: 404, Message: "Can't find item"})
		return
	}

	if !bots.IsBot(r.UserAgent()) {
		item.IncrementHits()
	}

	// Return template
	vars := itemVars{}
	vars.Item = item
	vars.Javascript = []string{"//platform.twitter.com/widgets.js"}
	vars.Flag = location.GetAmazonRegion(w, r)
	vars.Flags = regions
	vars.Path = r.URL.Path
	vars.WebPage = ITEM

	returnTemplate(w, "item", vars)
	return
}

type itemVars struct {
	Item       models.Item
	Javascript []string
	WebPage    string
	Flag       string
	Flags      map[string]string
	Path       string
}
