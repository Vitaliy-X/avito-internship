package app

import (
	"avito-internship/internal/transport"
	"log"
)

func Run() {
	router := transport.SetupRouter()
	log.Println("Starting server on port 8080")

	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
