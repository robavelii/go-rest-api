package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"example/rest-api/db"
	"example/rest-api/models"
	"example/rest-api/utils"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var payload models.CreateUserSchema

	// decode request body
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate payload
	errors := models.ValidateStruct(&payload)
	if errors != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}

	//hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create new user
	newUser := models.User{
		Username: payload.Username,
		Email:    payload.Email,
		FullName: payload.FullName,
		Password: string(hashedPassword),
		Role:     payload.Role,
	}

	// save the user
	if err := db.DB.Create(&newUser).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User created successfully",
		"data":    newUser,
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username *string `json:"username"`
		Email    string  `json:"email"`
		Password string  `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid email or password", http.StatusBadRequest)
		return
	}

	// find the user by email or username
	var user models.User
	if err := db.DB.Where("email = ? OR username = ?", credentials.Email, credentials.Username).First(&user).Error; err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	//verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// generate new jwt token
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   token,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get the JWT token from the request
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// No token provided
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "fail",
			"message": "No token provided",
		})
		return
	}

	// Parse and invalidate the token
	tokenString := strings.Split(authHeader, " ")[1]
	claims := &jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})

	if err != nil {
		// Invalid token
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "fail",
			"message": "Invalid token",
		})
		return
	}

	// Invalidate the token by setting its expiration time to the past
	(*claims)["exp"] = time.Now().Unix() - 1

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Logged out successfully",
	})
}
