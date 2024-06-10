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

	// create new rate limiter
	rl := middleware.NewRateLimiter(1, 5) // 1 request per second and a burst size of 5
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
	router.Handle("GET /api/notes/search", middleware.AuthMiddleware(http.HandlerFunc(handlers.SearchNote)))

	router.HandleFunc("GET /api/healthchecker", HealthCheckHandler)

	// Custom CORS configuration
	corsConfig := cors.New(cors.Options{
		AllowedHeaders:   []string{"Origin", "Authorization", "Accept", "Content-Type"},
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
		AllowCredentials: true,
	})

	// Wrap the router with the rate limiting middleware
	rateLimitedRouter := rl.RateLimiterMiddleware(router)

	// Wrap the rate limited router with the logRequests middleware
	loggedRouter := logRequests(rateLimitedRouter)

	// Create a new CORS handler
	corsHandler := corsConfig.Handler(loggedRouter)

	server := http.Server{
		Addr:    ":8750",
		Handler: corsHandler,
	}
	err := server.ListenAndServe()
	if err != nil {
		// Check if the error is due to the port already being in use
		if err.Error() == "listen tcp :8750: bind: Only one usage of each socket address (protocol/network address/port) is normally permitted." {
			log.Fatalf("Error: Port 8750 is already in use. Please choose a different port.")
		} else {
			log.Fatal(err)
		}
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
