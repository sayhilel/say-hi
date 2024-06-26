package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayhilel/say-hi/internal"
)

func LandingHandler(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

func ViewProjects(c *fiber.Ctx) error {
	return c.Render("layouts/projects", fiber.Map{})
}

func ViewAboutMe(c *fiber.Ctx) error {
	return c.Render("layouts/about-me", fiber.Map{})
}

func ViewContactMe(c *fiber.Ctx) error {
	return c.Render("layouts/contact-me", fiber.Map{})
}

func ViewResume(c *fiber.Ctx) error {
	return c.Render("layouts/resume", fiber.Map{})
}

// Project Handlers //
func SwitchProject(ps projects.Projects, c *fiber.Ctx) error {
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
