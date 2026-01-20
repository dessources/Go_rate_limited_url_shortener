package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {

	// initialize bucket
	bucket, err := NewMemoryBucket(50000, 50000)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// initialize limiter
	limiter := NewLimiter(10000, bucket)

	limiterMiddleware := createLimiterMiddleware(limiter)
	http.Handle("/", limiterMiddleware(indexHandler()))
	http.HandleFunc("/stop", func(w http.ResponseWriter, req *http.Request) {
		limiter.Stop()
		os.Exit(0)
	})

	http.ListenAndServe(":8090", nil)

}

func indexHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "This request has been allowed.\n")
	})

}

// func stop(w http.ResponseWriter, req *http.Request) {
// 	limiter.Stop()
// }

func createLimiterMiddleware(l *Limiter) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// fmt.Println("The size is ", r.ContentLength)

			if l.Allow(2) {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte("Request rejected by rate limiter"))
				return
			}
		})
	}

}
