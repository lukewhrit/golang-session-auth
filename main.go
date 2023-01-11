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
	s := NewServer()
	s.MountHandlers()
	http.ListenAndServe(":3000", s.Router)
}

func NewServer() *Server {
	s := &Server{}
	s.Router = chi.NewRouter()
	return s
}

func (s *Server) MountHandlers() {
	s.Router.Use(middleware.Logger)

	s.Router.Post("/signin", signin)
	s.Router.Get("/welcome", welcome)
}
