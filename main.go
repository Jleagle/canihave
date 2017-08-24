package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	r := chi.NewRouter()

	r.Get("/", searchHandler)
	r.Post("/", searchHandler)
	r.Get("/ajax", ajaxHandler)
	r.Get("/{url}", itemHandler)

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

func connectToSQL() (*sql.DB, error) {

	password := os.Getenv("SQL_PW")
	if len(password) > 0 {
		password = ":" + password
	}

	db, err := sql.Open("mysql", "root"+password+"@tcp(127.0.0.1:3306)/canihave")
	if err != nil {
		panic(err.Error())
	}

	return db, err
}

// item is the database row
type item struct {
	ID          string
	DateUpdated string
	DateCreated string
	TimesAdded  string
	Name        string
	Desc        string
	Source      string
}

func (i item) GetUKLink() string {
	return "https://www.amazon.co.uk/dp/" + i.ID + "?tag=canihaveone00-21"
}

func (i item) GetUKPixel() string {
	return "//ir-uk.amazon-adsystem.com/e/ir?t=canihaveone00-21&l=am2&o=2&a=B000J34HN4"
}

type source struct {
	Name string
}
