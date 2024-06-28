package projects

import (
	"github.com/gofiber/fiber/v2"

	"github.com/BurntSushi/toml"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Project struct {
	Name        string
	Description string
	Url         string
	Image       string
}

type Projects struct {
	PL []Project `toml:"project"`
}

func InitProjects() Projects {

	ps := Projects{}

	dat, err := os.ReadFile("./public/data/projects.toml")
	check(err)

	_, err = toml.Decode(string(dat), &ps)
	check(err)

	return ps
}

func (P Projects) ViewProjects(c *fiber.Ctx) error {
	return c.Render("layouts/projects", fiber.Map{
		"numProjects": len(P.PL),
	})
}

func (P Projects) HandleProjects(c *fiber.Ctx) error {
	index, err := c.ParamsInt("index", 0)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid project index")
	}

	p := P.PL[index]

	return c.Render("layouts/project", fiber.Map{
		"Name":        p.Name,
		"Description": p.Description,
		"Url":         p.Url,
		"Image":       p.Image,
	})
}
