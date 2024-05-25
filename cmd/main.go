package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/sayhilel/say-hi/handlers"
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{Views: engine})
	app.Static("/", "./public")

	app.Get("/", handlers.LandingHandler)
	app.Post("/project/:num", handlers.ProjectsHandler)

	log.Fatal(app.Listen(":3000"))
}
