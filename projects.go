package projects

import (
	"github.com/BurntSushi/toml"
	"github.com/gofiber/fiber/v2"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Project struct {
	Name        string
	Url         string
	Description string
	Image       string
}

type ListProjects struct {
	Pl []Project
}

func InitProjects() ListProjects {

	ps := ListProjects{}

	dat, err := os.ReadFile("./public/data/projects.toml")
	check(err)

	_, err = toml.Decode(string(dat), &ps)
	check(err)

	return ps
}

func (pl ListProjects) HandleProjects(c *fiber.Ctx) error {
	index, err := c.ParamsInt("index", 0)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid project index")
	}

	p := pl.Pl[index]

	return c.Render("layouts/project", fiber.Map{
		"Name":        p.Name,
		"Description": p.Description,
		"Url":         p.Url,
		"Image":       p.Image,
	})
}
