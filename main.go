package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/xiaofeiqiu/mlstock/handlers"
	"github.com/xiaofeiqiu/mlstock/lib/log"
	"net/http"
	"time"
)

const (
	Timeout  = 60
	Throttle = 10
	ApiKeyPath = ".apiKey"
)

func main() {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(log.NewMiddlewareLogger())
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Timeout(Timeout * time.Second))
	r.Use(middleware.Throttle(Throttle))

	r.Get("/mlstock/health", handlers.Health)
	http.ListenAndServe(":8080", r)

}