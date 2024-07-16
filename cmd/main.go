package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/template/html/v2"
	"github.com/sayhilel/say-hi/internal/handlers"
	"github.com/sayhilel/say-hi/internal/projects"
	"log"
	"os"
)

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Static("/", "./public")

	app.Use(favicon.New(favicon.Config{
		File: "public/css/images/favicon.ico",
		URL:  "/favicon.ico"}),
		cors.New(cors.Config{
			AllowHeaders: "Access-Control-Allow-Origin",
		}),
	)
	app.Get("/robots.txt", func(c *fiber.Ctx) error {
		return c.SendFile("./robots.txt")
	})

	app.Get("/sitemap.xml", func(c *fiber.Ctx) error {
		return c.SendFile("./sitemap.xml")
	})

	app.Get("/", handlers.LandingHandler)
	app.Get("/landing", handlers.ViewLanding)
	app.Post("/command", handlers.HandleCommands)

	ps := projects.InitProjects()
	app.Get("/projects/:index", func(c *fiber.Ctx) error {
		return ps.HandleProjects(c)
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}
