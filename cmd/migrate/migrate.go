package main

import (
	_ "database/sql"
	"github.com/andyjones11/graphql-users/internal/config"
	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func Up() {
	configuration := config.GetConfig()
	migrations_dir := "file://migrations"

	log.Printf("Loading migrations from %s", migrations_dir)

	m, err := migrate.New(migrations_dir, configuration.DatabaseUrl)

	if err != nil {
		log.Fatalf("Unable to initialize db migrator. Error: %s", err)
	}
	if err := m.Up(); err != nil {
		log.Fatalf("Up migrations failed. Error: %s", err)
	}
	log.Println("Successfully applied migrations")
}

func main() {
	if len(os.Args) > 1 {
		migrate_db := os.Args[1]
		if migrate_db == "up" {
			Up()
			return
		}
	}
}
