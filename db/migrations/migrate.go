package migrations

import (
	"log"
	"time"

	"github.com/Rakhulsr/foodcourt/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type MigrationRecord struct {
	ID        uint      `gorm:"primaryKey"`
	Version   string    `gorm:"size:20;uniqueIndex"`
	AppliedAt time.Time `gorm:"autoCreateTime"`
}

func Migrate(db *gorm.DB) error {
	if db.Migrator().HasTable(&MigrationRecord{}) {
		log.Println("Database already initialized. Skipping AutoMigrate.")
		return nil
	}

	log.Println("First time setup: Running AutoMigrate...")
	err := db.AutoMigrate(
		&model.Booth{},
		&model.Menu{},
		&model.Order{},
		&model.OrderItem{},
		&model.WhatsAppLog{},
		&model.Admin{},
		&MigrationRecord{},
	)
	if err != nil {
		return err
	}

	if err := seedAdmin(db); err != nil {
		log.Printf("Warning: Failed to seed admin: %v", err)
	} else {
		log.Println("Admin seeded: username=admin, password=admin123")
	}

	db.Create(&MigrationRecord{Version: "v1.0.0"})
	log.Println("AutoMigrate completed! Database ready.")
	return nil
}

func seedAdmin(db *gorm.DB) error {
	var count int64
	db.Model(&model.Admin{}).Count(&count)
	if count > 0 {
		return nil
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := model.Admin{
		Username: "admin",
		Password: string(hashed),
		IsActive: true,
	}

	return db.Create(&admin).Error
}
