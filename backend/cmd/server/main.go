package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"kb-backend/internal/handler"
	"kb-backend/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "kb_user"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "kb_password"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "kb_db"
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	log.Printf("Connecting to PostgreSQL at %s:%s/%s...", dbHost, dbPort, dbName)
	pool, err := pgxpool.New(ctx, connStr)
	var dbPool *pgxpool.Pool

	if err == nil {
		if pingErr := pool.Ping(ctx); pingErr == nil {
			log.Println("PostgreSQL connection established successfully.")
			dbPool = pool
			defer dbPool.Close()
		} else {
			log.Printf("PostgreSQL ping failed (%v). Falling back to In-Memory mode.", pingErr)
			pool.Close()
		}
	} else {
		log.Printf("PostgreSQL pool creation failed (%v). Falling back to In-Memory mode.", err)
	}

	if dbPool == nil {
		log.Println("⚡ Running in Standalone In-Memory Mode (Zero-Config DB). All document CRUD, categories, tags, and search operations will work seamlessly in memory!")
	}

	repo := repository.NewRepository(dbPool)
	h := handler.NewHandler(repo)

	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS Configuration
	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	var allowedOrigins []string
	if corsOrigins != "" {
		allowedOrigins = strings.Split(corsOrigins, ",")
	} else {
		allowedOrigins = []string{"*"}
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// API Routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			status := "ok"
			mode := "postgres"
			if repo.IsInMemory() {
				mode = "in-memory"
			}
			_, _ = w.Write([]byte(fmt.Sprintf(`{"status":"%s","mode":"%s","timestamp":"%s"}`, status, mode, time.Now().Format(time.RFC3339))))
		})

		r.Get("/documents", h.GetDocuments)
		r.Post("/documents", h.CreateDocument)
		r.Post("/documents/batch", h.BatchUpload)
		r.Get("/documents/{identifier}", h.GetDocumentByIDOrSlug)
		r.Put("/documents/{id}", h.UpdateDocument)
		r.Delete("/documents/{id}", h.DeleteDocument)

		r.Get("/categories", h.GetCategories)
		r.Post("/categories", h.CreateCategory)

		r.Get("/tags", h.GetTags)
		r.Post("/tags", h.CreateTag)

		r.Get("/search", h.SearchDocuments)
	})

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server listening on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
