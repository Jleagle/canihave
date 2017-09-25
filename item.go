package main

import (
	"net/http"
	"regexp"

	"github.com/Jleagle/canihave/bots"
	"github.com/Jleagle/canihave/location"
	"github.com/Jleagle/canihave/logger"
	"github.com/Jleagle/canihave/models"
	"github.com/Jleagle/canihave/scraper"
	"github.com/Jleagle/canihave/store"
	"github.com/Masterminds/squirrel"
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
	item, err := models.GetWithExtras(id, location.GetAmazonRegion(w, r), models.TYPE_MANUAL, scraper.SOURCE_Manual)
	if err != nil {
		logger.Err("Can't get with extras", err)

		returnTemplate(w, "error", errorVars{HTTPCode: 404, Message: "Can't find item"})
		return
	}

	if !bots.IsBot(r.UserAgent()) {
		go incrementHits(item.ID)
	}

	// Return template
	vars := itemVars{}
	vars.Item = item
	vars.Javascript = []string{"//platform.twitter.com/widgets.js"}
	vars.Flag = location.GetAmazonRegion(w, r)
	vars.Flags = location.GetRegions()
	vars.Path = r.URL.Path
	vars.WebPage = PAGE_ITEM
	vars.Similar = item.GetRelated(models.TYPE_SIMILAR)

	returnTemplate(w, "item", vars)
	return
}

func incrementHits(id string) (success bool, err error) {

	builder := squirrel.Update("items").Set("hits", squirrel.Expr("hits + 1")).Where("id = ?", id)
	err = store.Update(builder)
	if err == nil {
		return true, err
	}
	logger.Err("Cant increment hits query", err)
	return false, err
}

type itemVars struct {
	Item       models.Item
	Javascript []string
	WebPage    string
	Flag       string
	Flags      map[string]string
	Path       string
	Similar    []models.Item
}
