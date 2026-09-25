package projects

import (
	"github.com/BurntSushi/toml"
	"github.com/gofiber/fiber/v2"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Project struct {
	Name         string
	Technologies string
	Description  string
	Url          string
	Ascii        string
	Accent       string
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
	index, err := strconv.Atoi(c.Params("index"))
	if err != nil || index < 0 || index >= len(P.PL) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid project index")
	}

	p := P.PL[index]

	return c.Render("layouts/project", fiber.Map{
		"Name":         p.Name,
		"Technologies": p.Technologies,
		"Description":  p.Description,
		"Url":          p.Url,
		"Ascii":        p.Ascii,
		"Accent":       p.Accent,
	})
}
