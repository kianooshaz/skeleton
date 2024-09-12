package web

import "fmt"

func Serve(configPath string) error {
	cfg, err := newConfig(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	fmt.Println("version:", cfg.Version)

	return nil
}
