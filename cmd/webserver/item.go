package main

import (
	"regexp"

	"github.com/Jleagle/canihave/pkg/helpers"
	"github.com/Jleagle/canihave/pkg/location"
	"github.com/Jleagle/canihave/pkg/logger"
	"github.com/Jleagle/canihave/pkg/models"
	"github.com/Jleagle/canihave/pkg/mysql"
	"github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
)

func itemHandler(c *fiber.Ctx) error {

	id := c.Params("id")

	// Validate item ID
	match, err := regexp.MatchString("^[A-Z0-9]{10}$", id)
	if err != nil {
		logger.Err(err.Error(), err)
	}
	if !match {
		returnError(c, errorVars{HTTPCode: 404, Message: "Invalid Item ID"})
		return
	}

	// Get item details
	item, err := mysql.GetWithExtras(id, location.GetRegion(c), models.typeManual, main.SOURCE_Manual)
	if err != nil {
		logger.Err(err.Error(), err)

		returnError(c, errorVars{HTTPCode: 404, Message: "Can't find item"})
		return
	}

	if !helpers.IsBot(c.Get(fiber.HeaderUserAgent)) {
		go func() {
			_, err := incrementHits(item.ID)
			if err != nil {

			}
		}()
	}

	// Return template
	vars := itemVars{}
	vars.Item = item
	vars.Javascript = []string{"//platform.twitter.com/widgets.js"}
	vars.Flag = location.GetRegion(c)
	vars.Flags = location.GetRegions()
	vars.Path = c.Path()
	vars.Similar = item.GetRelated(models.typeSimilar)

	returnTemplate(c, "item", vars)
	return nil
}

func incrementHits(id string) (success bool, err error) {

	builder := squirrel.Update("items").Set("hits", squirrel.Expr("hits + 1")).Where("id = ?", id)
	err = mysql.Update(builder)
	if err == nil {
		return true, err
	}

	logger.Err(err.Error(), err)
	return false, err
}

type itemVars struct {
	commonTemplateVars
	Item    mysql.Item
	Similar []mysql.Item
}
