package migrations

import (
	"log"

	"github.com/Rakhulsr/foodcourt/internal/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	log.Println("Running AutoMigrate...")

	err := db.AutoMigrate(
		&models.Booth{},
		&models.Menu{},
		&models.Order{},
		&models.OrderItem{},
		&models.WhatsAppLog{},
	)

	if err != nil {
		return err
	}

	log.Println("AutoMigrate completed successfully!")
	return nil
}
