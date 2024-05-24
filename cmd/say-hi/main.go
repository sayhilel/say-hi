package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayhilel/say-hi/handlers"
)

func main() {

	app := fiber.New()

	app.Get("/", handlers.LandingHandler)
	app.Get("/secret", handlers.SecretHandler)

	app.Listen(":3000")

}
