package main

import (
	"github.com/andyjones11/graphql-users/internal/config"
	"github.com/andyjones11/graphql-users/internal/graph"
	"github.com/andyjones11/graphql-users/internal/http"
	"github.com/andyjones11/graphql-users/internal/storage"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	configuration := config.GetConfig()

	dbConn, err := storage.Connect(configuration.DatabaseUrl)
	if err != nil {
		panic("Unable to connect to database")
	}
	dbConn.LogMode(true)
	defer dbConn.Close()

	repositories, _ := storage.NewRepositories(dbConn)
	resolver := &graph.Resolver{Repos: repositories}

	http.InitializeApi(resolver)
}
