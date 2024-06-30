package main

import (
	"embed"
	"fmt"
	"net/http"
	"time"

	"rater/db"
	"rater/routes"
	"rater/utils"

	"github.com/rs/zerolog"
)

//go:embed assets/*
var assets embed.FS

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.DurationFieldUnit = time.Millisecond

	sqlite, err := utils.ConnectToSqlite()
	if err != nil {
		panic(err)
	}
	executor := db.New(sqlite)

	router := routes.CreateRouter(executor, assets)
	server := http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	fmt.Println("Starting server at http://127.0.0.1:3000")
	server.ListenAndServe()
}
