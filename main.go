package main

import (
	"log"
	"net/http"

	"github.com/Morizz00/self-destruct-share-api/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/upload", handlers.Upload)
	r.Get("/file/{id}", handlers.DownloadFile)

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", r)
}
