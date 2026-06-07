package main

import (
	"log"
	"os"
)

func main() {
	cfg := config{
		addr: ":8080",
		db:   dbConfig{},
	}

	api := application{
		config: cfg,
	}

	handler, err := api.mount()
	if handler != nil {
		err = api.run(handler)
	}

	if err != nil {
		log.Printf("Application closed with error: %s", err.Error())
		os.Exit(1)
	}
}
