package main

import (
	"log"

	"github.com/Rakhulsr/foodcourt/config"
	"github.com/Rakhulsr/foodcourt/internal/model"
)

func main() {
	db, err := config.GetDB()
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	booths := []model.Booth{
		// {Name: "Booth Makanan Pak Joko", WhatsApp: "6281234567890", IsActive: true},
		// {Name: "Booth Minuman Bu Rini", WhatsApp: "6289876543210", IsActive: true},
		{Name: "Kedai Suka-Suki", WhatsApp: "62898765432320", IsActive: false},
	}
	for _, b := range booths {
		db.Create(&b)
	}

	menus := []model.Menu{
		// {BoothID: 1, Name: "Nasi Goreng", Price: 25000, Category: "makanan", IsAvailable: true, ImagePath: "/uploads/menu/nasi.jpg"},
		// {BoothID: 1, Name: "Ayam Geprek", Price: 20000, Category: "makanan", IsAvailable: true, ImagePath: "/uploads/menu/geprek.jpg"},
		// {BoothID: 2, Name: "Es Teh", Price: 5000, Category: "minuman", IsAvailable: true, ImagePath: "/uploads/menu/esteh.jpg"},
		{BoothID: 6, Name: "Steak Ayam Kampus", Price: 350000, Category: "makanan", IsAvailable: false, ImagePath: "/uploads/menu/esteh.jpg"},
	}
	for _, m := range menus {
		db.Create(&m)
	}

	log.Println("Seeding completed! 2 booth, 3 menu.")
}
