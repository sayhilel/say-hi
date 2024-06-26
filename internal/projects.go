package projects

import (
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
