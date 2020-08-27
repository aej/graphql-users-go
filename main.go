package main

import (
	migrate "github.com/andyjones11/graphql-users/bin"
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
}

