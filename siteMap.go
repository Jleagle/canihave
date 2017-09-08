package main

import (
	"fmt"
	"net/http"

	"github.com/Jleagle/canihave/store"
	"github.com/Masterminds/squirrel"
	"github.com/ikeikeikeike/go-sitemap-generator/stm"
	"github.com/metal3d/go-slugify"
)

func siteMapHandler(w http.ResponseWriter, r *http.Request) {

	sm := stm.NewSitemap()
	sm.SetDefaultHost("https://canihave.one/")
	sm.SetCompress(true)
	sm.Create()

	// todo, cache
	query := squirrel.Select("*").From("items").OrderBy("dateCreated DESC").Limit(1000)
	rows := store.QueryRows(query)
	defer rows.Close()

	var ID, DateCreated, DateUpdated, Name, Link, Source, SalesRank, Photo, ProductGroup, Price, Region, Hits, Status string
	for rows.Next() {
		err := rows.Scan(&ID, &DateCreated, &DateUpdated, &Name, &Link, &Source, &SalesRank, &Photo, &ProductGroup, &Price, &Region, &Hits, &Status)
		if err != nil {
			fmt.Println(err)
		}

		sm.Add(stm.URL{
			"loc":        "/" + ID + "/" + slugify.Marshal(Name, true),
			"changefreq": "daily",
			"mobile":     true,
			//"title":            Name,
			//"publication_date": DateCreated,
		})
	}

	w.Write(sm.XMLContent())
	return
}
