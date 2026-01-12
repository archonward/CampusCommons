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

// this new func will handle POST topic requests for people looking to post.
func CreateTopic(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodPost {	// only POST
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Decode JSON request, simple struct with 3 fields, again more could be added later
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		CreatedBy   int    `json:"created_by"` // For now, i require user ID, 
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		log.Printf("Invalid JSON: %v", err)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate to help to prevent bugs, title cannot be empty, userID also cannot
	if input.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if input.CreatedBy <= 0 {
		http.Error(w, "Valid created_by user ID is required", http.StatusBadRequest)
		return
	}

	// Insert into database
	// .Exec is used for any commands that change data: insert, update, delete.
	result, err := database.DB.Exec(`
		INSERT INTO topics (title, description, created_by)
		VALUES (?, ?, ?)	
	`, input.Title, input.Description, input.CreatedBy)	//each of this will be inserted into the 3 placeholders.

	if err != nil {
		log.Printf("Database insert error: %v", err)
		http.Error(w, "Failed to create topic", http.StatusInternalServerError)
		return
	}

	// get the new topic ID that was inserted, result will contain the new row's ID
	topicID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed to get last insert ID: %v", err)
		http.Error(w, "Failed to retrieve topic ID", http.StatusInternalServerError)
		return
	}

	// Fetch the full topic (to return complete object with timestamps)
	var topic Topic
	err = database.DB.QueryRow(`
		SELECT id, title, description, created_by, created_at
		FROM topics
		WHERE id = ?
	`, topicID).Scan(&topic.ID, &topic.Title, &topic.Description, &topic.CreatedBy, &topic.CreatedAt)

	if err != nil {
		log.Printf("Failed to fetch created topic: %v", err)
		http.Error(w, "Failed to retrieve created topic", http.StatusInternalServerError)
		return
	}

	// Return 201 Created + JSON topic
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(topic)
}
