package main

import (
	"fmt"
	"net/http"
)

func productsHandler(w http.ResponseWriter, r *http.Request) {
	//url := chi.URLParam(r, "url")

	db, _ := connectToSQL()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM products ORDER BY date_created DESC")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// Make an array of requests for the template
	results := []product{}
	for rows.Next() {
		var id, name, desc string

		rows.Scan(&id, &name, &desc)

		product := product{}
		product.ID = id
		product.Name = name
		product.Desc = desc

		results = append(results, product)
	}

	vars := requestTemplateVars{}
	vars.Products = results

	returnTemplate(w, "products", vars)
}

// product is the database row
type product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type requestTemplateVars struct {
	Products []product
}
