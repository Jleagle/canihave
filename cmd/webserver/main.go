package main

import (
	"embed"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Jleagle/canihave/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html"
	"github.com/ngs/go-amazon-product-advertising-api/amazon"
)

//go:embed views/*
var viewsfs embed.FS

func main() {

	engine := html.NewFileSystem(http.FS(viewsfs), ".gohtml")
	engine.AddFunc("inc", func(i int) int { return i + 1 })
	engine.AddFunc("dec", func(i int) int { return i - 1 })
	engine.AddFunc("cmp", func(i interface{}, j interface{}) bool { return i == j })
	engine.AddFunc("startsWith", func(string string, prefix string) bool { return strings.HasPrefix(string, prefix) })

	app := fiber.New(fiber.Config{Views: engine})

	// Middleware
	if config.Environment == config.EnvProd {
		app.Use(cache.New(cache.Config{Expiration: time.Minute, KeyGenerator: func(c *fiber.Ctx) string { return c.OriginalURL() }}))
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

	return returnTemplate(c, "home", fiber.Map{
		"Title": "Hello, World!",
	})
}

func sitemapHandler(c *fiber.Ctx) error {

	return returnTemplate(c, "sitemap", fiber.Map{
		"Title": "Hello, World!",
	})
}

func returnTemplate(c *fiber.Ctx, page string, pageData fiber.Map) error {

	c.Response().Header.Set("x", "x")

	return c.Render(page, pageData)
}

func returnError(c *fiber.Ctx) error {

	return returnTemplate(c, "error", fiber.Map{
		"Title": "Hello, World!",
	})
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
