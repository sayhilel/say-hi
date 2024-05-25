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

	p := content.Project{Name: "PJX", Url: "X.com", Description: "X"}

	switch c.Params("num") {
	case "1":
		p.Name = "PJ1"
		p.Description = "Desc 1"
		p.Url = "a"

	case "2":
		p.Name = "PJ2"
		p.Description = "Desc 2"
		p.Url = "b"

	}

	return c.Render("project", fiber.Map{
		"Project": p.Name + " -> " + p.Description + " -> " + p.Url,
	})
}
