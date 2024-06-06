package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type Project struct {
	Name string
	Url string
	Description string
}

// Test Vector
type Plist []Project


func (pl Plist) SwitchProject(c *fiber.Ctx) error{
	index, err := c.ParamsInt("index", 0)
	if err != nil{
		return c.Status(fiber.StatusBadRequest).SendString("Invalid project index")
	}
	
	p := pl[index]

	return c.Render("layouts/project", fiber.Map{
		"Name":        p.Name,
		"Description": p.Description,
		"Url":         p.Url,
	})
}
