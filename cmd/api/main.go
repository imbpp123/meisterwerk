package main

import (
	"log"
	"net/http"

	"app/internal/quote"
)

func main() {
	router := quote.RouterAPIInitializer()
	if router == nil {
		log.Fatal("Failed to initialize router")
	}

	address := ":8080"
	log.Printf("Starting server on %s...", address)
	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
