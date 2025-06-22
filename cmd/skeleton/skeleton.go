package main

import (
	"log"

	"github.com/kianooshaz/skeleton/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Printf("Error encountered while running app: %s", err.Error())

		return
	}
}
