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

	// amazon.GetItems([]string{
	// 	"B01KMXS2TK",
	// 	"B00KAPFOEM",
	// })

	url := chi.URLParam(r, "url")

	// Validate item URL
	match, _ := regexp.MatchString("^[A-Z0-9]{10}$", url)
	if !match {
		returnTemplate(w, "404", nil)
		return
	}

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
	return
}

type itemVars struct {
	Item item
}
