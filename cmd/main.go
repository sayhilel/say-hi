package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/sayhilel/say-hi/content"
	"github.com/sayhilel/say-hi/handlers"
)

func main() {
	engine := html.New("./views", ".html")
	
	var pl[] content.Project
	
	pl = append(pl, content.Project{Name: "TEST", Url: "Test.com", Description: "Is a test"})


	pl = append(pl, content.Project{Name: "TEST 2", Url: "Test2.com", Description: "Is a test2"})

	app := fiber.New(fiber.Config{Views: engine})
	app.Static("/", "./public")

	app.Get("/", handlers.LandingHandler)
	app.Post("/project/:num", func(c* fiber.Ctx) error {
		return handlers.ProjectsHandler(c, pl)
		})

	log.Fatal(app.Listen(":3000"))
}
