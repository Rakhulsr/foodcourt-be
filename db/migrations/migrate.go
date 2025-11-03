package migrations

import (
	"log"
	"time"

	"github.com/Rakhulsr/foodcourt/internal/models"
	"gorm.io/gorm"
)

type Migration struct {
	ID        uint `gorm:"primaryKey"`
	Version   string
	AppliedAt time.Time
}

func Migrate(db *gorm.DB) error {

	if !db.Migrator().HasTable(&Migration{}) {
		log.Println("First time setup: Running AutoMigrate...")
		if err := db.AutoMigrate(
			&models.Booth{},
			&models.Menu{},
			&models.Order{},
			&models.OrderItem{},
			&models.WhatsAppLog{},
			&Migration{},
		); err != nil {
			return err
		}

		db.Create(&Migration{Version: "v1.0.0"})
		log.Println("AutoMigrate completed + migration tracked!")
		return nil
	}

	log.Println("Database already initialized. Skipping AutoMigrate.")

	return nil
}
