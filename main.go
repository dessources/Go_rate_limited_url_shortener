package main

import (
	"context"
	"fmt"
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

func main() {

	//create middleware
	limiter, limiterMiddleware := createLimiterMiddleware()

	//create server
	mux := http.NewServeMux()
	mux.Handle("/", limiterMiddleware(asHandler(index)))
	mux.Handle("POST /shorten", limiterMiddleware(asHandler(shorten)))

	server := &http.Server{
		Addr:    ":8090",
		Handler: mux,
	}

	idleConnsClosed := make(chan struct{})
	go enableGracefulExit(server, idleConnsClosed, limiter)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-idleConnsClosed
}

func index(ctx HTTPContext) {
	ctx.w.WriteHeader(http.StatusOK)
	fmt.Fprintf(ctx.w, "<h1>URL shortener is running!</h1>\n")
}

func shorten(request HTTPContext) {
	shortUrl := Shorten()
	request.w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(request.w, "Short URL: %s", shortUrl)
}

func limiteRate(limiter *Limiter, ctx HTTPContext, next http.Handler) {
	if limiter.Allow(1) {
		next.ServeHTTP(ctx.w, ctx.req)
	} else {
		ctx.w.WriteHeader(http.StatusTooManyRequests)
		_, _ = ctx.w.Write([]byte("Request rejected by rate limiter"))
		return
	}
}

func createLimiterMiddleware() (*Limiter, func(http.Handler) http.Handler) {

	limiter := NewLimiter(InMemory, 50000, 50000, 10000)
	return limiter, func(next http.Handler) http.Handler {
		return asHandler(func(ctx HTTPContext) { limiteRate(limiter, ctx, next) })
	}
}

func asHandler(handler func(request HTTPContext)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		request := HTTPContext{w, req}
		handler(request)
	})
}

func enableGracefulExit(s *http.Server, done chan struct{}, limiter *Limiter) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	//when interrupt received
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}
	limiter.Stop()

	close(done)
}
