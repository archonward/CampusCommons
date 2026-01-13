package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/archonward/CampusCommons/backend/database"
)

// just ID and Username for now, i may add 2 more fields, createdTime and whether the user is active or not. Depends on time
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// this func will handle POST /login, for now, we check if username is empty to prevent bugs.
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// parse request body
	var input struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("invalid JSON (login request): %v", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return		// if decode fails return void
	}

	username := input.Username

	// guard against empty inputs for username
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	var existingUser User
	err := database.DB.QueryRow(`
		SELECT id, username FROM users WHERE username = ?
	`, username).Scan(&existingUser.ID, &existingUser.Username)

	switch {
	case err == sql.ErrNoRows:	// user do not exist, so i create a new User
		result, err := database.DB.Exec(`
			INSERT INTO users (username) VALUES (?)
		`, username)
		if err != nil {
			log.Printf("Failed to create user: %v", err)
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}

		userID, _ := result.LastInsertId()
		existingUser = User{ID: int(userID), Username: username}

	case err != nil:
		// for all other errors
		log.Printf("login error: %v", err)
		http.Error(w, "Login failed", http.StatusInternalServerError)
		return

	default:
	}

	// Return user object
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingUser)
}
