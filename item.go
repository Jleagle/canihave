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

	db, _ := connectToSQL()
	defer db.Close()

	// Get an item
	item := item{}
	err := db.QueryRow("SELECT * FROM items WHERE id = ?", id).Scan(&item.ID, &item.DateCreated, &item.DateUpdated, &item.TimesAdded, &item.Name, &item.Desc, &item.Source)
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
	return
}

type itemVars struct {
	Item item
}
