package main

import (
	"log"
	"net/http"
)

func LimiteRate(limiter *Limiter, ctx HTTPContext, next http.Handler) {
	if limiter.Allow(1) {
		next.ServeHTTP(ctx.w, ctx.req)
	} else {
		ctx.w.WriteHeader(http.StatusTooManyRequests)
		_, _ = ctx.w.Write([]byte("Request rejected by rate limiter"))
		return
	}
}

// middleware utils
func createLimiterMiddleware() (*Limiter, LimiterMiddleware) {
	var limiter *Limiter
	l, err := NewLimiter(InMemory, 50000, 50000, 10000)
	if err != nil {
		log.Fatal(err)
	}

	return l, func(next http.Handler) http.Handler {
		return AsHandler(func(ctx HTTPContext) { LimiteRate(limiter, ctx, next) })
	}
}

func AsHandler(handler RouteHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := HTTPContext{w, req}
		handler(ctx)
	})
}
