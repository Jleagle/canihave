package stm

const (
	// MaxSitemapFiles defines max sitemap links per index file
	MaxSitemapFiles = 50000
	// MaxSitemapLinks defines max links per sitemap
	MaxSitemapLinks = 50000
	// MaxSitemapImages defines max images per url
	MaxSitemapImages = 1000
	// MaxSitemapNews defines max news sitemap per index_file
	MaxSitemapNews = 1000
	// MaxSitemapFilesize defines file size for sitemap.
	MaxSitemapFilesize = 10000000 // bytes
)

const (
	// SchemaGeo exists for geo sitemap
	SchemaGeo = "http://www.google.com/geo/schemas/sitemap/1.0"
	// SchemaImage exists for image sitemap
	SchemaImage = "http://www.google.com/schemas/sitemap-image/1.1"
	// SchemaMobile exists for mobile sitemap
	SchemaMobile = "http://www.google.com/schemas/sitemap-mobile/1.0"
	// SchemaNews exists for news sitemap
	SchemaNews = "http://www.google.com/schemas/sitemap-news/0.9"
	// SchemaPagemap exists for pagemap sitemap
	SchemaPagemap = "http://www.google.com/schemas/sitemap-pagemap/1.0"
	// SchemaVideo exists for video sitemap
	SchemaVideo = "http://www.google.com/schemas/sitemap-video/1.1"
)

var (
	// IndexXMLHeader exists for create sitemap xml as a specific sitemap document.
	IndexXMLHeader = []byte(`<?xml version="1.0" encoding="UTF-8"?>
      <sitemapindex
      xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
      xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9
        http://www.sitemaps.org/schemas/sitemap/0.9/siteindex.xsd"
      xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
    >`)
	// IndexXMLFooter and IndexXMLHeader will used from user together .
	IndexXMLFooter = []byte("</sitemapindex>")
)

var (
	// XMLHeader exists for create sitemap xml as a specific sitemap document.
	XMLHeader = []byte(`<?xml version="1.0" encoding="UTF-8"?>
      <urlset
        xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
        xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9
          http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd"
        xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
        xmlns:image="` + SchemaImage + `"
        xmlns:video="` + SchemaVideo + `"
        xmlns:geo="` + SchemaGeo + `"
        xmlns:news="` + SchemaNews + `"
        xmlns:mobile="` + SchemaMobile + `"
        xmlns:pagemap="` + SchemaPagemap + `"
        xmlns:xhtml="http://www.w3.org/1999/xhtml"
    >`)
	// XMLFooter and XMLHeader will used from user together .
	XMLFooter = []byte("</urlset>")
)
