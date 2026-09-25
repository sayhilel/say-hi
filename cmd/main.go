package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/sayhilel/say-hi/internal/handlers"
	"github.com/sayhilel/say-hi/internal/projects"
	"github.com/sayhilel/say-hi/internal/runtimeinfo"
	"github.com/sayhilel/say-hi/internal/store"
)

func main() {
	// JSON logs on stdout are collected by Container Apps into Log Analytics,
	// where every field becomes queryable with KQL.
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	st, err := store.FromEnv()
	if err != nil {
		slog.Error("init store", "err", err)
		os.Exit(1)
	}
	slog.Info("store ready", "backend", st.Backend())

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views:                 engine,
		DisableStartupMessage: true,
		// Container Apps ingress (Envoy) terminates TLS and forwards the client IP.
		ProxyHeader:  fiber.HeaderXForwardedFor,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	app.Use(
		recover.New(),
		requestLogger,
		helmet.New(helmet.Config{
			// The `sudo` easter egg embeds YouTube and HTMX loads from unpkg,
			// so cross-origin embedding stays allowed.
			CrossOriginEmbedderPolicy: "unsafe-none",
			CrossOriginResourcePolicy: "cross-origin",
			ReferrerPolicy:            "strict-origin-when-cross-origin",
			HSTSMaxAge:                31536000,
		}),
		favicon.New(favicon.Config{
			File: "public/css/images/favicon.ico",
			URL:  "/favicon.ico",
		}),
	)

	h := handlers.New(st)

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})
	app.Get("/readyz", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
		defer cancel()
		if err := st.Ping(ctx); err != nil {
			slog.Error("readiness check failed", "err", err)
			return c.Status(fiber.StatusServiceUnavailable).SendString("store unreachable")
		}
		return c.SendString("ok")
	})

	app.Get("/", h.LandingHandler)
	app.Get("/landing", handlers.ViewLanding)
	app.Post("/command", h.HandleCommands)
	app.Post("/contact", limiter.New(limiter.Config{
		Max:          5,
		Expiration:   10 * time.Minute,
		KeyGenerator: clientIP,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).Render("layouts/contact-result",
				fiber.Map{"Err": "too many messages, try again in a few minutes"})
		},
	}), h.SubmitContact)

	ps := projects.InitProjects()
	app.Get("/projects/:index", ps.HandleProjects)

	// Registered after the routes: a static miss resets the response, which
	// would drop the security headers set above.
	app.Static("/", "./public")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Container Apps sends SIGTERM when scaling in or rolling a revision;
	// finish in-flight requests before exiting.
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		slog.Info("shutting down")
		_ = app.ShutdownWithTimeout(10 * time.Second)
	}()

	slog.Info("listening", "port", port, "version", runtimeinfo.Current().Version)
	if err := app.Listen(":" + port); err != nil {
		slog.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

func requestLogger(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()

	// Health probes fire every few seconds; logging them only burns the
	// Log Analytics free ingestion allowance.
	if c.Path() == "/healthz" {
		return err
	}
	runtimeinfo.CountRequest()

	status := c.Response().StatusCode()
	if err != nil {
		if e, ok := err.(*fiber.Error); ok {
			status = e.Code
		} else {
			status = fiber.StatusInternalServerError
		}
	}

	slog.Info("request",
		"method", c.Method(),
		"path", c.Path(),
		"status", status,
		"durationMs", time.Since(start).Milliseconds(),
		"ip", clientIP(c),
		"userAgent", c.Get(fiber.HeaderUserAgent),
	)
	return err
}

// clientIP prefers X-Forwarded-For (set by Container Apps ingress) and falls
// back to the socket address when running without a proxy.
func clientIP(c *fiber.Ctx) string {
	if ip := c.IP(); ip != "" {
		return ip
	}
	return c.Context().RemoteIP().String()
}
