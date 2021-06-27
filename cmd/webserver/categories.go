package main

import (
	"github.com/Jleagle/canihave/pkg/location"
	"github.com/Jleagle/canihave/pkg/models"
	"github.com/Jleagle/canihave/pkg/mysql"
	"github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/stvp/rollbar"
)

func categoriesHandler(c *fiber.Ctx) error {

	builder := squirrel.Select("nodeName AS category, count(nodeName) AS count").From("items").GroupBy("nodeName").OrderBy("count DESC").Where("type = ?", models.typeScraper)
	rows := mysql.Query(builder)

	//goland:noinspection GoUnhandledErrorResult
	defer rows.Close()

	var results []mysql.Category
	item := mysql.Category{}
	for rows.Next() {
		err := rows.Scan(&item.Category, &item.Count)
		if err != nil {
			rollbar.Error("error", err)
		}

		results = append(results, item)
	}

	vars := categoriesVars{}
	vars.Flag = location.GetRegion(c)
	vars.Flags = location.GetRegions()
	vars.Items = results
	vars.Path = c.Path()

	returnTemplate(c, "categories", vars)
}

type categoriesVars struct {
	commonTemplateVars
	Items []mysql.Category
}
