package projects

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func TestHandleProjectsIndexBounds(t *testing.T) {
	ps := Projects{PL: []Project{{Name: "only"}}}
	app := fiber.New(fiber.Config{Views: html.New("../../views", ".html")})
	app.Get("/projects/:index", ps.HandleProjects)

	cases := map[string]int{
		"/projects/0":   fiber.StatusOK,
		"/projects/1":   fiber.StatusBadRequest,
		"/projects/-1":  fiber.StatusBadRequest,
		"/projects/abc": fiber.StatusBadRequest,
	}
	for path, want := range cases {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, path, nil))
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != want {
			t.Errorf("%s: got %d, want %d", path, resp.StatusCode, want)
		}
	}
}
