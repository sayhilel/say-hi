package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	projects "github.com/sayhilel/say-hi/internal"
	"github.com/sayhilel/say-hi/internal/handlers"
)

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Static("/", "./public")

	app.Get("/", handlers.LandingHandler)
	app.Get("/about-me", handlers.ViewAboutMe)

	ps := projects.InitProjects()

	print(ps.PL[0].Description)

	app.Get("/projects", handlers.ViewProjects)
	app.Get("/projects/:index", func(c *fiber.Ctx) error {
		return handlers.SwitchProject(ps, c)
	})

	app.Get("/contact-me", handlers.ViewContactMe)
	app.Get("/resume", handlers.ViewResume)

	log.Fatal(app.Listen(":3000"))
}
