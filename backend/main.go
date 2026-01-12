package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Backend is running!")
	})

	fmt.Println("🚀 Backend server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
