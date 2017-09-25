package main

import (
	"net/http"

	"github.com/Jleagle/canihave/location"
	"github.com/Jleagle/canihave/logger"
	"github.com/Jleagle/canihave/models"
	"github.com/Jleagle/canihave/store"
	"github.com/Masterminds/squirrel"
)

func categoriesHandler(w http.ResponseWriter, r *http.Request) {

	location.DetectLanguageChange(w, r)

	builder := squirrel.Select("nodeName AS category, count(nodeName) AS count").From("items").GroupBy("nodeName").OrderBy("count DESC").Where("type = ?", models.TYPE_SCRAPE)
	rows := store.Query(builder)
	defer rows.Close()

	results := []category{}
	item := category{}
	for rows.Next() {
		err := rows.Scan(&item.Category, &item.Count)
		if err != nil {
			logger.Err("Can't scan category", err)
		}

		results = append(results, item)
	}

	vars := categoriesVars{}
	vars.Flag = location.GetAmazonRegion(w, r)
	vars.Flags = location.GetRegions()
	vars.Items = results
	vars.Path = r.URL.Path
	vars.WebPage = PAGE_CATEGORIES

	returnTemplate(w, "categories", vars)
}

type category struct {
	Category string
	Count    string
}

type categoriesVars struct {
	Path    string
	Name    string
	Size    int
	Flag    string
	Flags   map[string]string
	Items   []category
	WebPage string
}
