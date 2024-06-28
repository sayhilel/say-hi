package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	projects "github.com/sayhilel/say-hi/internal"
	"github.com/sayhilel/say-hi/internal/handlers"
	"log"
)

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Static("/", "./public")

	app.Get("/", handlers.LandingHandler)
	app.Get("/about-me", handlers.ViewAboutMe)
	app.Get("/about-me", handlers.ViewAboutMe)
	app.Get("/contact-me", handlers.ViewContactMe)

	app.Post("/command", func(c *fiber.Ctx) error {
		command := c.FormValue("command")
		execFunc := handlers.CommandHandler[command]

		if execFunc == nil {
			return handlers.HandleInvalid(c)
		}
		return execFunc(c)
	})

	ps := projects.InitProjects()

	app.Get("/projects/:index", func(c *fiber.Ctx) error {
		return handlers.HandleProjects(ps, c)
	})

	log.Fatal(app.Listen(":3000"))
}
