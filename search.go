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
	"database/sql"
	"math"
)

var perPage int = 94

func searchHandler(w http.ResponseWriter, r *http.Request) {

	location.ChangeLanguage(w, r)

	// Get data
	params := r.URL.Query()

	region := location.GetAmazonRegion(w, r)
	search := params.Get("search")
	category := params.Get("cat")

	pageLimit := getPageLimit(search, category, region)

	page := params.Get("page")
	if page == "" {
		page = "1"
	}
	pageInt, _ := strconv.Atoi(page)

	if pageInt < 1 {
		pageInt = 1
	} else if pageInt > pageLimit {
		pageInt = pageLimit
	}

	// Make cache key
	//cacheMd5 := md5.Sum([]byte("p" + page + ":c" + category + ":s" + search))
	//cacheKey := string(cacheMd5)

	// Return template
	vars := searchVars{}
	vars.Items = getResults(search, category, region, pageInt)
	vars.Search = search
	vars.Search64 = base64.StdEncoding.EncodeToString([]byte(search))
	vars.Javascript = []string{"//platform.twitter.com/widgets.js"}
	vars.Flag = location.GetAmazonRegion(w, r)
	vars.Flags = regions
	vars.Page = pageInt
	vars.PageLimit = pageLimit
	vars.Path = r.URL.Path

	returnTemplate(w, "search", vars)
}

func getResults(search string, category string, region string, page int) []models.Item {

	offset := uint64((page - 1) * perPage)

	query := squirrel.Select("*").From("items").OrderBy("dateCreated DESC").Limit(uint64(perPage)).Offset(offset)
	query = filter(query, search, category, region)

	rows := runQueryRows(query)
	defer rows.Close()

	// Convert to types
	results := []models.Item{}
	item := models.Item{}
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.DateCreated, &item.DateUpdated, &item.Name, &item.Link, &item.Source, &item.SalesRank, &item.Photo, &item.ProductGroup, &item.Price, &item.Region)
		if err != nil {
			fmt.Println(err)
		}

		results = append(results, item)
	}

	return results
}

func getPageLimit(search string, category string, region string) int {

	query := squirrel.Select("count(id) as count").From("items")
	query = filter(query, search, category, region)

	var count int
	err := runQueryRow(query).Scan(&count)
	if err != nil {
		fmt.Println(err)
	}

	d := float64(count) / float64(perPage)
	return int(math.Ceil(d))
}

func filter(query squirrel.SelectBuilder, search string, category string, region string) (squirrel.SelectBuilder) {

	query = query.Where("region = ?", region)

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if category != "" {
		query = query.Where("cat = ?", category)
	}

	return query
}

func runQueryRows(queryBuilder squirrel.SelectBuilder) (*sql.Rows) {

	rawSQL, args, err := queryBuilder.ToSql()
	if err != nil {
		fmt.Println(err)
	}

	// Run SQL
	rows, err := store.GetMysqlConnection().Query(rawSQL, args...)
	if err != nil {
		fmt.Println(err)
	}

	return rows
}

func runQueryRow(queryBuilder squirrel.SelectBuilder) (*sql.Row) {

	rawSQL, args, err := queryBuilder.ToSql()
	if err != nil {
		fmt.Println(err)
	}

	// Run SQL
	return store.GetMysqlConnection().QueryRow(rawSQL, args...)
}

type searchVars struct {
	Path       string
	Items      []models.Item
	Page       int
	PageLimit  int
	Category   string
	Search     string
	Search64   string
	Javascript []string
	Flag       string
	Flags      map[string]string
}
