package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/Jleagle/canihave/location"
	"github.com/Jleagle/canihave/models"
	"github.com/Jleagle/canihave/store"
	"github.com/Masterminds/squirrel"
)

var perPage int = 94

func searchHandler(w http.ResponseWriter, r *http.Request) {

	location.DetectLanguageChange(w, r)

	// Get data
	params := r.URL.Query()

	region := location.GetAmazonRegion(w, r)
	search := params.Get("search")
	category := params.Get("cat")

	pageLimit := getPageLimit(search, category, region)
	pageLimit = int(math.Max(float64(pageLimit), 1))

	page := params.Get("page")
	if page == "" {
		page = "1"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		log.Fatal("Error converting string to int")
	}

	if pageInt < 1 {
		pageInt = 1
	} else if pageInt > pageLimit {
		pageInt = pageLimit
	}

	// Make cache key
	//cacheMd5 := md5.Sum([]byte("p" + page + ":c" + category + ":s" + search))
	//cacheKey := string(cacheMd5)

	// Return template
	vars := searchVars{}
	vars.Items = getResults(search, category, region, pageInt)
	vars.Search = search
	vars.Search64 = base64.StdEncoding.EncodeToString([]byte(search))
	vars.Javascript = []string{"//platform.twitter.com/widgets.js"}
	vars.Flag = region
	vars.Flags = location.GetRegions()
	vars.Page = pageInt
	vars.PageLimit = pageLimit
	vars.Path = r.URL.Path
	vars.WebPage = PAGE_SEARCH

	returnTemplate(w, "search", vars)
}

func getResults(search string, category string, region string, page int) []models.Item {

	offset := uint64((page - 1) * perPage)

	query := squirrel.Select("*").From("items").OrderBy("region = '" + region + "' DESC, dateUpdated DESC").Limit(uint64(perPage)).Offset(offset)
	query = filter(query, search, category, region)

	rows := store.Query(query)
	defer rows.Close()

	// Convert to types
	results := []models.Item{}
	i := models.Item{}
	for rows.Next() {
		err := rows.Scan(&i.ID, &i.DateCreated, &i.DateUpdated, &i.DateScanned, &i.Name, &i.Link, &i.Source, &i.SalesRank, &i.Photo, &i.Node, &i.NodeName, &i.Price, &i.Region, &i.Hits, &i.Status, &i.Type, &i.CompanyName)
		if err != nil {
			fmt.Println(err)
		}

		results = append(results, i)
	}

	return results
}

func getPageLimit(search string, category string, region string) int {

	query := squirrel.Select("count(id) as count").From("items")
	query = filter(query, search, category, region)

	var count int
	err := store.QueryRow(query).Scan(&count)
	if err != nil {
		fmt.Println(err)
	}

	d := float64(count) / float64(perPage)
	return int(math.Ceil(d))
}

func filter(query squirrel.SelectBuilder, search string, category string, region string) squirrel.SelectBuilder {

	query = query.Where("type = ?", "scrape")

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if category != "" {
		query = query.Where("cat = ?", category)
	}

	return query
}

type searchVars struct {
	Path       string
	Items      []models.Item
	Page       int
	PageLimit  int
	Category   string
	Search     string
	Search64   string
	Javascript []string
	Flag       string
	Flags      map[string]string
	WebPage    string
}
