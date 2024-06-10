package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

var JWT_SECRET string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	JWT_SECRET = os.Getenv("JWT_SECRET")
}

// Generate jwt token with user id
func GenerateJWT(userID, username, role string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix() //token valid for 2 hour

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(JWT_SECRET)
	// log.Println("SecretKey: ", secretKey)
	return token.SignedString(secretKey)
}

// Verify jwt token
func VerifyJWT(tokenString string) (jwt.MapClaims, error) {
	//parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
