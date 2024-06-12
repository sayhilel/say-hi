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
	app.Get("/about-me", handlers.ViewAboutMe)

	pl := handlers.Plist{{
		Name:        "X",
		Url:         "X.com",
		Description: "No",
	}}

	pl = append(pl, handlers.Project{Name: "Y", Url: "Y.com", Description: "Yes"})

	app.Get("/projects", handlers.ViewProjects)
	app.Get("/projects/:index", pl.SwitchProject)

	app.Get("/contact-me", handlers.ViewContactMe)
	app.Get("/resume", handlers.ViewResume)

	log.Fatal(app.Listen(":3000"))
}
