package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayhilel/say-hi/content"
)

func LandingHandler(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

func ProjectsHandler(c *fiber.Ctx) error {

	p := content.Project{Name: "PJX", Url: "X.com", Description: "X"}

	switch c.Params("num") {
	case "1":
		p.Name = "PJ1"
		p.Description = "A simple app designed to be simple."
		p.Url = "xyz.com"

	case "2":
		p.Name = "PJ2"
		p.Description = "A complex app designed with rust"
		p.Url = "abc.com"

	}

	return c.Render("project", fiber.Map{
		"Name":        p.Name,
		"Description": p.Description,
		"Url":         p.Url,
	})
}
