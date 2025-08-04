package rest

import (
	"context"
	"log/slog"
	"time"

	"github.com/kianooshaz/skeleton/foundation/session"
	"github.com/kianooshaz/skeleton/internal/app/web/protocol"
	passwordproto "github.com/kianooshaz/skeleton/services/authentication/password/proto"
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
		Enable           bool     `yaml:"enable"`
		AllowedOrigins   []string `yaml:"allowed_origins"`
		AllowedHeaders   []string `yaml:"allowed_headers"`
		AllowedMethods   []string `yaml:"allowed_methods"`
		AllowCredentials bool     `yaml:"allow_credentials"`
		ExposedHeaders   []string `yaml:"exposed_headers"`
		MaxAge           int      `yaml:"max_age"`
	}
	RateLimit struct {
		Enable   bool          `yaml:"enable"`
		Rate     float64       `yaml:"rate"`
		Burst    int           `yaml:"burst"`
		Duration time.Duration `yaml:"duration"`
	}
}

type Services struct {
	userService passwordproto.PasswordService
	// Add other services as needed
}

type server struct {
	core     *echo.Echo
	address  string
	logger   *slog.Logger
	services Services
}

func New(cfg Config, logger *slog.Logger, services Services) (protocol.WebService, error) {
	e := echo.New()

	e.Debug = cfg.Debug
	e.HideBanner = true
	e.IPExtractor = echo.ExtractIPFromXFFHeader()
	e.Server.ReadTimeout = cfg.ReadTimeout
	e.Server.WriteTimeout = cfg.WriteTimeout
	e.Server.IdleTimeout = cfg.IdleTimeout
	e.Server.ErrorLog = slog.NewLogLogger(logger.Handler(), slog.LevelError)
	e.HTTPErrorHandler = ErrorResponse

	// Middlewares
	e.Use(echomw.Recover())
	e.Use(echomw.RequestIDWithConfig(echomw.RequestIDConfig{
		RequestIDHandler: session.SetRequestIDEcho(),
	}))
	e.Use(echomw.Secure())

	if cfg.CORS.Enable {
		e.Use(echomw.CORSWithConfig(echomw.CORSConfig{
			AllowOrigins:     cfg.CORS.AllowedOrigins,
			AllowMethods:     cfg.CORS.AllowedMethods,
			AllowHeaders:     cfg.CORS.AllowedHeaders,
			AllowCredentials: cfg.CORS.AllowCredentials,
			ExposeHeaders:    cfg.CORS.ExposedHeaders,
			MaxAge:           cfg.CORS.MaxAge,
		}))
	}

	if cfg.RateLimit.Enable {
		// TODO implement rate limiter with redis and echo
	}

	if cfg.BodyLimitSize != "" {
		e.Use(echomw.BodyLimit(cfg.BodyLimitSize))
	}

	server := &server{
		core:     e,
		address:  cfg.Address,
		logger:   logger,
		services: services,
	}

	server.registerRoutes()

	return server, nil
}

func (s *server) Start() error {
	return s.core.Start(s.address)
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.core.Shutdown(ctx)
}

func (s *server) Close() error {
	return s.core.Close()
}
