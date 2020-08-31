package main

import (
	"github.com/andyjones11/graphql-users/api"
	migrate "github.com/andyjones11/graphql-users/bin"
	"github.com/andyjones11/graphql-users/db"
	"github.com/andyjones11/graphql-users/graph"
	"github.com/andyjones11/graphql-users/repo"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		migrate_db := os.Args[1]
		if migrate_db == "migrate" {
			migrate.Up()
			return
		}
	}

	db_conn, err := db.Connect()
	if err != nil {
		panic("Unable to connect to database")
	}
	db_conn.LogMode(true)
	defer db_conn.Close()

	repositories, _ := repo.NewRepositories(db_conn)
	resolver := &graph.Resolver{Repos: repositories}

	api.InitializeApi(resolver)
}
