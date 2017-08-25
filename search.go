package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Jleagle/canihave/helpers"
	cache "github.com/patrickmn/go-cache"
)

var perPage = 12
var maxPage = 10

func searchHandler(w http.ResponseWriter, r *http.Request) {

	c := cache.New(5*time.Minute, 10*time.Minute)

	foo, found := c.Get("foo")
	if found {
		//fmt.Println(foo)
	}

	query := r.URL.Query()
	page := query.Get("page")
	pageInt, _ := strconv.Atoi(page)
	limit := pageInt * perPage

	options := queryOptions{}
	options.limit = strconv.Itoa(limit)
	options.page = strconv.Itoa(helpers.Min([]int{pageInt, maxPage}))
	options.search = query.Get("search")

	// Return template
	vars := searchVars{}
	vars.Items = handleQuery(options)
	vars.Page = options.page

	returnTemplate(w, "search", vars)
}

func ajaxHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	pageInt, _ := strconv.Atoi(query.Get("page"))

	options := queryOptions{}
	options.limit = strconv.Itoa(perPage)
	options.page = strconv.Itoa(helpers.Min([]int{pageInt, maxPage}))
	options.search = query.Get("search")

	// Return template
	vars := searchVars{}
	vars.Items = handleQuery(options)
	vars.Page = options.page

	returnTemplate(w, "search_ajax", vars)
}

func handleQuery(options queryOptions) []item {

	// Connect to SQL
	db, _ := connectToSQL()
	defer db.Close()

	// Run the query
	rows, err := db.Query("SELECT * FROM items ORDER BY date_created DESC LIMIT 12")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// Convert to types
	results := []item{}
	item := item{}
	for rows.Next() {
		rows.Scan(&item.ID, &item.DateCreated, &item.DateUpdated, &item.TimesAdded, &item.Name, &item.Desc, &item.Source)
		results = append(results, item)
	}

	return results
}

type searchVars struct {
	Items []item
	Page  string
}

type queryOptions struct {
	limit    string
	page     string
	search   string
	category string
}
