package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sayhilel/say-hi/content"
)

func LandingHandler(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

func ProjectsHandler(c *fiber.Ctx, pl []content.Project) error {
	pn, err := strconv.Atoi(c.Params("num")) 
	if err != nil{
		fmt.Print("Error")
	}

	p := pl[pn]


	return c.Render("project", fiber.Map{
		"Name":        p.Name,
		"Description": p.Description,
		"Url":         p.Url,
	})
}
