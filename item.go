package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/go-chi/chi"
)

func itemHandler(w http.ResponseWriter, r *http.Request) {

	// getItems([]string{
	// 	"B01KMXS2TK",
	// 	"B00KAPFOEM",
	// })

	// Get the string associated with the key "foo" from the cache

	// return

	//importItems()

	id := chi.URLParam(r, "id")

	// Validate item ID
	match, _ := regexp.MatchString("^[A-Z0-9]{10}$", id)
	if !match {
		returnTemplate(w, "error", errorVars{HTTPCode: 404, Message: "Invalid Item ID"})
		return
	}

	db := connectToSQL()

	// Get an item
	item := item{}
	error := db.QueryRow("SELECT * FROM items WHERE id = ? LIMIT 1", id).Scan(&item.ID, &item.DateCreated, &item.DateUpdated, &item.TimesAdded, &item.Name, &item.Desc, &item.Source)
	if error != nil {
		fmt.Println(error)
	}

	// Handle some errors
	switch {
	case error == sql.ErrNoRows:
		fmt.Fprintf(w, "No such item")
	case error != nil:
		log.Fatal(error)
	}

	// Return template
	itemVars := itemVars{}
	itemVars.Item = item

	returnTemplate(w, "item", itemVars)
	return
}

type itemVars struct {
	Item item
}
