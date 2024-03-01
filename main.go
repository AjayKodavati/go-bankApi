package main

import (
	"log"

	"main.go/api"
)

func main() {

	store, err := api.NewPostgresStore()

	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":3000", store)
	server.Run()
}
