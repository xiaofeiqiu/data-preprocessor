package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/xiaofeiqiu/mlstock/handlers"
	"github.com/xiaofeiqiu/mlstock/lib/logger"
	"net/http"
	"time"
)

const (
	Timeout  = 60
	Throttle = 10
)

func main() {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(logger.NewMiddlewareLogger())
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Timeout(Timeout * time.Second))
	r.Use(middleware.Throttle(Throttle))

	r.Get("/mlstock/health", handlers.Health)
	http.ListenAndServe(":8080", r)

}