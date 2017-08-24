package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {

	rows := handleQuery(r)

	// Make an array of requests for the template
	results := []item{}
	item := item{}
	for rows.Next() {
		rows.Scan(&item.ID, &item.DateCreated, &item.DateUpdated, &item.TimesAdded, &item.Name, &item.Desc, &item.Source)
		results = append(results, item)
	}

	// Return template
	vars := searchVars{}
	vars.Items = results

	returnTemplate(w, "search", vars)
}

func ajaxHandler(w http.ResponseWriter, r *http.Request) {

	rows := handleQuery(r)

	// Make an array of requests for the template
	results := []item{}
	item := item{}
	for rows.Next() {
		rows.Scan(&item.ID, &item.DateCreated, &item.DateUpdated, &item.TimesAdded, &item.Name, &item.Desc, &item.Source)
		results = append(results, item)
	}

	// Return template
	vars := searchVars{}
	vars.Items = results

	returnTemplate(w, "search_ajax", vars)
}

func handleQuery(r *http.Request) *sql.Rows {

	db, _ := connectToSQL()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM items ORDER BY date_created DESC LIMIT 12")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	return rows
}

type searchVars struct {
	Items []item
}
