package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
	"github.com/sayhilel/say-hi/internal/quotes"
	"github.com/sayhilel/say-hi/internal/runtimeinfo"
	"github.com/sayhilel/say-hi/internal/store"
)

const visitedCookie = "sayhi_visited"

type Handler struct {
	Store store.Store
}

func New(s store.Store) *Handler {
	return &Handler{Store: s}
}

func (h *Handler) HandleCommands(c *fiber.Ctx) error {
	switch strings.TrimSpace(c.FormValue("command")) {

	default:
		return HandleInvalid(c)
	case "sudo":
		return secretUser(c)
	case "whoami":
		return ViewAboutMe(c)
	case "exp":
		return ViewExperience(c)
	case "showcase":
		return ViewProjects(c)
	case "ping":
		return ViewContactMe(c)
	case "cloud":
		return h.ViewCloud(c)
	case "more":
		return OpenConfirmation(c)
	case "clear":
		return ClearField(c)
	}
}

func secretUser(c *fiber.Ctx) error {
	return c.Render("layouts/secret", fiber.Map{})
}

func ClearField(c *fiber.Ctx) error {
	return c.Render("layouts/prompt", fiber.Map{})
}

func HandleInvalid(c *fiber.Ctx) error {
	return c.Render("layouts/invalid", fiber.Map{})
}

func ViewLanding(c *fiber.Ctx) error {
	return c.Render("layouts/landing", quote.GetQuote())
}

func (h *Handler) LandingHandler(c *fiber.Ctx) error {
	userAgent := c.Get("User-Agent")
	if strings.Contains(strings.ToLower(userAgent), "mobile") {
		return c.Render("layouts/mobile", fiber.Map{})
	}

	return c.Render("index", fiber.Map{
		"q":       "Use (C-c) to refresh this page for a random quote.",
		"a":       "Sahil Sinha",
		"visitor": h.countVisitor(c),
	})
}

// countVisitor increments the visitor counter once per browser per day. A
// storage failure never blocks the page; the counter is simply hidden.
func (h *Handler) countVisitor(c *fiber.Ctx) int64 {
	if c.Cookies(visitedCookie) != "" {
		return 0
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
	defer cancel()

	n, err := h.Store.IncrementVisits(ctx)
	if err != nil {
		slog.Error("increment visits", "err", err)
		return 0
	}

	c.Cookie(&fiber.Cookie{
		Name:     visitedCookie,
		Value:    "1",
		MaxAge:   int((24 * time.Hour).Seconds()),
		HTTPOnly: true,
		Secure:   c.Protocol() == "https",
		SameSite: fiber.CookieSameSiteLaxMode,
	})
	return n
}

func OpenConfirmation(c *fiber.Ctx) error {
	return c.Render("layouts/dialog-box", fiber.Map{
		"appl": "Open resume?",
	})
}

func ViewAboutMe(c *fiber.Ctx) error {
	return c.Render("layouts/about-me", fiber.Map{})
}

func ViewExperience(c *fiber.Ctx) error {
	return c.Render("layouts/exp", fiber.Map{})
}

func ViewContactMe(c *fiber.Ctx) error {
	return c.Render("layouts/contact-me", fiber.Map{})
}

func ViewProjects(c *fiber.Ctx) error {
	return c.Render("layouts/projects", fiber.Map{
		"shown": "hidden",
	})
}

func (h *Handler) ViewCloud(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
	defer cancel()

	storeStatus := "connected"
	if err := h.Store.Ping(ctx); err != nil {
		storeStatus = "unreachable"
	}

	return c.Render("layouts/cloud", fiber.Map{
		"Info":        runtimeinfo.Current(),
		"Backend":     h.Store.Backend(),
		"StoreStatus": storeStatus,
	})
}

// ContactForm is the payload of the `ping` form. Website is a honeypot field
// that is hidden from people; bots that fill it in are silently dropped.
type ContactForm struct {
	Name    string `form:"name"`
	Email   string `form:"email"`
	Message string `form:"message"`
	Website string `form:"website"`
}

func (f ContactForm) Validate() string {
	switch {
	case f.Name == "" || utf8.RuneCountInString(f.Name) > 100:
		return "name must be between 1 and 100 characters"
	case len(f.Email) > 254:
		return "email is too long"
	case f.Message == "" || utf8.RuneCountInString(f.Message) > 2000:
		return "message must be between 1 and 2000 characters"
	}
	if addr, err := mail.ParseAddress(f.Email); err != nil || addr.Address != f.Email {
		return "email address is not valid"
	}
	return ""
}

func (h *Handler) SubmitContact(c *fiber.Ctx) error {
	var f ContactForm
	if err := c.BodyParser(&f); err != nil {
		return c.Status(fiber.StatusBadRequest).Render("layouts/contact-result",
			fiber.Map{"Err": "could not read form"})
	}
	f.Name = strings.TrimSpace(f.Name)
	f.Email = strings.TrimSpace(f.Email)
	f.Message = strings.TrimSpace(f.Message)

	if f.Website != "" {
		slog.Info("contact honeypot triggered", "ip", c.IP())
		return c.Render("layouts/contact-result", fiber.Map{})
	}

	if msg := f.Validate(); msg != "" {
		return c.Status(fiber.StatusUnprocessableEntity).Render("layouts/contact-result",
			fiber.Map{"Err": msg})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	err := h.Store.SaveMessage(ctx, store.Message{
		ID:        newID(),
		Kind:      "contact",
		Name:      f.Name,
		Email:     f.Email,
		Body:      f.Message,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		slog.Error("save contact message", "err", err)
		return c.Status(fiber.StatusServiceUnavailable).Render("layouts/contact-result",
			fiber.Map{"Err": "message could not be delivered, please use email instead"})
	}

	slog.Info("contact message saved")
	return c.Render("layouts/contact-result", fiber.Map{})
}

func newID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
