package quote

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func GetQuote() map[string]string {
	agent := fiber.Get("https://zenquotes.io/api/random")
	_, body, errs := agent.Bytes()

	if len(errs) > 0 {
		return map[string]string{
			"Err": "Unable to Fetch, API might be down?",
		}
	}

	var quotes []map[string]string
	json.Unmarshal(body, &quotes)

	if len(quotes) == 0 {
		fmt.Println("QUOTES ARE BROKEN")
		return map[string]string{
			"Err": "No quotes found",
		}
	}

	quote := quotes[0]
	result := map[string]string{
		"q": quote["q"],
		"a": quote["a"],
	}

	return result
}
