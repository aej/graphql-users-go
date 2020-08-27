package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DatabaseUrl string
	SecretKey   string
	HashKey     string
	Domain      string
	SecureSessionCookie bool
}

func GetConfig() Config {
	loadDotenv()

	config := Config{}

	config.DatabaseUrl = os.Getenv("DATABASE_URL")
	config.SecretKey = os.Getenv("SECRET_KEY")
	config.HashKey = os.Getenv("HASH_KEY")
	config.Domain = os.Getenv("HOST")

	if os.Getenv("APP_ENV") == "dev" {
		config.SecureSessionCookie = false
	} else {
		config.SecureSessionCookie = true
	}

	return config
}

func loadDotenv() {
	if os.Getenv("APP_ENV") == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}
