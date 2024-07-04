package quote

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

func GetQuote() map[string]string {
	agent := fiber.Get("https://api.quotable.io/random")
	_, body, errs := agent.Bytes()

	result := map[string]string{
		"author":  "Sahil Sinha",
		"content": "Woopsie",
	}

	if len(errs) > 0 {
		return result
	}

	json.Unmarshal(body, &result)

	return result
}
