package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func HandleCommands(c *fiber.Ctx) error {
	switch c.FormValue("command") {

	default:
		return HandleInvalid(c)
	case "whoami":
		return ViewAboutMe(c)
	case "showcase":
		return ViewProjects(c)
	case "ping":
		return ViewContactMe(c)

	}
}

func HandleInvalid(c *fiber.Ctx) error {
	return c.Render("layouts/invalid", fiber.Map{})
}

func ViewLanding(c *fiber.Ctx) error {
	return c.Render("layouts/landing", fiber.Map{})
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
