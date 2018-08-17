package main

import (
	"net/http"

	"github.com/Jleagle/canihave/links"
	"github.com/Jleagle/canihave/location"
	"github.com/Jleagle/canihave/models"
	"github.com/Jleagle/canihave/store"
	"github.com/Masterminds/squirrel"
	"github.com/stvp/rollbar"
)

func categoriesHandler(w http.ResponseWriter, r *http.Request) {

	location.DetectLanguageChange(w, r)

	builder := squirrel.Select("nodeName AS category, count(nodeName) AS count").From("items").GroupBy("nodeName").OrderBy("count DESC").Where("type = ?", models.TYPE_SCRAPE)
	rows := store.Query(builder)
	defer rows.Close()

	var results []models.Category
	item := models.Category{}
	for rows.Next() {
		err := rows.Scan(&item.Category, &item.Count)
		if err != nil {
			rollbar.Error("error", err)
		}

		results = append(results, item)
	}

	vars := categoriesVars{}
	vars.Flag = location.GetAmazonRegion(w, r)
	vars.Flags = location.GetRegions()
	vars.Items = results
	vars.Path = r.URL.Path
	vars.Links = links.GetHeaderLinks(r)

	returnTemplate(w, "categories", vars)
}

type categoriesVars struct {
	commonTemplateVars
	Items []models.Category
}
