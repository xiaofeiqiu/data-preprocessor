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
	"github.com/xiaofeiqiu/data-preprocessor/services/dbservice"
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

	log.Info("Init", "Get app config successful")

	defaultClient := &restutils.HttpClient{
		Client: &http.Client{},
	}
	alphaVantageApi := &alphavantage.AlphaVantageClient{
		Host:       AlphavantageHost,
		ApiKey:     config.AlphaVantageKey,
		HttpClient: defaultClient,
	}

	dbClient, err := db.NewPostgresDBClient(DBHost, DBName, config.DBUsername, config.DBPassword)
	if err != nil {
		log.Fatal("NewPostgresDBClient", "error connecting to db, "+err.Error())
	}
	log.Info("Init", "Create db client successful")

	dbService := dbservice.NewDBService(dbClient)

	apiHandler := handlers.ApiHandler{
		AlphaVantageClient: alphaVantageApi,
		DBService:          dbService,
		DefaultClient:      defaultClient,
	}

	err = apiHandler.DBService.InitDBTableMapping()
	if err != nil {
		log.Fatal("Init", "error init db table mapping")
	}

	r.Get("/preprocessor/health", restutils.Health)

	r.Post("/preprocessor/candle/dailyadjusted", handlers.ErrorHandler(apiHandler.InsertDailyCandle))
	r.Post("/preprocessor/candle/missingdailyadjusted", handlers.ErrorHandler(apiHandler.InsertMissingDailyCandle))
	r.Post("/preprocessor/processpr/doall", handlers.ErrorHandler(apiHandler.Doall))

	r.Put("/preprocessor/ema/dailyadjusted", handlers.ErrorHandler(apiHandler.FillDailyEMA))
	r.Put("/preprocessor/cci/dailyadjusted", handlers.ErrorHandler(apiHandler.FillDailyCCI))
	r.Put("/preprocessor/aroon/dailyadjusted", handlers.ErrorHandler(apiHandler.FillDailyAroon))
	r.Put("/preprocessor/macd/dailyadjusted", handlers.ErrorHandler(apiHandler.FillDailyMacd))

	r.Delete("/preprocessor/processpr/dailyrawdata", handlers.ErrorHandler(apiHandler.ClearRawData))

	http.ListenAndServe(":8080", r)
	log.Info("Init", "Server started")
}
