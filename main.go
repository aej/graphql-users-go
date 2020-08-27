package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	migrate "github.com/andyjones11/graphql-users/bin"
	db "github.com/andyjones11/graphql-users/db"
	"github.com/andyjones11/graphql-users/graph"
	"github.com/andyjones11/graphql-users/graph/generated"
	"github.com/andyjones11/graphql-users/repo"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
)

const defaultPort = ":8080"

func graphqlHandler(repo *repo.Repositories) gin.HandlerFunc {
	resolvers := &graph.Resolver{Repo: repo}
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolvers}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("Graphql", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}


func main() {
	if len(os.Args) > 1 {
		migrate_db := os.Args[1]
		if migrate_db == "migrate" {
			migrate.Migrate()
			return
		}
	}

	db, _ := db.Connect()
	defer db.Close()

	repositories, _ := repo.NewRepositories(db)

	r := gin.Default()
	r.POST("/query", graphqlHandler(repositories))
	r.GET("/", playgroundHandler())
	r.Run(defaultPort)
}

