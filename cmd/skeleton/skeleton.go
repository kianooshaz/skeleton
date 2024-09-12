package main

import (
	"flag"
	"log"

	"github.com/kianooshaz/skeleton/internal/app/web"
)

func main() {
	configPath := flag.String("config", "config.yaml", "yaml config path file")
	flag.Parse()

	if err := web.Serve(*configPath); err != nil {
		log.Printf("error at serving: %s", err.Error())
		return
	}
}
