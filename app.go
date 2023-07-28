package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rob-sokolowski/api/lib"
	"io"
	"net/url"

	"net/http"
)

var vellumAiUrl, _ = url.Parse("https://predict.vellum.ai/v1/generate")
var secretsWrapper *lib.SecretsWrapper

func configureMiddleware(ctx context.Context, r *chi.Mux) {
	r.Use(middleware.Logger)

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// TODO: Need to toggle this value based on env
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		//AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-API-KEY"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
}

func configureRouteHandlers(r *chi.Mux, sw *lib.SecretsWrapper) {
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

	//vellumAiProxy := httputil.NewSingleHostReverseProxy(vellumAiUrl)
	//vellumAiProxy.Director = func(req *http.Request) {
	//	req.URL.Scheme = vellumAiUrl.Scheme
	//	req.URL.Host = vellumAiUrl.Host
	//	req.URL.Path = "/v1/generate"
	//	req.Header.Add("X-API-KEY", "<API_KEY>")
	//}

	// TODO: This is not a good solution, but it unblocks me. I would like to use the built-in http proxy
	//       to solve this issue more generally. When I tried that, it appeared to be mutating the request in some way
	//       that wasn't compatible with vellum
	r.Post("/vellum-ai", func(w http.ResponseWriter, r *http.Request) {
		outReq, err := http.NewRequestWithContext(r.Context(), r.Method, vellumAiUrl.String(), r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		apiKey, err := sw.VellumApiKey()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		outReq.Header.Add("X-API-KEY", apiKey)
		outReq.Header.Add("content-type", "application/json")
		outReq.Header.Add("accept", "application/json")

		client := &http.Client{}
		resp, err := client.Do(outReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(resp.StatusCode)
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func main() {
	ctx := context.Background()
	secretsWrapper, err := lib.InitSecretsWrapper(ctx)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	configureMiddleware(ctx, r)
	configureRouteHandlers(r, secretsWrapper)

	fmt.Println("Server starting..")
	http.ListenAndServe(":8080", r)
}
