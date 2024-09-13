package rest

import (
	"context"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

type Config struct {
	Debug         bool          `yaml:"debug"`
	Address       string        `yaml:"address" validate:"required"`
	ReadTimeout   time.Duration `yaml:"read_timeout"`
	WriteTimeout  time.Duration `yaml:"write_timeout"`
	IdleTimeout   time.Duration `yaml:"idle_timeout"`
	BodyLimitSize string        `yaml:"body_limit_size"`
	CORS          struct {
		AllowedOrigins   []string `yaml:"allowed_origins"`
		AllowedHeaders   []string `yaml:"allowed_headers"`
		AllowedMethods   []string `yaml:"allowed_methods"`
		AllowCredentials bool     `yaml:"allow_credentials"`
		ExposedHeaders   []string `yaml:"exposed_headers"`
		MaxAge           int      `yaml:"max_age"`
	} `yaml:"cors"`
	RateLimit struct {
		Rate     float64       `yaml:"rate"`
		Burst    int           `yaml:"burst"`
		Duration time.Duration `yaml:"duration"`
	} `yaml:"rate_limit"`
}

type Server struct {
	core    *echo.Echo
	address string
}

func New(cfg *Config) *Server {
	e := echo.New()

	e.HideBanner = true
	e.IPExtractor = echo.ExtractIPFromXFFHeader()
	e.Server.ReadTimeout = cfg.ReadTimeout
	e.Server.WriteTimeout = cfg.WriteTimeout
	e.Server.IdleTimeout = cfg.IdleTimeout
	e.Server.ErrorLog = slog.NewLogLogger(slog.Default().Handler(), slog.LevelError)

	// Middlewares
	e.Use(echomw.Recover())
	e.Use(echomw.CORSWithConfig(echomw.CORSConfig{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		ExposeHeaders:    cfg.CORS.ExposedHeaders,
		MaxAge:           cfg.CORS.MaxAge,
	}))
	e.Use(echomw.Secure())
	e.Use(echomw.BodyLimit(cfg.BodyLimitSize))

	server := &Server{
		core:    e,
		address: cfg.Address,
	}

	server.registerRoutes()

	return server
}

func (s *Server) Start() error {
	return s.core.Start(s.address)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.core.Shutdown(ctx)
}

func (s *Server) Close() error {
	return s.core.Close()
}
