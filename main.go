package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	amaz "github.com/Jleagle/canihave/amazon"
	"github.com/Jleagle/canihave/location"
	"github.com/Jleagle/canihave/logger"
	"github.com/Jleagle/canihave/scraper"
	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
)

const (
	SEARCH     string = "search"
	CATEGORIES string = "cats"
	INFO       string = "info"
	ITEM       string = "item"
)

var regions map[string]string

func main() {

	logger.Notice()

	amaz.RateLimit = time.Tick(time.Millisecond * 1200)

	regions = map[string]string{
		location.US: "United States",
		location.UK: "United Kingdom",
		//location.DE: "Deutschland",
		//location.FR: "France",
		//location.JP: "Japan",
		//location.CA: "Canada",
		//location.CN: "China",
		//location.IT: "Italia",
		//location.ES: "España",
		//location.IN: "India",
		//location.BR: "Brazil",
		//location.MX: "Mexico",
	}

	scrape := flag.Bool("scrape", false, "Grab new items from websites")
	social := flag.Bool("social", false, "Add items to social media")
	flag.Parse()
	if *scrape {
		scraper.ScrapeHandler(*social)
		return
	}

	r := chi.NewRouter()

	r.Get("/", searchHandler)
	r.Post("/", searchHandler)
	r.Get("/info", infoHandler)
	r.Get("/sitemap.xml", siteMapHandler)
	r.Get("/categories", categoriesHandler)
	r.Get("/{id}", itemHandler)
	r.Get("/{id}/{slug}", itemHandler)

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
	always := []string{folder + "/templates/header.html", folder + "/templates/footer.html", folder + "/templates/" + page + ".html"}

	t, err := template.New("t").Funcs(getTemplateFuncMap()).ParseFiles(always...)
	if err != nil {
		panic(err)
	}

	// Write a respone
	err = t.ExecuteTemplate(w, page, pageData)
	if err != nil {
		panic(err)
	}
}

func getTemplateFuncMap() map[string]interface{} {
	return template.FuncMap{
		"avail": func(name string, data interface{}) bool {
			v := reflect.ValueOf(data)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			if v.Kind() != reflect.Struct {
				return false
			}
			return v.FieldByName(name).IsValid()
		},
		"inc": func(i int) int { return i + 1 },
		"dec": func(i int) int { return i - 1 },
		"cmp": func(i interface{}, j interface{}) bool { return i == j },
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

type errorVars struct {
	HTTPCode int
	Message  string
}
