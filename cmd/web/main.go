package main

import (
	"log"

	"github.com/Rakhulsr/foodcourt/pkg/server"
)

func main() {

	if err := server.Run(); err != nil {
		log.Fatal("Server failed:", err)
	}
}
