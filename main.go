package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HTTPContext struct {
	w   http.ResponseWriter
	req *http.Request
}

type RouteHandler func(ctx HTTPContext)
type LimiterMiddleware func(http.Handler) http.Handler

type UrlShortenerPayload struct {
	Original string `json:"original"`
}

type StorageType int

const (
	InMemory StorageType = iota
	Redis
)

const maxUrlLength = 2048

func main() {
	server := &http.Server{
		Addr: ":8090",
	}

	idleConnsClosed := make(chan struct{})
	limiterMiddleware, shorten, retrieve := initialize(idleConnsClosed, server)

	//create server
	mux := http.NewServeMux()
	mux.Handle("/", limiterMiddleware(AsHandler(Index)))
	mux.Handle("GET /{shortUrl}", limiterMiddleware(AsHandler(retrieve)))
	mux.Handle("POST /shorten", limiterMiddleware(AsHandler(shorten)))
	server.Handler = mux

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

	//wait for graceful shutdown
	<-idleConnsClosed
}

func createUrlShortener() (UrlShortener, RouteHandler, RouteHandler) {
	s, err := NewUrlShortener(InMemory, 100000, time.Hour)
	if err != nil {
		log.Fatal(err)
	}
	return s,
		func(ctx HTTPContext) {
			ShortenUrl(s, ctx)
		},
		func(ctx HTTPContext) {
			RetrieveUrl(s, ctx)
		}
}

func createLimiterMiddleware() (*Limiter, LimiterMiddleware) {
	var limiter *Limiter
	if l, err := NewLimiter(InMemory, 50000, 50000, 10000); err != nil {
		log.Fatal(err)
	} else {
		limiter = l
	}

	return limiter, func(next http.Handler) http.Handler {
		return AsHandler(func(ctx HTTPContext) { LimiteRate(limiter, ctx, next) })
	}
}

func initialize(done chan struct{}, server *http.Server) (LimiterMiddleware, RouteHandler, RouteHandler) {
	//create limiter & middleware
	limiter, limiterMiddleware := createLimiterMiddleware()
	//create url shortener
	shortener, shorten, retrieve := createUrlShortener()

	// enable Graceful Exit
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		//when interrupt received
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		limiter.Stop()
		shortener.Offline()
		close(done)
	}()

	return limiterMiddleware, shorten, retrieve
}
