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
	app.Get("/contact-me", handlers.ViewContactMe)
	app.Post("/command", handlers.HandleCommands)

	ps := projects.InitProjects()
	app.Get("/projects/:index", func(c *fiber.Ctx) error {
		return ps.HandleProjects(c)
	})

	log.Fatal(app.Listen(":3000"))
}
