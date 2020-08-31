package main

import (
	"github.com/andyjones11/graphql-users/api"
	migrate "github.com/andyjones11/graphql-users/bin"
	"github.com/andyjones11/graphql-users/db"
	"github.com/andyjones11/graphql-users/repo"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		migrate_db := os.Args[1]
		if migrate_db == "migrate" {
			migrate.Migrate()
			return
		}
	}

	db_conn, _ := db.Connect()
	db_conn.LogMode(true)
	defer db_conn.Close()

	repositries, _ := repo.NewRepositories(db_conn)

	schema, err := api.GetSchema(repositries)
	if err != nil {
		log.Fatalf("failed to create new schema. Error: %v", err)
	}

	api.InitializeApi(schema, ":8080")
}
