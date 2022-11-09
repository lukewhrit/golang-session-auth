package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis/v8"
)

const (
	databaseURI = "redis://localhost:6379/0"
)

var rdb *redis.Client
var ctx = context.Background()

func init() {
	opt, err := redis.ParseURL(databaseURI)
	if err != nil {
		panic(err)
	}

	rdb = redis.NewClient(opt)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/signin", signin)
	r.Get("/welcome", welcome)

	http.ListenAndServe(":3000", r)
}
