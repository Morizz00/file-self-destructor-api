package main

import (
	"log"
	"net/http"

	"github.com/Morizz00/self-destruct-share-api/blockchain/handlers"
	"github.com/Morizz00/self-destruct-share-api/blockchain/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CORS())

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Blockchain API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/proof/upload", handlers.RegisterUpload)
		r.Post("/proof/download", handlers.RegisterDownload)
		r.Get("/proof/verify/{file_id}", handlers.VerifyProof)
		r.Get("/proof/{file_id}", handlers.GetProof)
	})

	log.Printf("Starting blockchain service on %s", cfg.BindAddress)
	log.Fatal(http.ListenAndServe(cfg.BindAddress, r))
}
