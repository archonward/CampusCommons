package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/archonward/CampusCommons/backend/database"
	"github.com/archonward/CampusCommons/backend/handlers"
	"github.com/rs/cors"
)

func main() {
	
	database.InitDB()
	// I use a ServeMux here so that when we have more routes, the code will not be so confusing
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Backend is running, database connected")
	})
	
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/topics/{id}", handlers.DeleteTopic)
	
	mux.HandleFunc("/posts/{id}/comments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetCommentsByPost(w, r)
		case http.MethodPost:
			handlers.CreateComment(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	
	mux.HandleFunc("/topics/{id}/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetPostsByTopic(w, r)
		case http.MethodPost:
			handlers.CreatePost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/posts/{id}", func(w http.ResponseWriter, r *http.Request) {
    		switch r.Method {
    		case http.MethodGet:
      			handlers.GetPostByID(w, r)
    		case http.MethodDelete:
        		handlers.DeletePost(w, r)
    		default:
        		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    		}
	})	

	//mux.HandleFunc("/topics", handlers.GetTopics) // any request on get Topics handled here
	//mux.HandleFunc("/topics", handlers.CreateTopic)

	mux.HandleFunc("/topics", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method { 	// switch here chooses which option to run based on the value of r.Method
		case http.MethodGet:
			handlers.GetTopics(w, r)
		case http.MethodPost:
			handlers.CreateTopic(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	
	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // React dev server
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		Debug:          false, // may want to set to true to log CORS-related issues
	})
	
	// Wrap the mux with CORS
	handler := c.Handler(mux)

	port := ":8080"
	fmt.Printf("🚀 Server starting on http://localhost%s\n", port)
	fmt.Printf("🔌 CORS enabled for http://localhost:3000\n")
	log.Fatal(http.ListenAndServe(port, handler))
}
