package main

import (
	"time"

	"github.com/Jleagle/canihave/pkg/logger"
	"github.com/Jleagle/canihave/pkg/mysql"
	"github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/ikeikeikeike/go-sitemap-generator/stm"
	"go.uber.org/zap"
)

func siteMapHandler(c *fiber.Ctx) error {

	sm := stm.NewSitemap()
	sm.SetDefaultHost("https://canihave.one/")
	sm.SetCompress(true)
	sm.Create()

	query := squirrel.Select("*").From("items").OrderBy("type = '" + mysql.TypeScraper + "' DESC, dateCreated DESC").Limit(1000)
	rows := mysql.Query(query)
	defer rows.Close()

	i := mysql.Item{}
	for rows.Next() {
		err := rows.Scan(&i.ID, &i.DateCreated, &i.DateUpdated, &i.DateScanned, &i.Name, &i.Link, &i.Source, &i.SalesRank, &i.Photo, &i.Node, &i.NodeName, &i.Price, &i.Region, &i.Hits, &i.Type, &i.CompanyName)
		if err != nil {
			logger.Logger.Error("Can't scan site map result", zap.Error(err))
		}

		sm.Add(stm.URL{
			"loc":        "/" + i.ID + "/" + slug.Make(i.Name),
			"changefreq": "daily",
			"mobile":     true,
			"news": stm.URL{
				"publication": stm.URL{
					"name":     "Canihave.one/",
					"language": i.Region,
				},
				"title":            i.Name,
				"publication_date": time.Unix(i.DateCreated, 0).Format("2006-01-02 15:04:05"),
				"genres":           i.NodeName,
			},
		})
	}

	w.Write(sm.XMLContent())

	return nil
}
