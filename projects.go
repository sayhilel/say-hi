package projects

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
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
}

type Plist []Project

func InitProjects() Plist {
	dat, err := os.ReadFile("./public/data/projects.toml")
	check(err)
	pl := Plist{}
	p := Project{}
	_, err = toml.Decode(string(dat), &p)
	check(err)
	fmt.Println(p.Name)

	return pl
}
