package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"runtime"
)

func main() {

	http.HandleFunc("/", homeHandler)

	log.Fatal(http.ListenAndServe(":8083", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	// 404 page
	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	// Get current app path
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	folder := path.Dir(file)

	// Load templates needed
	t, err := template.ParseFiles(folder + "/templates/home.html")
	if err != nil {
		panic(err)
	}

	// Write a respone
	err = t.ExecuteTemplate(w, "home", nil)
	if err != nil {
		panic(err)
	}
}
