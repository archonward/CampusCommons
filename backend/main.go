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
	
	mux.HandleFunc("/topics", handlers.GetTopics) // any request on Topics handled here.
	
	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // React dev server
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		Debug:          false, // Set to true to log CORS-related issues during dev
	})
	
	// Wrap the mux with CORS
	handler := c.Handler(mux)

	port := ":8080"
	fmt.Printf("🚀 Server starting on http://localhost%s\n", port)
	fmt.Printf("🔌 CORS enabled for http://localhost:3000\n")
	log.Fatal(http.ListenAndServe(port, handler))
}
