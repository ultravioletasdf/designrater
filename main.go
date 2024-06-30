package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"

	"rater/frontend"
	"rater/utils"

	"github.com/go-http-utils/etag"
)

//go:embed assets/*
var assets embed.FS
var ctx = context.Background()

func main() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		frontend.Page().Render(ctx, w)
	})
	router.Handle("GET /assets/", etag.Handler(utils.Gzip(http.FileServer(http.FS(assets))), false))

	server := http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	fmt.Println("Starting server at http://127.0.0.1:3000")
	server.ListenAndServe()
}
