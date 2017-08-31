package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Jleagle/canihave/models"
	"github.com/Jleagle/canihave/store"
	"github.com/Masterminds/squirrel"
	"github.com/Jleagle/canihave/location"
)

var perPage int = 94

func searchHandler(w http.ResponseWriter, r *http.Request) {

	params :=r.URL.Query()

	// Country override
	flag := r.URL.Query().Get("flag")
	if flag != "" {
		location.SetCookie(w, flag)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	// Get data
	search := params.Get("search")
	category := params.Get("cat")
	page := params.Get("page")
	if page == "" {
		page = "1"
	}

	// Make SQL
	pageInt, _ := strconv.Atoi(page)
	if pageInt < 1 {
		pageInt = 1
	}

	conn := store.GetMysqlConnection()
	query := squirrel.Select("*").From("items")
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	if category != "" {
		query = query.Where("cat = ?", category)
	}
	query = query.OrderBy("dateCreated DESC").Limit(uint64(perPage)).Offset(uint64((pageInt - 1) * perPage))

	sql, args, err := query.ToSql()
	if err != nil {
		fmt.Println(err)
	}

	// Run SQL
	rows, err := conn.Query(sql, args...)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	// Convert to types
	results := []models.Item{}
	item := models.Item{}
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.DateCreated, &item.DateUpdated, &item.Name, &item.Link, &item.Source, &item.SalesRank, &item.Photo, &item.ProductGroup, &item.Price, &item.Currency)
		if err != nil {
			fmt.Println(err)
		}

		results = append(results, item)
	}

	// Return template
	vars := searchVars{}
	vars.Items = results
	vars.Search = search
	vars.Search64 = base64.StdEncoding.EncodeToString([]byte(search))
	vars.Javascript = []string{"//platform.twitter.com/widgets.js"}
	vars.Flag = location.GetAmazonRegion(w, r)
	vars.Flags = regions
	vars.Page = page

	returnTemplate(w, "search", vars)
}

type searchVars struct {
	Items      []models.Item
	Page       string
	Category   string
	Search     string
	Search64   string
	Javascript []string
	Flag       string
	Flags      map[string]string
}
