package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayhilel/say-hi/content"
)

func LandingHandler(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Name": "PJ1",
	})
}

func ProjectsHandler(c *fiber.Ctx) error {
	p := content.Project{Name: "PJ1", Url: "PJ1.com", Description: "It does it all"}

	return c.Render("index", fiber.Map{
		"Name": p.Name + " -> " + p.Description + " -> " + p.Url,
	})
}

func CollapseHandler(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Name": "PJ1",
	})
}
