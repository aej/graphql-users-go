package migrate

import (
	_ "database/sql"
	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/lib/pq"
	"log"
)

func Up() {
	migrations_dir := "file://migrations"
	database_url := "postgres://postgres:password@localhost:5432/graphql-users?sslmode=disable"
	log.Printf("Loading migrations from %s", migrations_dir)
	m, err := migrate.New(migrations_dir, database_url)
	if err != nil {
		log.Fatalf("Unable to initialize db migrator. Error: %s", err)
	}
	if err := m.Up(); err != nil {
		log.Fatalf("Up migrations failed. Error: %s", err)
	}
	log.Println("Successfully applied migrations")
}
