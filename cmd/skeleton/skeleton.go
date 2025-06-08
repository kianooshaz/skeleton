package main

import (
	"log"

	"github.com/kianooshaz/skeleton/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Printf("error at serving: %s", err.Error())

		return
	}
}
