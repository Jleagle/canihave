package main

import (
	"fmt"
	"net/http"

	"github.com/Jleagle/canihave/models"
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

	query := squirrel.Select("*").From("items").OrderBy("dateCreated DESC").Limit(1000)
	rows := store.Query(query)

	i := models.Item{}
	for rows.Next() {
		err := rows.Scan(&i.ID, &i.DateCreated, &i.DateUpdated, &i.DateScanned, &i.Name, &i.Link, &i.Source, &i.SalesRank, &i.Photo, &i.Node, &i.NodeName, &i.Price, &i.Region, &i.Hits, &i.Status, &i.Type, &i.CompanyName)
		if err != nil {
			fmt.Println(err)
		}

		sm.Add(stm.URL{
			"loc":        "/" + i.ID + "/" + slugify.Marshal(i.Name, true),
			"changefreq": "daily",
			"mobile":     true,
			"news": stm.URL{
				"publication": stm.URL{
					"name":     "Canihave.one/",
					"language": i.Region,
				},
				"title":            i.Name,
				"publication_date": i.DateCreated,
				//"access":           "Subscription",
				//"genres":           "PressRelease",
				//"keywords":         "my article, articles about myself",
			},
		})
	}

	w.Write(sm.XMLContent())
	return
}
