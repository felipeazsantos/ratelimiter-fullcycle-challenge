package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/config/app"
	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/config/getenv"
	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/config/redisdb"
	"github.com/go-chi/chi/v5"
)

func main() {
	if err := getenv.InitConfig("../../.env", ".env"); err != nil {
		log.Fatal("unable to load app config: ", err)
	}

	router := chi.NewRouter()

	rdb := redisdb.NewRedisClient(getenv.AppConfig.RedisAddr)
	defer rdb.Close()

	if err := app.StartDependencies(router, rdb); err != nil {
		log.Fatal("failed to start dependencies: ", err)
	}

	fmt.Printf("running the server on port: %s", "8080")

	log.Fatal(http.ListenAndServe(":8080", router))
}
