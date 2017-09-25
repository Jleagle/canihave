package main

import (
	"encoding/base64"
	"encoding/binary"
	"log"
	"math"
	"net/http"
	"strconv"

	"crypto/md5"
	"encoding/hex"

	"github.com/Jleagle/canihave/location"
	"github.com/Jleagle/canihave/logger"
	"github.com/Jleagle/canihave/models"
	"github.com/Jleagle/canihave/store"
	"github.com/Masterminds/squirrel"
	"github.com/bradfitz/gomemcache/memcache"
)

var perPage int = 94

func searchHandler(w http.ResponseWriter, r *http.Request) {

	location.DetectLanguageChange(w, r)

	// Get data
	params := r.URL.Query()

	region := location.GetAmazonRegion(w, r)
	search := params.Get("search")
	category := params.Get("cat")

	pageLimit, err := getPageLimit(search, category, region)
	if err != nil {
		logger.Err("Can't count all items", err)
		returnTemplate(w, "error", errorVars{HTTPCode: 503})
		return
	}
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
		err := rows.Scan(&i.ID, &i.DateCreated, &i.DateUpdated, &i.DateScanned, &i.Name, &i.Link, &i.Source, &i.SalesRank, &i.Photo, &i.Node, &i.NodeName, &i.Price, &i.Region, &i.Hits, &i.Type, &i.CompanyName)
		if err != nil {
			logger.Err("Can't scan search result", err)
		}

		results = append(results, i)
	}

	return results
}

func getPageLimit(search string, category string, region string) (ret int, err error) {

	// Get memcache key
	md5ByteArray := md5.Sum([]byte(search))
	searchHash := hex.EncodeToString(md5ByteArray[:])
	mcKey := "total-items-count-" + category + "-" + region + "-" + searchHash

	mcItem, err := store.GetMemcacheItem(mcKey)

	if err == memcache.ErrCacheMiss {

		query := squirrel.Select("count(id) as count").From("items")
		query = filter(query, search, category, region)

		var count int
		err := store.QueryRow(query).Scan(&count)
		ret := math.Ceil(float64(count) / float64(perPage))

		store.SetMemcacheItem(mcKey, float64bytes(ret))
		return int(ret), err

	} else if err != nil {

		return int(binary.BigEndian.Uint64(mcItem.Value)), err
	}

	return 0, err
}

func float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func filter(query squirrel.SelectBuilder, search string, category string, region string) squirrel.SelectBuilder {

	//query = query.Where("type = ?", "scrape")

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
