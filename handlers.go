package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UrlShortenerPayload struct {
	Original string `json:"original"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type Metrics struct {
	TokenBucketCap  int `json:"tokenBucketCap"`
	TokenUsed       int `json:"tokenUsed"`
	ActiveUsers     int `json:"activeUsers"`
	CurrentUrlCount int `json:"currentUrlCount"`
}

//--------- Index route -------------------

func MakeIndexHandler() http.Handler {
	fsys := FileHidingFileSystem{http.Dir("./frontend/out/")}
	return http.FileServer(fsys)

}

//------- shortener routes ------------------------

type App struct {
	shortener        UrlShortener
	globalLimiter    *GlobalLimiter
	perClientLimiter *PerClientLimiter
}

func (app *App) RetrieveUrl(w http.ResponseWriter, r *http.Request) {
	short := r.PathValue("shortUrl")

	if short != "" {
		if original, err := app.shortener.RetrieveUrl(short); err != nil {

			//setting this header makes Go warn me that ServeFile tries to set
			//status again to 200 internally but fails silently.
			//TODO: load 404.html then return text/html with status 404 instead of ServeFile
			w.WriteHeader(http.StatusNotFound)
			http.ServeFile(w, r, "frontend/out/404.html")

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

	} else if message, ok := ValidateUrl(payload.Original); !ok {
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
			tokenBucketCap := app.globalLimiter.bucket.Cap()
			tokenUsed := tokenBucketCap - app.globalLimiter.bucket.Len()
			activeUsers := app.perClientLimiter.timeLogStore.Len()
			currentUrlCount := app.shortener.Len()

			jsonData, err := json.Marshal(&Metrics{tokenBucketCap, tokenUsed, activeUsers, currentUrlCount})
			if err != nil {
				errorCount++
				if errorCount > 2 {
					json.NewEncoder(w).Encode(&ErrorResponse{"Metrics Streaming is currently unsupported."})
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
