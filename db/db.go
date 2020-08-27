package db

import (
	"github.com/jinzhu/gorm"
	"log"
)

func Connect() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", "postgres://postgres@localhost:5432/graphql-users?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to the db. Error: %s", err)
	}

	log.Print("Database connection established")

	return db, nil
}