package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Jleagle/canihave/scraper"
	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	r := chi.NewRouter()

	r.Get("/", searchHandler)
	r.Post("/", searchHandler)
	r.Get("/info", infoHandler)
	r.Get("/scrape", scraper.ScrapeHandler)
	r.Get("/ajax", ajaxHandler)
	r.Get("/{id}", itemHandler)

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "assets")
	fileServer(r, "/assets", http.Dir(filesDir))

	log.Fatal(http.ListenAndServe(":8083", r))
}

func returnTemplate(w http.ResponseWriter, page string, pageData interface{}) {

	// Get current app path
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	folder := path.Dir(file)

	// Load templates needed
	t, err := template.ParseFiles(folder+"/templates/header.html", folder+"/templates/footer.html", folder+"/templates/card.html", folder+"/templates/"+page+".html")
	if err != nil {
		panic(err)
	}

	// Write a respone
	err = t.ExecuteTemplate(w, page, pageData)
	if err != nil {
		panic(err)
	}
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {

	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

type source struct {
	Name string
}

type errorVars struct {
	HTTPCode int
	Message  string
}
