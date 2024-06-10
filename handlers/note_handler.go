package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"example/rest-api/db"
	"example/rest-api/models"

	"gorm.io/gorm"
)

// ! CREATE
func CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	var payload models.CreateNoteSchema

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
	errors := models.ValidateStruct(&payload)
	if errors != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errors)
		return
	}

	now := time.Now()
	newNote := models.Note{
		Title:     payload.Title,
		Content:   payload.Content,
		Category:  payload.Category,
		Published: payload.Published,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// save new note
	result := db.DB.Create(&newNote)
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

// ! GET ALL
func FindNotes(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}

	intPage, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, "Invalid page parameter", http.StatusBadRequest)
		return
	}
	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}
	offset := (intPage - 1) * intLimit

	var notes []models.Note
	results := db.DB.Limit(intLimit).Offset(offset).Find(&notes)
	if results.Error != nil {
		http.Error(w, results.Error.Error(), http.StatusBadGateway)
		return
	}

	// Return success response
	response := map[string]interface{}{
		"status":  "success",
		"results": len(notes),
		"notes":   notes,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ! GET ONE
func FindNoteById(w http.ResponseWriter, r *http.Request) {
	noteID := r.PathValue("noteId")

	var note models.Note
	result := db.DB.First(&note, "id = ?", noteID)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			response := map[string]interface{}{
				"status":  "fail",
				"message": "No note with that ID exists",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		response := map[string]interface{}{
			"status":  "fail",
			"message": err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"note": note,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func SearchNote(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	var notes []models.Note

	if query != "" {
		db.DB.Where("title ILIKE ? OR content ILIKE ?", "%"+query+"%", "%"+query+"%").Find(&notes)
	} else {
		db.DB.Find(&notes)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// ! PUT
func UpdateNote(w http.ResponseWriter, r *http.Request) {
	noteID := r.PathValue("noteId")

	var payload models.UpdateNoteSchema
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var note models.Note
	result := db.DB.First(&note, "id = ?", noteID)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			errorResponse := map[string]interface{}{
				"status":  "fail",
				"message": "No note with that ID exists",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	updates := make(map[string]interface{})
	if payload.Title != "" {
		updates["title"] = payload.Title
	}
	if payload.Category != "" {
		updates["category"] = payload.Category
	}
	if payload.Content != "" {
		updates["content"] = payload.Content
	}
	if payload.Published != nil {
		updates["published"] = payload.Published
	}
	updates["updated_at"] = time.Now()

	db.DB.Model(&note).Updates(updates)

	response := map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"note": note,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ! DELETE
func DeleteNote(w http.ResponseWriter, r *http.Request) {
	noteID := r.PathValue("noteId")

	result := db.DB.Delete(&models.Note{}, "id = ?", noteID)

	if result.RowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		response := map[string]interface{}{
			"status":  "fail",
			"message": "No note with that ID exists",
		}
		json.NewEncoder(w).Encode(response)
		return
	} else if result.Error != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		response := map[string]interface{}{
			"status":  "error",
			"message": result.Error.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status":  "success",
		"message": "Note deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}
