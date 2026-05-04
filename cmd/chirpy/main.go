package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"chirpy/internal/api"
	"chirpy/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	dbQueries := database.New(db)
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")

	cfg := api.ApiConfig{
		DB:        dbQueries,
		Platform:  platform,
		JwtSecret: jwtSecret,
		PolkaKey:  polkaKey,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /admin/metrics", cfg.MetricsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.ResetHandler)

	mux.HandleFunc("GET /api/healthz", api.ReadinessHandler)

	mux.HandleFunc("POST /api/users", cfg.CreateUserHandler)
	mux.HandleFunc("PUT /api/users", cfg.UpdateUserHandler)
	mux.HandleFunc("POST /api/login", cfg.LoginHandler)

	mux.HandleFunc("POST /api/chirps", cfg.CreateChirpHandler)
	mux.HandleFunc("GET /api/chirps", cfg.GetChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.GetChirpHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.DeleteChirpHandler)

	mux.HandleFunc("POST /api/refresh", cfg.RefreshHandler)
	mux.HandleFunc("POST /api/revoke", cfg.RevokeHandler)

	mux.HandleFunc("POST /api/polka/webhooks", cfg.WebhookHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving file from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
