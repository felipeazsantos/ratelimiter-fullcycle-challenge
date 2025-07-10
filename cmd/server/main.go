package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/config/app"
	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	if err := app.StartDependencies(router); err != nil {
		log.Fatal("failed to start dependencies: ", err)
	}

	fmt.Printf("running the server on port: %s", "8080")

	log.Fatal(http.ListenAndServe(":8080", router))
}
