package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/archonward/CampusCommons/backend/database"
)

type Post struct {
	ID       int       `json:"id"`
	TopicID  int       `json:"topic_id"`
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	CreatedBy int      `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

// this func handles GET /topics/id/posts
func GetPostsByTopic(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	// extract topic ID from URL path
	topicIDStr := request.PathValue("id")
	if topicIDStr == "" {
		http.Error(writer, "topic ID is required", http.StatusBadRequest)
		return
	}

	topicID, err := strconv.Atoi(topicIDStr)	// convert to int
	if err != nil || topicID <= 0 {
		http.Error(writer, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	// Check if topic exists first
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM topics WHERE id = ?)", topicID).Scan(&exists)
	// standard SQL query to check for presence of topic
	if err != nil {
		log.Printf("Error checking topic existence: %v", err)
		http.Error(writer, "Database error", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(writer, "Topic not found", http.StatusNotFound)	//send to react
		return
	}

	// Fetch all posts under the topic, starting from the oldest
	rows, err := database.DB.Query(`SELECT id, topic_id, title, body, created_by, created_at
		FROM posts
		WHERE topic_id = ?
		ORDER BY created_at ASC`, topicID)

	if err != nil {
		log.Printf("fail to fetch posts: %v", err)
		http.Error(writer, "failed to fetch posts", http.StatusInternalServerError)	//send to react
		return
	}
	// close rows when done to avoid leaking DB resources
	defer rows.Close()

	postList := []Post{}	
	for rows.Next() {	// loop runs once per row, then moves onto next
		var p Post
		err := rows.Scan(&p.ID, &p.TopicID, &p.Title, &p.Body, &p.CreatedBy, &p.CreatedAt)
		if err != nil {
			// if scanning fails, return 500
			log.Printf("row error: %v", err)
			http.Error(writer, "parsing error", http.StatusInternalServerError)	// send the http error back to react
			return
		}
		postList = append(postList, p)	// add the Post to the list
	}
	
	err = rows.Err();
	if err != nil {
		log.Printf("rows iteration error: %v", err)
		http.Error(writer, "retrieval error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(postList)	// converts the list of Post into JSON, writes direct to ResponseWriter
}

func GetPostByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")	// tell client response will be JSON

	postIDStr := r.PathValue("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var p Post
	//query for a single post by ID
	err = database.DB.QueryRow(`
		SELECT id, topic_id, title, body, created_by, created_at
		FROM posts
		WHERE id = ?
	`, postID).Scan(&p.ID, &p.TopicID, &p.Title, &p.Body, &p.CreatedBy, &p.CreatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Failed to fetch post: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(p)
}

// this func handles POST /topics/id/posts
func CreatePost(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	// Extract topic ID from URL
	topicIDStr := request.PathValue("id")
	if topicIDStr == "" {
		http.Error(writer, "topic ID is required", http.StatusBadRequest)
		return
	}

	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil || topicID <= 0 {
		http.Error(writer, "invalid ID", http.StatusBadRequest)
		return
	}

	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM topics WHERE id = ?)", topicID).Scan(&exists)
	if err != nil {
		log.Printf("error checking existence of topic: %v", err)
		http.Error(writer, "database error", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(writer, "topic not found", http.StatusNotFound)
		return
	}

	// expected JSON body for creeating post
	var input struct {
		Title     string `json:"title"`
		Body      string `json:"body"`
		CreatedBy int    `json:"created_by"`
	}

	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		log.Printf("invalid JSON: %v", err)
		http.Error(writer, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	if input.Title == "" {
		http.Error(writer, "title can't be empty", http.StatusBadRequest)
		return
	}
	if input.Body == "" {
		http.Error(writer, "body can't be empty", http.StatusBadRequest)
		return
	}
	if input.CreatedBy <= 0 {
		http.Error(writer, "valid created_by user ID is required", http.StatusBadRequest)
		return
	}

	// insert post
	result, err := database.DB.Exec(`INSERT INTO posts (topic_id, title, body, created_by)
		VALUES (?, ?, ?, ?)`, topicID, input.Title, input.Body, input.CreatedBy) //SQL INSERT to create a new row

	if err != nil {
		log.Printf("Failed to create post: %v", err)
		http.Error(writer, "Failed to create post", http.StatusInternalServerError)
		return
	}
	
	// get the auto generated post ID from database
	postID, err := result.LastInsertId()
	if err != nil {
		log.Printf("failed to get post ID: %v", err)
		http.Error(writer, "failed to get post ID", http.StatusInternalServerError)
		return
	}

	// Read the inserted post back from the database so we can return it as JSON.
	var post Post
	row := database.DB.QueryRow(`
		SELECT id, topic_id, title, body, created_by, created_at
		FROM posts
		WHERE id = ?
	`, postID)	// using WHERE id = ? so that we will take in the row with all the fields updated by the database
	err = row.Scan(&post.ID, &post.TopicID, &post.Title, &post.Body, &post.CreatedBy, &post.CreatedAt)

	if err != nil {
		log.Printf("failed to fetch created post: %v", err)
		http.Error(writer, "failed to retrieve created post", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(post)
}

// this func handles DELETE /posts/{id}
func DeletePost(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodDelete {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	postIDStr := request.PathValue("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		http.Error(writer, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// start off, delete all comments under this post
	_, err = database.DB.Exec(`DELETE FROM comments 
		WHERE post_id = ?`, postID)
	if err != nil {
		log.Printf("failed to delete comments: %v", err)
		http.Error(writer, "failed to delete comments", http.StatusInternalServerError)
		return
	}

	// Then, delete the post
	result, err := database.DB.Exec(`DELETE FROM posts 
		WHERE id = ?`, postID)
	if err != nil {
		log.Printf("failed to delete post: %v", err)
		http.Error(writer, "failed to delete post", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(writer, "Post does not exist", http.StatusNotFound)
		return
	}

	writer.WriteHeader(http.StatusNoContent) // 204
}

