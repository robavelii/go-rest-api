package main

import (
	"encoding/json"
	"example/rest-api/db"
	"example/rest-api/handlers"
	"example/rest-api/middleware"
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"
)

func init() {
	// init database
	err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

}

func main() {
	router := http.NewServeMux()

	// auth routes
	router.Handle("POST /api/auth/register", http.HandlerFunc(handlers.RegisterHandler))
	router.Handle("POST /api/auth/login", http.HandlerFunc(handlers.LoginHandler))
	router.Handle("POST /api/auth/logout", middleware.AuthMiddleware(http.HandlerFunc(handlers.LogoutHandler)))

	// note routes
	router.Handle("PATCH /api/notes/{noteId}", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateNote)))
	router.Handle("GET /api/notes/{noteId}", middleware.AuthMiddleware(http.HandlerFunc(handlers.FindNoteById)))
	router.Handle("DELETE /api/notes/{noteId}", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteNote)))
	router.Handle("POST /api/notes/", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateNoteHandler)))
	router.Handle("GET /api/notes/", middleware.AuthMiddleware(http.HandlerFunc(handlers.FindNotes)))

	router.HandleFunc("GET /api/healthchecker", HealthCheckHandler)

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

	log.Println("Starting server on port: 8750")

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
