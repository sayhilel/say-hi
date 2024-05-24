package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func LandingHandler(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func SecretHandler(c *fiber.Ctx) error {
	fmt.Println("Secret")
	return c.SendString("Nice")

}
