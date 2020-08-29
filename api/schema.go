package api

import (
	"github.com/andyjones11/graphql-users/repo"
	"github.com/graphql-go/graphql"
)

func GetSchema(repositories *repo.Repositories) (graphql.Schema, error) {
	resolver := Resolver{Repo: repositories}

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

	return schema, err
}
