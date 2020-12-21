package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/xiaofeiqiu/data-preprocessor/handlers"
	"github.com/xiaofeiqiu/data-preprocessor/lib/log"
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"github.com/xiaofeiqiu/data-preprocessor/services/alphavantage"
	"net/http"
	"os"
	"time"
)

const (
	Timeout    = 60
	Throttle   = 10
	AlphaVantageKey = "ALPHA_VANTAGE"
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

	apiKey := os.Getenv(AlphaVantageKey)
	if apiKey == "" {
		log.Fatal("GetApiKey","Api key not found")
	}

	alphaVantageApi := &alphavantage.AlphaVantageApi{
		Host:   "https://www.alphavantage.co",
		ApiKey: apiKey,
		HttpClient: &restutils.HttpClient{
			Client: &http.Client{},
		},
	}

	apiHandler := handlers.ApiHandler{
		AlphaVantageApi: alphaVantageApi,
	}

	r.Get("/preprocessor/health", restutils.Health)
	r.Post("/preprocessor/candle/dailyadjusted", handlers.ErrorHandler(apiHandler.InsertDailyCandle))
	r.Put("/preprocessor/ema8/dailyadjusted", handlers.ErrorHandler(apiHandler.FillDailyEMA))
	http.ListenAndServe(":8080", r)
}
