package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/archonward/CampusCommons/backend/database"
	"strconv"
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
func GetTopics(writer http.ResponseWriter, request *http.Request) {
	// Set content type
	writer.Header().Set("Content-Type", "application/json")

	// Use the built in Query to find the rows required, error if there is no such topics
	rows, err := database.DB.Query(`
		SELECT id, title, description, created_by, created_at 
		FROM topics 
		ORDER BY created_at DESC
	`)
	if err != nil {
		log.Printf("Database query error: %v", err)
		http.Error(writer, "No such topics in database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var topics []Topic		// empty list
	for rows.Next() {		// for each of the item inside rows, create a Topic, then append to the list as required
		var t Topic
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.CreatedBy, &t.CreatedAt)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			http.Error(writer, "Data parsing error", http.StatusInternalServerError)
			return
		}
		topics = append(topics, t)
	}

	// Check for iteration errors
	if err = rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		http.Error(writer, "Data retrieval error", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	json.NewEncoder(writer).Encode(topics)
}

// this new func will handle POST topic requests for people looking to post.
func CreateTopic(writer http.ResponseWriter, request *http.Request) {
	
	if request.Method != http.MethodPost {	// only POST
		http.Error(writer, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set content type
	writer.Header().Set("Content-Type", "application/json")

	// Decode JSON request, simple struct with 3 fields, again more could be added later
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		CreatedBy   int    `json:"created_by"` // For now, i require user ID, 
	}

	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&input); err != nil {
		log.Printf("Invalid JSON: %v", err)
		http.Error(writer, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate to help to prevent bugs, title cannot be empty, userID also cannot
	if input.Title == "" {
		http.Error(writer, "Title is required", http.StatusBadRequest)
		return
	}
	if input.CreatedBy <= 0 {
		http.Error(writer, "Valid created_by user ID is required", http.StatusBadRequest)
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
		http.Error(writer, "Failed to create topic", http.StatusInternalServerError)
		return
	}

	// get the new topic ID that was inserted, result will contain the new row's ID
	topicID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed to get last insert ID: %v", err)
		http.Error(writer, "Failed to retrieve topic ID", http.StatusInternalServerError)
		return
	}

	// Fetch the full topic (to return complete object with timestamps)
	var topic Topic
	err = database.DB.QueryRow(`SELECT id, title, description, created_by, created_at
		FROM topics
		WHERE id = ?`, topicID).Scan(&topic.ID, &topic.Title, &topic.Description, &topic.CreatedBy, &topic.CreatedAt)

	if err != nil {
		log.Printf("Failed to fetch created topic: %v", err)
		http.Error(writer, "Failed to retrieve created topic", http.StatusInternalServerError)
		return
	}

	// Return 201 Created + JSON topic
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(topic)
}

// this func handles DELETE /topics/{id}
func DeleteTopic(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodDelete {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	topicIDStr := request.PathValue("id")
	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil || topicID <= 0 {
		http.Error(writer, "invalid topic ID", http.StatusBadRequest)
		return
	}

	// start off by delete all comments under posts in this topic
	_, err = database.DB.Exec(`DELETE FROM comments 
		WHERE post_id IN (SELECT id FROM posts WHERE topic_id = ?)`, topicID)
	if err != nil {
		log.Printf("Failed to delete comments: %v", err)
		http.Error(writer, "Failed to delete associated comments", http.StatusInternalServerError)
		return
	}

	// delete all posts in this topic
	_, err = database.DB.Exec(`
		DELETE FROM posts 
		WHERE topic_id = ?
	`, topicID)
	if err != nil {
		log.Printf("failed to delete posts: %v", err)
		http.Error(writer, "failed to delete associated posts", http.StatusInternalServerError)
		return
	}

	// here we delete the topic
	result, err := database.DB.Exec(`DELETE FROM topics 
		WHERE id = ?`, topicID)
	if err != nil {
		log.Printf("failed to delete topic: %v", err)
		http.Error(writer, "failed to delete topic", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(writer, "Topic not found", http.StatusNotFound)
		return
	}

	writer.WriteHeader(http.StatusNoContent) // status 204
}

// this func handles PUT PUT /topics/{id} requests
func UpdateTopic(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPut {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	topicIDStr := request.PathValue("id")
	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil || topicID <= 0 {
		http.Error(writer, "invalid topic ID", http.StatusBadRequest)
		return
	}

	// Ensure topic exists
	var existingTitle string
	err = database.DB.QueryRow("SELECT title FROM topics WHERE id = ?", topicID).Scan(&existingTitle)
	if err == sql.ErrNoRows {
		http.Error(writer, "topic not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("DB error checking topic: %v", err)
		http.Error(writer, "server error", http.StatusInternalServerError)
		return
	}

	var input struct {							 
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {		// Parsing
		http.Error(writer, "invalid JSON", http.StatusBadRequest)
		return
	}
	if input.Title == "" {
		http.Error(writer, "title is required", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec(`UPDATE topics 
		SET title = ?, description = ? 
		WHERE id = ?`, input.Title, input.Description, topicID)		// Update row
	if err != nil {
		log.Printf("failed to update topic: %v", err)
		http.Error(writer, "failed to update topic", http.StatusInternalServerError)
		return
	}

	var updatedTopic Topic
	err = database.DB.QueryRow(`SELECT id, title, description, created_by, created_at
		FROM topics
		WHERE id = ?`, topicID).Scan(&updatedTopic.ID, &updatedTopic.Title, &updatedTopic.Description, &updatedTopic.CreatedBy, &updatedTopic.CreatedAt)	// Return updated topic

	if err != nil {
		log.Printf("failed to fetch updated topic: %v", err)
		http.Error(writer, "failed to retrieve updated topic", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(updatedTopic)
}

