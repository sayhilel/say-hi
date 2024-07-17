package quote

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

func GetQuote() map[string]string {
	agent := fiber.Get("https://api.quotable.io/random?tags=technology|education|Inspirational|Science")
	_, body, errs := agent.Bytes()

	result := map[string]string{}

	if len(errs) > 0 {
		return map[string]string{
			"Err": "Unable to Fetch, API might be down?",
		}
	}

	json.Unmarshal(body, &result)

	return result
}
