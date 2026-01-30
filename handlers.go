package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

type UrlShortenerPayload struct {
	Original string `json:"original"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type Metrics struct {
	GlobalTokenBucketCap int `json:"globalTokenBucketCap"`
	GlobalTokensUsed     int `json:"globalTokensUsed"`
	ActiveUsers          int `json:"activeUsers"`
	CurrentUrlCount      int `json:"currentUrlCount"`
}

//--------- Index route -------------------

func MakeIndexHandler() http.Handler {
	fsys := FileHidingFileSystem{http.Dir("./frontend/out/")}
	return http.FileServer(fsys)

}

//------- shortener routes ------------------------

type App struct {
	cfg                  *Config
	shortener            UrlShortener
	globalRateLimiter    *GlobalRateLimiter
	perClientRateLimiter *PerClientRateLimiter
}

var page404HTMLText = Load404Page()

func (app *App) RetrieveUrl(w http.ResponseWriter, r *http.Request) {
	short := r.PathValue("shortUrl")

	if short != "" {
		if original, err := app.shortener.RetrieveUrl(short); err != nil {

			w.Header().Add("Content-Type", "text/html")
			w.WriteHeader(http.StatusNotFound)
			if page404HTMLText != "" {
				fmt.Fprintf(w, "%s", page404HTMLText)
			} else {
				fmt.Fprintf(w, app.cfg.Fallback404HTML)
			}

		} else {
			http.Redirect(w, r, original, http.StatusTemporaryRedirect)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}

}

func (app *App) ShortenUrl(w http.ResponseWriter, r *http.Request) {
	var payload UrlShortenerPayload
	errorMessage := ""

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorMessage = "Oops, we couldn't process your request. Please try again later."
		json.NewEncoder(w).Encode(&ErrorResponse{errorMessage})
		return

	} else if message, ok := ValidateUrl(payload.Original, app.cfg); !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&ErrorResponse{message})
		return
	} else {

		shortUrl, err := Shorten(app.shortener, payload.Original)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errorMessage = "Something broke on our end. Please try again later."
			json.NewEncoder(w).Encode(&ErrorResponse{errorMessage})
			return
		} else {

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			response := map[string]string{
				"shortCode": shortUrl,
			}
			json.NewEncoder(w).Encode(response)
		}
	}

}

func (app *App) StreamMetrics(w http.ResponseWriter, r *http.Request) {

	flusher, ok := w.(http.Flusher)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&ErrorResponse{"Metrics Streaming is currently unsupported."})
		return
	}

	w.Header().Add("Content-Type", "text/event-stream")
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Connection", "keep-alive")

	metricsTicker := time.NewTicker(time.Second)
	defer metricsTicker.Stop()
	errorCount := 0

	for {
		select {
		case <-metricsTicker.C:
			globalTokenBucketCap := app.globalRateLimiter.bucket.Cap()
			globalTokensUsed := globalTokenBucketCap - app.globalRateLimiter.bucket.Len()
			activeUsers := app.perClientRateLimiter.timeLogStore.Len()
			currentUrlCount := app.shortener.Len()

			jsonData, err := json.Marshal(&Metrics{globalTokenBucketCap, globalTokensUsed, activeUsers, currentUrlCount})
			if err != nil {
				errorCount++
				if errorCount > 2 {
					SendSSEErrorEvent(w, "Metrics Streaming is currently unavailable.", flusher)

					return
				}
			}

			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			flusher.Flush()

		case <-r.Context().Done():
			fmt.Println("Client closed connection. Stopping Metrics stream.")
			return
		}
	}

}

func (app *App) StressTest(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&ErrorResponse{"Unexpected error occured while running tests. Please try again later."})
		return
	}

	w.Header().Add("Content-Type", "text/event-stream")
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Connection", "keep-alive")

	if testServer, app, err := StartTestServer(app.cfg); err != nil {
		//TODO log server start failure error message
		SendSSEErrorEvent(w, "Failed to start test server. Please try again later.", flusher)

		return
	} else {
		defer app.shortener.Offline()
		defer app.perClientRateLimiter.Offline()
		defer app.globalRateLimiter.Offline()
		defer testServer.Shutdown(context.Background())

		serverStopUnexpected := make(chan struct{})

		go func() {
			if err := testServer.ListenAndServe(); err != http.ErrServerClosed {
				close(serverStopUnexpected)
			}
		}()

		testCommand := exec.Command("./production_stress_test.sh")
		stdout, err := testCommand.StdoutPipe()
		if err != nil {
			fmt.Println(err)
			SendSSEErrorEvent(w, "Unexpected error occured while running tests. Please try again later.", flusher)
			return
		}

		testCommand.Stderr = testCommand.Stdout

		if err := testCommand.Start(); err != nil {
			fmt.Println(err)

			SendSSEErrorEvent(w, "Unexpected error occured while running tests. Please try again later.", flusher)
			return
		}

		scanner := bufio.NewScanner(stdout)

		for scanner.Scan() {
			select {

			case <-r.Context().Done():
				fmt.Println("Client closed connection. Killing test...")
				testCommand.Process.Kill()
				return

			case <-serverStopUnexpected:
				fmt.Println("Test Server stopped unexpectedly")

				SendSSEErrorEvent(w, "Test server stoped unexpectedly. Please try again later.", flusher)
				testCommand.Process.Kill()
				return

			default:
				jsonData, err := json.Marshal(map[string]string{"outputLine": scanner.Text()})
				if err != nil {
					SendSSEErrorEvent(w, "Unexpected error occured while reading test output. Please try again later.", flusher)
					return
				}

				fmt.Fprintf(w, "data: %s\n\n", jsonData)
				flusher.Flush()
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Test Server stopped unexpectedly")
			SendSSEErrorEvent(w, "Unexpected error occured while reading test output. Please try again later.", flusher)
			return
		}

		if err := testCommand.Wait(); err != nil {
			SendSSEErrorEvent(w, "Unexpected error occured while running tests. Please try again later.", flusher)

			return
		} else {
			fmt.Println("Stress test completed successfully.")
			fmt.Fprintf(w, "event: done\n")
			fmt.Fprintf(w, "data: {\"outputLine\": \"Tests completed successfully.\"}\n\n")
			flusher.Flush()
		}
	}

}
