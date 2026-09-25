package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/sayhilel/say-hi/internal/store"
)

func newTestApp(t *testing.T) (*fiber.App, *store.Memory) {
	t.Helper()
	mem := store.NewMemory()
	h := New(mem)
	app := fiber.New(fiber.Config{Views: html.New("../../views", ".html")})
	app.Get("/", h.LandingHandler)
	app.Post("/command", h.HandleCommands)
	app.Post("/contact", h.SubmitContact)
	return app, mem
}

func post(t *testing.T, app *fiber.App, path string, form url.Values) (int, string) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(body)
}

func TestCommandRouting(t *testing.T) {
	app, _ := newTestApp(t)

	cases := map[string]string{
		"whoami":   "Whoami: About-Me",
		"exp":      "Professional Experience",
		"showcase": "Showcase: Projects",
		"ping":     "contact-form",
		"cloud":    "Live runtime",
		" cloud ":  "Live runtime",
		"rm -rf /": "Invalid Command",
		"":         "Invalid Command",
	}
	for cmd, want := range cases {
		status, body := post(t, app, "/command", url.Values{"command": {cmd}})
		if status != fiber.StatusOK {
			t.Errorf("%q: status %d", cmd, status)
		}
		if !strings.Contains(body, want) {
			t.Errorf("%q: body missing %q", cmd, want)
		}
	}
}

func TestLandingCountsVisitorOnce(t *testing.T) {
	app, _ := newTestApp(t)

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/", nil))
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "viewer #1") {
		t.Fatalf("first visit not counted: %s", body)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	for _, c := range resp.Cookies() {
		req.AddCookie(c)
	}
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ = io.ReadAll(resp.Body)
	if strings.Contains(string(body), "viewer #") {
		t.Fatalf("returning visitor counted again")
	}
}

func TestSubmitContact(t *testing.T) {
	valid := url.Values{
		"name":    {"Recruiter"},
		"email":   {"recruiter@example.com"},
		"message": {"Let's talk about the cloud role."},
	}

	t.Run("valid message is stored", func(t *testing.T) {
		app, mem := newTestApp(t)
		status, body := post(t, app, "/contact", valid)
		if status != fiber.StatusOK || !strings.Contains(body, "delivered") {
			t.Fatalf("status %d body %s", status, body)
		}
		if len(mem.Messages) != 1 || mem.Messages[0].Email != "recruiter@example.com" {
			t.Fatalf("message not stored: %+v", mem.Messages)
		}
	})

	t.Run("honeypot is dropped silently", func(t *testing.T) {
		app, mem := newTestApp(t)
		form := url.Values{}
		for k, v := range valid {
			form[k] = v
		}
		form.Set("website", "http://spam.example")
		status, _ := post(t, app, "/contact", form)
		if status != fiber.StatusOK || len(mem.Messages) != 0 {
			t.Fatalf("status %d, stored %d", status, len(mem.Messages))
		}
	})

	invalid := map[string]url.Values{
		"missing name":  {"name": {""}, "email": valid["email"], "message": valid["message"]},
		"bad email":     {"name": valid["name"], "email": {"not-an-email"}, "message": valid["message"]},
		"display email": {"name": valid["name"], "email": {"Bob <bob@example.com>"}, "message": valid["message"]},
		"long message":  {"name": valid["name"], "email": valid["email"], "message": {strings.Repeat("a", 2001)}},
	}
	for name, form := range invalid {
		t.Run(name, func(t *testing.T) {
			app, mem := newTestApp(t)
			status, _ := post(t, app, "/contact", form)
			if status != fiber.StatusUnprocessableEntity || len(mem.Messages) != 0 {
				t.Fatalf("status %d, stored %d", status, len(mem.Messages))
			}
		})
	}
}
