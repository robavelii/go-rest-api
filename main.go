package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"
)

func init() {
	// init database
	err := ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /api/healthchecker", HealthCheckHandler)
	router.HandleFunc("PATCH /api/notes/{noteId}", UpdateNote)
	router.HandleFunc("GET /api/notes/{noteId}", FindNoteById)
	router.HandleFunc("DELETE /api/notes/{noteId}", DeleteNote)
	router.HandleFunc("POST /api/notes/", CreateNoteHandler)
	router.HandleFunc("GET /api/notes/", FindNotes)

	// Custom CORS configuration
	corsConfig := cors.New(cors.Options{
		AllowedHeaders:   []string{"Origin", "Authorization", "Accept", "Content-Type"},
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
		AllowCredentials: true,
	})

	// Wrap the router with the logRequests middleware
	loggedRouter := logRequests(router)

	// Create a new CORS handler
	corsHandler := corsConfig.Handler(loggedRouter)

	server := http.Server{
		Addr:    ":8750",
		Handler: corsHandler,
	}

	log.Println("Starting server on port: 8080")

	server.ListenAndServe()
}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		elapsed := time.Since(start)
		log.Printf("Received request: %d %s %s %s", wrapped.statusCode, r.Method, r.URL.Path, elapsed)
	})
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "ok",
		"message": "Go server running!",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
