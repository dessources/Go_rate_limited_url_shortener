package main

import (
	"log"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type StorageType int

const (
	InMemory StorageType = iota
	Redis
)

const maxUrlLength = 4096

var Cfg, cfgErr = LoadConfig()

func main() {

	if cfgErr != nil {
		//ToDo log that config was unable to  be loaded
		log.Fatal(cfgErr)
	}
	server := &http.Server{
		Addr: Cfg.ServerAddr,
	}

	idleConnsClosed := make(chan struct{})
	EnableGracefulShutdown(idleConnsClosed, server)

	//create global limiter & middleware
	rateLimitGlobally, globalRateLimiter, err := MakeGlobalRateLimitMiddleware(InMemory, Cfg.GlobalLimiterCount, Cfg.GlobalLimiterCap, Cfg.GlobalLimiterRate)
	if err != nil {
		log.Fatal(err)
	}
	defer globalRateLimiter.Offline()

	//create per client limiter & middleware
	rateLimitPerClient, perClientRateLimiter, err := MakePerClientRateLimitMiddleware(InMemory, Cfg.PerClientLimiterCap, Cfg.PerClientLimiterLimit, Cfg.PerClientLimiterWindow)

	if err != nil {
		log.Fatal(err)
	}
	defer perClientRateLimiter.Offline()

	//middleware composers
	withMiddlewares := ComposeMiddlewares(rateLimitGlobally, rateLimitPerClient)
	//composed middleware for stress test route
	stressTestMiddlewares, cleanup, err := MakeStressTestRouteMiddlewares()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	//url shortener struct
	shortener, err := NewUrlShortener(InMemory, Cfg.ShortenerCap, Cfg.ShortenerTTL)
	if err != nil {
		log.Fatal(err)
	}
	defer shortener.Offline()

	//create app struct with methods for api handler logic
	app := &App{shortener, globalRateLimiter, perClientRateLimiter}

	//Route handlers
	mux := http.NewServeMux()
	mux.Handle("/", rateLimitGlobally(MakeIndexHandler()))
	mux.Handle("GET /{shortUrl}", rateLimitGlobally(http.HandlerFunc(app.RetrieveUrl)))
	mux.Handle("POST /api/shorten", withMiddlewares(http.HandlerFunc(app.ShortenUrl)))
	mux.Handle("GET /api/metrics/stream", rateLimitGlobally(http.HandlerFunc(app.StreamMetrics)))
	mux.Handle("GET /api/stress-test/stream", stressTestMiddlewares(http.HandlerFunc(StressTest)))
	server.Handler = SetupCors(mux)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

	//wait for graceful shutdown
	<-idleConnsClosed

}
