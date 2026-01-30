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

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		//ToDo log that config was unable to  be loaded
		log.Fatal(cfg)
	}
	server := &http.Server{
		Addr: cfg.ServerAddr,
	}

	idleConnsClosed := make(chan struct{})
	EnableGracefulShutdown(idleConnsClosed, server)

	//create global limiter & middleware
	rateLimitGlobally, globalRateLimiter, err := MakeGlobalRateLimitMiddleware(InMemory, cfg.GlobalLimiterCount, cfg.GlobalLimiterCap, cfg.GlobalLimiterRate)
	if err != nil {
		log.Fatal(err)
	}
	defer globalRateLimiter.Offline()

	//create per client limiter & middleware
	rateLimitPerClient, perClientRateLimiter, err := MakePerClientRateLimitMiddleware(InMemory, cfg.PerClientLimiterCap, cfg.PerClientLimiterLimit, cfg.PerClientLimiterWindow, cfg.PerClientLimiterClientTtl)

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
	shortener, err := NewUrlShortener(InMemory, cfg.ShortenerCap, cfg.ShortenerTTL, cfg.ShortCodeLength)
	if err != nil {
		log.Fatal(err)
	}
	defer shortener.Offline()

	//create app struct with methods for api handler logic
	app := &App{cfg, shortener, globalRateLimiter, perClientRateLimiter}

	//Route handlers
	mux := http.NewServeMux()
	mux.Handle("/", rateLimitGlobally(MakeIndexHandler()))
	mux.Handle("GET /{shortUrl}", rateLimitGlobally(http.HandlerFunc(app.RetrieveUrl)))
	mux.Handle("POST /api/shorten", withMiddlewares(http.HandlerFunc(app.ShortenUrl)))
	mux.Handle("GET /api/metrics/stream", rateLimitGlobally(http.HandlerFunc(app.StreamMetrics)))
	mux.Handle("GET /api/stress-test/stream", stressTestMiddlewares(http.HandlerFunc(app.StressTest)))
	server.Handler = SetupCors(mux, cfg)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

	//wait for graceful shutdown
	<-idleConnsClosed

}
