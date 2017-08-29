package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Jleagle/canihave/helpers"
	sq "github.com/Masterminds/squirrel"
)

var perPage = 12
var maxPage = 10

func searchHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	query := r.URL.Query()
	page := query.Get("page")
	pageInt, _ := strconv.Atoi(page)
	limit := pageInt * perPage

	options := queryOptions{}
	options.limit = strconv.Itoa(limit)
	options.page = strconv.Itoa(helpers.Min([]int{pageInt, maxPage}))
	options.search = r.Form.Get("search")

	// Return template
	vars := searchVars{}
	vars.Items = handleQuery(options)
	vars.Page = options.page
	vars.Search = options.search
	vars.Javascript = []string{"/assets/search.js"}

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
	vars.Search = options.search

	returnTemplate(w, "search_ajax", vars)
}

func handleQuery(options queryOptions) []item {

	// Connect to SQL
	db := connectToSQL()

	// Make the query
	query := sq.Select("*").From("items")
	if options.search != "" {
		query = query.Where("name LIKE ?", "%"+options.search+"%")
	}
	query = query.OrderBy("dateCreated DESC").Limit(12)

	sql, args, error := query.ToSql()
	if error != nil {
		fmt.Println(error)
	}

	// Run the query
	rows, error := db.Query(sql, args...)
	if error != nil {
		fmt.Println(error)
	}
	defer rows.Close()

	// Convert to types
	results := []item{}
	item := item{}
	for rows.Next() {
		rows.Scan(&item.ID, &item.DateCreated, &item.DateUpdated, &item.Name, &item.Desc, &item.Source)
		results = append(results, item)
	}

	return results
}

type searchVars struct {
	Items      []item
	Page       string
	Search     string
	Javascript []string
}

type queryOptions struct {
	limit    string
	page     string
	search   string
	category string
}
