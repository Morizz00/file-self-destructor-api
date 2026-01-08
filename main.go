package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/Morizz00/self-destruct-share-api/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func main() {
	r := chi.NewRouter()
	
	// Structured logging middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(structuredLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration - restrict to specific origins
	allowedOrigins := getCORSOrigins()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link", "X-File-Name", "X-File-Size", "X-Downloads-Left"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Rate limiting - 100 requests per minute per IP
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"file-self-destruct-api"}`))
	})

	// API routes with stricter rate limiting for uploads
	r.Group(func(r chi.Router) {
		// Upload endpoint: 10 requests per minute per IP
		r.With(httprate.LimitByIP(10, 1*time.Minute)).Post("/upload", handlers.Upload)
		r.Get("/file/{id}", handlers.DownloadFile)
		r.Get("/preview/{id}", handlers.Preview)
		r.Get("/meta/{id}", handlers.GetMeta)
	})

	// Static file serving
	workDir, _ := os.Getwd()
	log.Printf("Working directory: %s", workDir)
	
	// Check if static files exist
	staticFiles := []string{"index.html", "download.html", "styles.css", "script.js"}
	for _, file := range staticFiles {
		path := filepath.Join(workDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Printf("WARNING: Static file not found: %s", path)
		}
	}
	
	fs := http.FileServer(http.Dir(workDir))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(workDir, "index.html"))
	})
	r.Get("/download.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(workDir, "download.html"))
	})
	r.Get("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, filepath.Join(workDir, "styles.css"))
	})
	r.Get("/script.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, filepath.Join(workDir, "script.js"))
	})
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Get port from environment variable or use 8000 as default (Koyeb default)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// structuredLogger provides structured logging middleware
func structuredLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		log.Printf(
			"%s %s %s %d %d %s %s",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			ww.Status(),
			ww.BytesWritten(),
			time.Since(start),
			r.UserAgent(),
		)
	})
}

// getCORSOrigins returns allowed CORS origins from environment or defaults
func getCORSOrigins() []string {
	corsEnv := os.Getenv("CORS_ORIGINS")
	if corsEnv == "" {
		// Default: allow all origins for development (use specific origins in production)
		// Set CORS_ORIGINS environment variable in production
		return []string{"*"}
	}
	
	// Split by comma and trim spaces
	origins := strings.Split(corsEnv, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	return origins
}
