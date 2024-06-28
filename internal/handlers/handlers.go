package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayhilel/say-hi/internal"
)

var CommandHandler = map[string]func(c *fiber.Ctx) error{
	"projects": ViewProjects,
	"aboutme":  ViewAboutMe,
}

func HandleInvalid(c *fiber.Ctx) error {
	return c.Render("layouts/invalid", fiber.Map{})
}

func LandingHandler(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

func ViewAboutMe(c *fiber.Ctx) error {
	return c.Render("layouts/about-me", fiber.Map{})
}

func ViewContactMe(c *fiber.Ctx) error {
	return c.Render("layouts/contact-me", fiber.Map{})
}

func ViewProjects(c *fiber.Ctx) error {
	return c.Render("layouts/projects", fiber.Map{})
}

func HandleProjects(ps projects.Projects, c *fiber.Ctx) error {
	index, err := c.ParamsInt("index", 0)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid project index")
	}

	p := ps.PL[index]

	return c.Render("layouts/project", fiber.Map{
		"Name":        p.Name,
		"Description": p.Description,
		"Url":         p.Url,
		"Image":       p.Image,
	})
}
