package web

import (
	"fmt"

	"github.com/kianooshaz/skeleton/internal/app/web/rest"
)

func Serve(configPath string) error {
	cfg, err := newConfig(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	e := rest.New(cfg.Rest)

	return e.Start()
}
