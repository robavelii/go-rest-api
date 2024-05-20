package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateNoteSchema

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}
	// validate payload struct
	errors := ValidateStruct(&payload)
	if errors != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errors)
		return
	}

	now := time.Now()
	newNote := Note{
		Title:     payload.Title,
		Content:   payload.Content,
		Category:  payload.Category,
		Published: payload.Published,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// save new note
	result := DB.Create(&newNote)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "UNIQUE constraint field") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  "fail",
				"message": "Title already exists!",
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": result.Error.Error(),
		})
		return
	}

	//Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"note": newNote,
		},
	})

}
