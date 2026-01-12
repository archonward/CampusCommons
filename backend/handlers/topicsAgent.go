package handlers

import (
	//"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/archonward/CampusCommons/backend/database"
)

// A simple struct to start off, more fields can be added into the struct if necessary in the future.
type Topic struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	CreatedBy   int       `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// This func will handle GET requests for the topic of the post.
func GetTopics(w http.ResponseWriter, r *http.Request) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Use the built in Query to find the rows required, error if there is no such topics
	rows, err := database.DB.Query(`
		SELECT id, title, description, created_by, created_at 
		FROM topics 
		ORDER BY created_at DESC
	`)
	if err != nil {
		log.Printf("Database query error: %v", err)
		http.Error(w, "No such topics in database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var topics []Topic		// empty list
	for rows.Next() {		// for each of the item inside rows, create a Topic, then append to the list as required
		var t Topic
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.CreatedBy, &t.CreatedAt)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			http.Error(w, "Data parsing error", http.StatusInternalServerError)
			return
		}
		topics = append(topics, t)
	}

	// Check for iteration errors
	if err = rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		http.Error(w, "Data retrieval error", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	json.NewEncoder(w).Encode(topics)
}
