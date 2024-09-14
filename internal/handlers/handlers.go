package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sayhilel/say-hi/internal/quotes"
)

func HandleCommands(c *fiber.Ctx) error {
	switch c.FormValue("command") {

	default:
		return HandleInvalid(c)
	case "sudo":
		return secretUser(c)
	case "whoami":
		return ViewAboutMe(c)
	case "showcase":
		return ViewProjects(c)
	case "ping":
		return ViewContactMe(c)
	case "more":
		return OpenConfirmation(c)
	case "clear":
		return ClearField(c)
	}
}

// get pwned lol
func secretUser(c *fiber.Ctx) error {
	return c.Render("layouts/secret", fiber.Map{})
}

func ClearField(c *fiber.Ctx) error {
	return c.Render("layouts/prompt", fiber.Map{})
}

func HandleInvalid(c *fiber.Ctx) error {
	return c.Render("layouts/invalid", fiber.Map{})
}

func ViewLanding(c *fiber.Ctx) error {
	return c.Render("layouts/landing", quote.GetQuote())
}

func LandingHandler(c *fiber.Ctx) error {
	userAgent := c.Get("User-Agent")
	if strings.Contains(strings.ToLower(userAgent), "mobile") {
		return c.Render("layouts/mobile", fiber.Map{})
	}

	return c.Render("index", fiber.Map{
		"q": "Use (C-c) to refresh this page for a random quote.",
		"a": "Sahil Sinha",
	})
}

func OpenConfirmation(c *fiber.Ctx) error {
	return c.Render("layouts/dialog-box", fiber.Map{
		"appl": "Open resume?",
	})
}

func ViewAboutMe(c *fiber.Ctx) error {
	return c.Render("layouts/about-me", fiber.Map{})
}

func ViewContactMe(c *fiber.Ctx) error {
	return c.Render("layouts/contact-me", fiber.Map{})
}

func ViewProjects(c *fiber.Ctx) error {
	return c.Render("layouts/projects", fiber.Map{
		"shown": "hidden",
	})
}
