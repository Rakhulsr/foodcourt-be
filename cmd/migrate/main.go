package main

import (
	"log"

	"github.com/Rakhulsr/foodcourt/config"
	"github.com/Rakhulsr/foodcourt/db/migrations"
)

func main() {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	if err := migrations.Migrate(db); err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Migration completed!")
}
