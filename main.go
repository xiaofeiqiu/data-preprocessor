package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/xiaofeiqiu/data-preprocessor/handlers"
	"github.com/xiaofeiqiu/data-preprocessor/lib/db"
	"github.com/xiaofeiqiu/data-preprocessor/lib/log"
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"github.com/xiaofeiqiu/data-preprocessor/services/alphavantage"
	"github.com/xiaofeiqiu/data-preprocessor/internal"
	"net/http"
	"time"
)

const (
	Timeout          = 60
	Throttle         = 10
	AlphavantageHost = "https://www.alphavantage.co"
	DBName           = "golddigger"
	DBHost           = "localhost"
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

	config, err := internal.GetAppConfig()
	if err != nil {
		log.Fatal("GetAppConfig", "error getting app config")
	}

	alphaVantageApi := &alphavantage.AlphaVantageClient{
		Host:   AlphavantageHost,
		ApiKey: config.AlphaVantageKey,
		HttpClient: &restutils.HttpClient{
			Client: &http.Client{},
		},
	}

	dbClient, err := db.NewPostgresDBClient(DBHost, DBName, config.DBUsername, config.DBPassword)
	if err != nil {
		log.Fatal("NewPostgresDBClient", "error connecting to db, "+err.Error())
	}

	apiHandler := handlers.ApiHandler{
		AlphaVantageClient: alphaVantageApi,
		DBClient:           dbClient,
	}

	apiHandler.InitDBTableMapping()

	r.Get("/preprocessor/health", restutils.Health)
	r.Post("/preprocessor/candle/dailyadjusted", handlers.ErrorHandler(apiHandler.InsertDailyCandle))
	r.Post("/preprocessor/candle/missingdailyadjusted", handlers.ErrorHandler(apiHandler.InsertMissingDailyCandle))
	r.Put("/preprocessor/ema8/dailyadjusted", handlers.ErrorHandler(apiHandler.FillDailyEMA))
	http.ListenAndServe(":8080", r)
}
