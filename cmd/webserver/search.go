package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"log"
	"math"
	"strconv"

	"github.com/Jleagle/canihave/pkg/location"
	"github.com/Jleagle/canihave/pkg/logger"
	"github.com/Jleagle/canihave/pkg/mysql"
	"github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var perPage = 64

func searchHandler(c *fiber.Ctx) error {

	// Get data
	region := location.GetRegion(c)
	search := c.Query("search")
	category := c.Query("category")

	pageLimit, err := getPageLimit(search, category, region)
	if err != nil {
		logger.Logger.Error("Can't count all items", zap.Error(err))
		returnError(c, errorVars{HTTPCode: 503})
		return
	}
	pageLimit = int(math.Max(float64(pageLimit), 1))

	page := c.Query("page")
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
	// cacheMd5 := md5.Sum([]byte("p" + page + ":c" + category + ":s" + search))
	// cacheKey := string(cacheMd5)

	// Return template
	vars := searchVars{}
	vars.Items = getResults(search, category, region, pageInt)
	vars.Search = search
	vars.Search64 = base64.StdEncoding.EncodeToString([]byte(search))
	vars.Javascript = []string{"//platform.twitter.com/widgets.js", "assets/search.js"}
	vars.Flag = region
	vars.Flags = location.GetRegions()
	vars.Page = pageInt
	vars.PageLimit = pageLimit
	vars.Path = c.Path()

	// Hidden search form
	vars.Category = category
	vars.Sort = c.Query("sort")

	returnTemplate(c, "search", vars)
}

func getResults(search string, category string, region string, page int) []mysql.Item {

	offset := uint64((page - 1) * perPage)

	query := squirrel.Select("*").From("items").Limit(uint64(perPage)).Offset(offset)

	if search != "" {
		query = query.OrderBy("salesRank ASC")
	} else {
		query = query.OrderBy("region = '" + region + "' DESC, dateCreated DESC")
	}

	query = filter(query, search, category, region)

	rows := mysql.Query(query)

	//goland:noinspection GoUnhandledErrorResult
	defer rows.Close()

	// Convert to types
	var results []mysql.Item
	i := mysql.Item{}
	for rows.Next() {
		err := rows.Scan(&i.ID, &i.DateCreated, &i.DateUpdated, &i.DateScanned, &i.Name, &i.Link, &i.Source, &i.SalesRank, &i.Photo, &i.Node, &i.NodeName, &i.Price, &i.Region, &i.Hits, &i.Type, &i.CompanyName)
		if err != nil {
			logger.Logger.Error("Can't scan search result", zap.Error(err))
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

	mcItem, err := memcache.GetMemcacheItem(mcKey)

	if err == memcache.ErrCacheMiss {

		// Calculate from MySQL
		query := squirrel.Select("count(id) as count").From("items")
		query = filter(query, search, category, region)

		var count int
		err := mysql.QueryRow(query).Scan(&count)
		ret := math.Ceil(float64(count) / float64(perPage))

		memcache.SetMemcacheItem(mcKey, float64bytes(ret))
		return int(ret), err

	} else if err == nil {

		// Get from memcache
		return int(binary.BigEndian.Uint64(mcItem.Value)), err
	}

	logger.Logger.Error("Getting page limit from memcache", zap.Error(err))
	return 1, err
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

	query = query.Where("type = ?", models.typeScraper).Where("photo != ''")

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if category != "" {
		// query = query.Where("cat = ?", category) // todo
	}

	return query
}

type searchVars struct {
	commonTemplateVars
	Items     []mysql.Item
	Page      int
	PageLimit int
	Category  string
	Search    string
	Search64  string
}
