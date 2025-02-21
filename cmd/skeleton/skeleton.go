package main

import (
	"flag"
	"log"

	"github.com/kianooshaz/skeleton/internal/app"
)

func main() {
	configPath := flag.String("config", "config.yml", "yaml config path file")
	flag.Parse()

	if err := app.Run(*configPath); err != nil {
		log.Printf("error at serving: %s", err.Error())

		return
	}

}
