package handlers

import (
	"github.com/gofiber/fiber/v2"
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
