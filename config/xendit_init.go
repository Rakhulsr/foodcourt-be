package config

import (
	"os"

	"github.com/joho/godotenv"
	xendit "github.com/xendit/xendit-go/v7"
)

var XenditClient *xendit.APIClient

func InitXendit() {
	godotenv.Load()

	apiKey := os.Getenv("XENDIT_SECRET_KEY")
	if apiKey == "" {
		panic("XENDIT_SECRET_KEY is missing in .env")
	}

	XenditClient = xendit.NewClient(apiKey)
}
