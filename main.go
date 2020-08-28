package main

import (
	migrate "github.com/andyjones11/graphql-users/bin"
	"github.com/andyjones11/graphql-users/db"
	"github.com/andyjones11/graphql-users/repo"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
)

const defaultPort = ":8080"

func graphqlHandler(h *handler.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

type Resolver struct {
	Repo *repo.Repositories
}

func (r *Resolver) ListUsersResolver(p graphql.ResolveParams) (interface{}, error) {
	users, _ := r.Repo.User.ListAllUsers()
	return users, nil
}

func (r *Resolver) RegisterUserResolver(p graphql.ResolveParams) (interface{}, error) {
	user, err := r.Repo.User.CreateUser(p.Args["email"].(string))
	return user, err
}

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
	resolver := Resolver{Repo: repositries}

	userType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "User",
		Description: "A User",
		Fields: graphql.Fields{
			"id":        &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"email":     &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"createdAt": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"updatedAt": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		},
	})

	queryFields := graphql.Fields{
		"listUsers": &graphql.Field{
			Type:        graphql.NewList(userType),
			Description: "List of all users",
			Resolve:     resolver.ListUsersResolver,
		},
	}

	mutationFields := graphql.Fields{
		"registerUser": &graphql.Field{
			Type:        userType,
			Description: "Register a new user",
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: resolver.RegisterUserResolver,
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: queryFields}
	rootMutation := graphql.ObjectConfig{Name: "RootMutation", Fields: mutationFields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery), Mutation: graphql.NewObject(rootMutation)}
	schema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		log.Fatalf("failed to create new schema. Error: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: false,
	})

	r := gin.Default()
	r.POST("/graphql", graphqlHandler(h))
	r.GET("/graphql", graphqlHandler(h))
	r.Run(defaultPort)
}
