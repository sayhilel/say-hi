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
	"strings"
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

		func(c *fiber.Ctx) error {
			userAgent := c.Get("User-Agent")
			if strings.Contains(strings.ToLower(userAgent), "mobile") {
				return c.Redirect("/mobile")
			}
			return c.Next()
		},
	)

	app.Get("/mobile", handlers.HandleInvalid)
	app.Get("/", handlers.LandingHandler)
	app.Get("/about-me", handlers.ViewAboutMe)
	app.Get("/landing", handlers.ViewLanding)
	app.Get("/contact-me", handlers.ViewContactMe)
	app.Post("/command", handlers.HandleCommands)

	ps := projects.InitProjects()
	app.Get("/projects/:index", func(c *fiber.Ctx) error {
		return ps.HandleProjects(c)
	})
	app.Get("/projects/:index/view", func(c *fiber.Ctx) error {
		return c.Redirect("https://gofiber.net", 301)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
