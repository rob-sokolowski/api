package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rob-sokolowski/cors-proxy/lib"

	"net/http"
)

func configureMiddleware(ctx context.Context, r *chi.Mux) {
	r.Use(middleware.Logger)
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// TODO: Need to toggle this value based on env
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
}

func configureRouteHandlers(r *chi.Mux) {
	// healthchecks
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		data := lib.PingResponse{
			Message: "PONG!",
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	r.Post("/vellum-ai", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}

func main() {
	r := chi.NewRouter()
	configureMiddleware(context.Background(), r)
	configureRouteHandlers(r)

	fmt.Println("Server starting..")
	http.ListenAndServe(":8080", r)
}
