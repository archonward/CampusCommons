package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/archonward/CampusCommons/backend/database"
)

// Comment represents a forum comment
type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	Body      string    `json:"body"`
	CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

// this func handles GET /posts/{id}/comments
func GetCommentsByPost(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	postIDStr := request.PathValue("id")
	if postIDStr == "" {
		http.Error(writer, "Post ID is required", http.StatusBadRequest)  // if no ID in the path
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		http.Error(writer, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Check if post exists in database
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking post existence: %v", err)
		http.Error(writer, "Database error", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(writer, "Post not found", http.StatusNotFound)
		return
	}

	// Fetch all comments that is under this post, from oldest to the latest
	rows, err := database.DB.Query(`
		SELECT id, post_id, body, created_by, created_at
		FROM comments
		WHERE post_id = ?
		ORDER BY created_at ASC
	`, postID)

	if err != nil {		//if query fails
		log.Printf("Failed to fetch comments: %v", err)
		http.Error(writer, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}
	
	//close rows when we're done to avoid leaking DB resources
	defer rows.Close()

	//var comments []Comment
	comments := []Comment{} // special fix, so that backend does not return null when there are no comments
	for rows.Next() {	// loop through every row once
		var c Comment
		// copy row columns into Comment struct fields
		err := rows.Scan(&c.ID, &c.PostID, &c.Body, &c.CreatedBy, &c.CreatedAt)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			http.Error(writer, "Data parsing error", http.StatusInternalServerError)
			return
		}
		comments = append(comments, c)	// append the comment to the list
	}

	//check for any errors during iteration
	if err = rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		http.Error(writer, "Data retrieval error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(comments)	// send back as JSON
}

// CreateComment handles POST /posts/{id}/comments
func CreateComment(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	postIDStr := request.PathValue("id")
	if postIDStr == "" {
		http.Error(writer, "Post ID is required", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		http.Error(writer, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// validate post exists before allowing comments to be created
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking post existence: %v", err)
		http.Error(writer, "Database error", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(writer, "Post not found", http.StatusNotFound)
		return
	}

	//JSON body shape for creating a comment
	var input struct {
		Body      string `json:"body"`
		CreatedBy int    `json:"created_by"`
	}

	// Decode JSON request body into input struct
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		log.Printf("Invalid JSON: %v", err)
		http.Error(writer, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if input.Body == "" {
		http.Error(writer, "Comment body is required", http.StatusBadRequest)
		return
	}
	if input.CreatedBy <= 0 {
		http.Error(writer, "Valid created_by user ID is required", http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(`
		INSERT INTO comments (post_id, body, created_by)
		VALUES (?, ?, ?)
	`, postID, input.Body, input.CreatedBy)

	if err != nil {
		log.Printf("Failed to create comment: %v", err)
		http.Error(writer, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed to get comment ID: %v", err)
		http.Error(writer, "Failed to retrieve comment ID", http.StatusInternalServerError)
		return
	}

	//fetch the newly created comment 
	var comment Comment
	err = database.DB.QueryRow(`SELECT id, post_id, body, created_by, created_at
		FROM comments
		WHERE id = ?`, commentID).Scan(&comment.ID, &comment.PostID, &comment.Body, &comment.CreatedBy, &comment.CreatedAt)

	if err != nil {
		log.Printf("Failed to fetch created comment: %v", err)
		http.Error(writer, "Failed to retrieve created comment", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)	// return 201 Created
	json.NewEncoder(writer).Encode(comment)
}

