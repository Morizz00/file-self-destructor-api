package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Morizz00/self-destruct-share-api/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link", "X-File-Name", "X-File-Size", "X-Downloads-Left"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Post("/upload", handlers.Upload)
	r.Get("/file/{id}", handlers.DownloadFile)
	r.Get("/preview/{id}", handlers.Preview)
	r.Get("/meta/{id}", handlers.GetMeta)
	fs := http.FileServer(http.Dir("."))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	r.Get("/download.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "download.html")
	})
	r.Get("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "styles.css")
	})
	r.Get("/script.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "script.js")
	})
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Get port from environment variable or use 8080 as default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on :%s", port)
	http.ListenAndServe(":"+port, r)
}
