package routes

import (
	"context"
	"embed"
	"net/http"
	"rater/db"
	"rater/frontend"
	"rater/utils"

	"github.com/go-http-utils/etag"
)

var ctx = context.Background()
var assets embed.FS
var executor *db.Queries

func CreateRouter(e *db.Queries, a embed.FS) *http.ServeMux {
	executor = e
	assets = a

	router := http.NewServeMux()

	router.HandleFunc("GET /", utils.Logger(home))
	router.HandleFunc("GET /sign/up/", utils.Logger(signUp))
	router.HandleFunc("GET /sign/in/", utils.Logger(signIn))
	router.Handle("GET /assets/", (etag.Handler(utils.Gzip(http.FileServer(http.FS(assets))), false)))

	return router
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	frontend.Page().Render(ctx, w)
}
func signUp(w http.ResponseWriter, _ *http.Request) {
	frontend.Sign("/sign/up/", "Create new account", "Enter your details below to get started", "Log In", "Already have an account?", "/sign/in").Render(ctx, w)
}
func signIn(w http.ResponseWriter, _ *http.Request) {
	frontend.Sign("/sign/in/", "Login to your account", "Enter your details below to continue", "Sign Up", "Don't have an account?", "/sign/up").Render(ctx, w)
}

// var static =
