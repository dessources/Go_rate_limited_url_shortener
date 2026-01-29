package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/cors"
)

//---------------Middleware utils ----------------

func SetupCors(mux *http.ServeMux) http.Handler {

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://appurl.com", "http://localhost:8090", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-API-Key"},
		AllowCredentials: true,
		Debug:            false,
	})

	return c.Handler(mux)
}

func ComposeMiddlewares(r ...Middleware) Middleware {
	return func(h http.Handler) http.Handler {
		//start by wrapping handler with last middleware
		acc := r[len(r)-1](h)

		//then compose from last to first
		for i := len(r) - 2; i >= 0; i-- {
			acc = r[i](acc)
		}
		return acc
	}
}

//------------- Handler utils----------------------

func containsTxtFile(name string) bool {
	parts := strings.Split(name, "/")
	for _, part := range parts {
		if strings.HasSuffix(part, ".txt") {
			return true
		}
	}
	return false
}

type FileHidingFileSystem struct {
	http.FileSystem
}

type FileHidingFile struct {
	http.File
}

func (fsys FileHidingFileSystem) Open(name string) (http.File, error) {
	if containsTxtFile(name) {
		// If txt file, return 403 error
		return nil, fs.ErrPermission
	}

	file, err := fsys.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	return FileHidingFile{file}, nil
}

func ValidateUrl(s string) (string, bool) {
	if len(s) > maxUrlLength {
		return fmt.Sprintf("Provided url exceeds max-length of %d", maxUrlLength), false
	}

	u, err := url.Parse(s)

	if err != nil || u.Host == "" {
		return "Invalid url provided", false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return "Invalid protocol provided. Only http:// or https:// allowed", false
	}

	return "", true
}
