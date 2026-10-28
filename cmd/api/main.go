package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := env("DATABASE_PATH", "studybuddy.db")
	db, err := sql.Open("sqlite3", dbPath+"?_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		log.Fatalf("initialize database: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable", "database": "disconnected"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "database": "connected"})
	})

	handler := withCORS(mux, env("CORS_ORIGINS", "http://localhost:5173,http://localhost:4173,http://127.0.0.1:5173"))

	addr := env("API_ADDR", ":8080")
	server := &http.Server{Addr: addr, Handler: handler, ReadHeaderTimeout: 5 * time.Second}
	log.Printf("StudyBuddy API listening on %s", addr)
	log.Fatal(server.ListenAndServe())
}

func withCORS(next http.Handler, allowedOrigins string) http.Handler {
	allowed := make(map[string]bool)
	for _, origin := range strings.Split(allowedOrigins, ",") {
		if origin = strings.TrimSpace(origin); origin != "" {
			allowed[origin] = true
		}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && allowed[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}
		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write response: %v", err)
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
