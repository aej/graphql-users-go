package storage

import (
	"github.com/jinzhu/gorm"
	"log"
)

func Connect(databaseUrl string) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatalf("Failed to connect to the db. Error: %s", err)
	}
	log.Print("Database connection established")
	return db, nil
}
