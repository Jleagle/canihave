package main

import (
	"github.com/Jleagle/canihave/pkg/location"
	"github.com/Jleagle/canihave/pkg/mysql"
	"github.com/gofiber/fiber/v2"
)

func categoriesHandler(c *fiber.Ctx) error {

	rows := mysql.GetCategoryCounts(mysql.TypeScraper)

	vars := categoriesVars{}
	vars.Flag = location.GetRegion(c)
	vars.Flags = location.GetRegions()
	vars.Items = results
	vars.Path = c.Path()


}

type categoriesVars struct {
	commonTemplateVars
	Items []mysql.Category
}
