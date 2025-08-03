package main

import (
	"log/slog"

	"github.com/kianooshaz/skeleton/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		slog.Error("Error encountered while running app", "error", err)

		return
	}
}
