package web

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Serve(configPath string) error {
	cfg, err := newConfig(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	e := echo.New()
	e.GET("/version", func(c echo.Context) error {
		return c.String(http.StatusOK, cfg.Version)
	})

	return e.Start(":1323")
}
