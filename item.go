package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func itemHandler(w http.ResponseWriter, r *http.Request) {

	url := chi.URLParam(r, "url")

	db, _ := connectToSQL()
	defer db.Close()

	// Get an item
	item := item{}
	err := db.QueryRow("SELECT * FROM items WHERE id = ?", url).Scan(&item.ID, &item.DateCreated, &item.DateUpdated, &item.TimesAdded, &item.Name, &item.Desc, &item.Source)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// Handle some errors
	switch {
	case err == sql.ErrNoRows:
		fmt.Fprintf(w, "No such item")
	case err != nil:
		log.Fatal(err)
	}

	// Return template
	itemVars := itemVars{}
	itemVars.Item = item

	returnTemplate(w, "item", itemVars)
}

type itemVars struct {
	Item item
}
