package server

import (
	"log"

	"github.com/Rakhulsr/foodcourt/config"
	"github.com/Rakhulsr/foodcourt/pkg/router"
)

func Run() error {

	db, err := config.GetDB()
	if err != nil {
		return err
	}

	r := router.NewRouter(db)

	log.Println("Server starting on :8080")
	return r.Run(":8081")
}
