package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/Jleagle/canihave/pkg/config"
	"github.com/Jleagle/canihave/pkg/location"
	"github.com/Jleagle/canihave/pkg/logger"
	"github.com/gofiber/fiber/v2"
	fiverCache "github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ngs/go-amazon-product-advertising-api/amazon"
)

func main() {

	app := fiber.New()

	// Middleware
	if config.Environment == config.EnvProd {
		app.Use(fiverCache.New(fiverCache.Config{Expiration: time.Minute, KeyGenerator: func(c *fiber.Ctx) string { return c.OriginalURL() }}))
	}
	app.Use(compress.New(compress.Config{Level: compress.LevelBestSpeed}))
	app.Use(cors.New(cors.Config{AllowOrigins: "*", AllowMethods: "GET"}))

	// Routes
	app.Get("/", homeHandler)
	app.Get("/sitemap.xml", sitemapHandler)
	app.Get("/info", infoHandler)
	app.Get("/categories", categoriesHandler)
	app.Get("/:id", itemHandler)
	app.Get("/:id/:slug", itemHandler)

	// Serve
	err := app.Listen("0.0.0.0:" + config.Port)
	fmt.Println(err)
}

func homeHandler(c *fiber.Ctx) error {
	return c.SendString("OK")
}

func sitemapHandler(c *fiber.Ctx) error {
	return c.SendString("OK")
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

func returnTemplate(c *fiber.Ctx, page string, pageData interface{}) {

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
		logger.Logger.Error(err.Error())
	}

	// Write a respone
	err = t.ExecuteTemplate(w, page, pageData)
	if err != nil {
		logger.Logger.Error(err.Error())
	}
}

func returnError(c *fiber.Ctx, vars errorVars) {

	logger.Logger.Info("Showing error template")

	vars.Flag = location.GetRegion(c)
	vars.Flags = location.GetRegions()
	vars.Path = c.Path()

	vars.Category = c.Query("category")
	vars.Search = c.Query("search")
	vars.Sort = c.Query("sort")

	returnTemplate(c, "error", vars)
}

type commonTemplateVars struct {
	Flag       amazon.Region
	Flags      []amazon.Region
	Path       string
	Javascript []string

	// For hidden search form fields
	Category string
	Search   string
	Sort     string
}

type errorVars struct {
	commonTemplateVars
	HTTPCode int
	Message  string
}
