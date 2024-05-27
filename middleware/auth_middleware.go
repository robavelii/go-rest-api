package middleware

import (
	"encoding/json"
	"log"
	"os"

	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

func AuthMiddleware(next http.Handler) http.Handler {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	JWT_SECRET := os.Getenv("JWT_SECRET")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get jwt token form the auth header
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Unauthorized! Please login.",
			})
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Invalid token format",
			})
			return
		}

		tokenString := tokenParts[1]
		log.Println("tokenString: ", tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(JWT_SECRET), nil
		})

		if err != nil || !token.Valid {
			log.Println("Token error:", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Invalid token. Unauthorized.",
			})
			return
		}

		// call the next handler
		next.ServeHTTP(w, r)

	})
}
