package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	amaz "github.com/Jleagle/canihave/amazon"
	"github.com/Jleagle/canihave/environment"
	"github.com/Jleagle/canihave/links"
	"github.com/Jleagle/canihave/location"
	"github.com/Jleagle/canihave/logger"
	"github.com/Jleagle/canihave/scraper"
	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// Setup
	location.SetRegions()
	amaz.SetRateLimit()
	environment.SetEnv()

	// CLI
	scrape := flag.Bool("scrape", false, "Grab new items from websites")
	social := flag.Bool("social", false, "Add items to social media")
	flag.Parse()
	if *scrape {
		scraper.ScrapeHandler(*social)
		return
	}

	// Routes
	r := chi.NewRouter()
	r.Get("/", searchHandler)
	r.Post("/", searchHandler)
	r.Get("/info", infoHandler)
	r.Get("/sitemap.xml", siteMapHandler)
	r.Get("/categories", categoriesHandler)
	r.Get("/{id}", itemHandler)
	r.Get("/{id}/{slug}", itemHandler)
	r.MethodFunc(http.MethodGet, "/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/assets/icons/favicon-64.png", 301)
	})

	// Assets
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "assets")
	fileServer(r, "/assets", http.Dir(filesDir))

	// Serve
	log.Fatal(http.ListenAndServe(":8080", r))
}

func returnTemplate(w http.ResponseWriter, page string, pageData interface{}) {

	// Load templates needed
	folder := os.Getenv("CANIHAVE_PATH")
	if folder == "" {
		folder = "/root"
	}

	templates := []string{
		folder + "/templates/header.html",
		folder + "/templates/footer.html",
		folder + "/templates/" + page + ".html",
	}

	t, err := template.New("t").Funcs(getTemplateFuncMap()).ParseFiles(templates...)
	if err != nil {
		logger.ErrExit(err.Error())
	}

	// Write a respone
	err = t.ExecuteTemplate(w, page, pageData)
	if err != nil {
		logger.ErrExit(err.Error())
	}
}

func returnError(w http.ResponseWriter, r *http.Request, vars errorVars) {

	logger.Info("Showing error template")

	vars.Links = links.GetHeaderLinks(r)
	vars.Flag = location.GetAmazonRegion(w, r)
	vars.Flags = location.GetRegions()
	vars.Path = r.URL.Path

	params := r.URL.Query()

	vars.Category = params.Get("category")
	vars.Search = params.Get("search")
	vars.Sort = params.Get("sort")

	returnTemplate(w, "error", vars)
}

type commonTemplateVars struct {
	Links      map[string]string
	Flag       string
	Flags      map[string]string
	Path       string
	Javascript []string

	// For hidden search form fields
	Category string
	Search   string
	Sort     string
}

func getTemplateFuncMap() map[string]interface{} {
	return template.FuncMap{
		"inc": func(i int) int { return i + 1 },
		"dec": func(i int) int { return i - 1 },
		"cmp": func(i interface{}, j interface{}) bool { return i == j },
		"startsWith": func(string string, prefix string) bool {
			return strings.HasPrefix(string, prefix)
		},
	}
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {

	if strings.ContainsAny(path, "{}*") {
		logger.ErrExit("FileServer does not permit URL parameters.")
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

type errorVars struct {
	commonTemplateVars
	HTTPCode int
	Message  string
}
