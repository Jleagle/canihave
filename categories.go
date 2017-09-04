package main

import (
	"net/http"
	"github.com/Jleagle/canihave/location"
	"github.com/Jleagle/canihave/store"
	"fmt"
)

func categoriesHandler(w http.ResponseWriter, r *http.Request) {

	conn := store.GetMysqlConnection()

	rows, err := conn.Query("SELECT productGroup AS category, count(productGroup) AS count FROM items GROUP BY productGroup ORDER BY count DESC")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	results := []category{}
	item := category{}
	for rows.Next() {
		err := rows.Scan(&item.Category, &item.Count)
		if err != nil {
			fmt.Println(err)
		}

		results = append(results, item)
	}

	vars := categoriesVars{}
	vars.Flag = location.GetAmazonRegion(w, r)
	vars.Flags = regions
	vars.Items = results

	returnTemplate(w, "categories", vars)
}

type category struct {
	Category string
	Count    string
}

type categoriesVars struct {
	Name  string
	Size  int
	Flag  string
	Flags map[string]string
	Items []category
}
