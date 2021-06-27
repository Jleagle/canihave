package main

import (
	"github.com/Jleagle/canihave/pkg/location"
	"github.com/gofiber/fiber/v2"
)

func infoHandler(c *fiber.Ctx) error {

	vars := infoVars{}
	vars.Javascript = []string{"//platform.twitter.com/widgets.js"}
	vars.Flag = location.GetRegion(c)
	vars.Flags = location.GetRegions()
	vars.Path = c.Path()

	returnTemplate(c, "info", vars)

	return nil
}

type infoVars struct {
	commonTemplateVars
}
