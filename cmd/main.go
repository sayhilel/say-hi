package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/sayhilel/say-hi/handlers"
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", handlers.LandingHandler)
	app.Post("/projects", handlers.ProjectsHandler)
	app.Post("/collapse", handlers.CollapseHandler)
	log.Fatal(app.Listen(":3000"))
}
