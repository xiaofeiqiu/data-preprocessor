package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/xiaofeiqiu/mlstock/handlers"
	"github.com/xiaofeiqiu/mlstock/lib/log"
	"github.com/xiaofeiqiu/mlstock/lib/restutils"
	"github.com/xiaofeiqiu/mlstock/services/ioutils"
	"net/http"
	"time"
)

const (
	Timeout    = 60
	Throttle   = 10
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

	apiKey, err := ioutils.LoadApiKey(ApiKeyPath)
	if err != nil {
		log.Panic("Load api key", err.Error())
	}

	apiHandler := handlers.ApiHandler{
		Host:   "https://www.alphavantage.co",
		ApiKey: apiKey,
		HttpClient: &restutils.HttpClient{
			Client: &http.Client{},
		},
	}

	r.Get("/mlstock/health", restutils.Health)
	r.Get("/mlstock/dailyadjusted", handlers.ErrorHandler(apiHandler.GetDailyAdjusted))
	http.ListenAndServe(":8080", r)
}
