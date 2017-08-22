package main

import (
	"fmt"
	"net/http"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {

	db, _ := connectToSQL()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM items ORDER BY date_created DESC")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

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

type searchVars struct {
	Items []item
}
