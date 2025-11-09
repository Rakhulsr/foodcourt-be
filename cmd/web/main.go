package main

import (
	"log"

	"github.com/Rakhulsr/foodcourt/pkg/server"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found or failed to load")
	}

	if err := server.Run(); err != nil {
		log.Fatal("Server failed:", err)
	}
}
