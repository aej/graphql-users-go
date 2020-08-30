package api

import (
	"github.com/andyjones11/graphql-users/repo"
	"github.com/graph-gophers/graphql-go"
)

func GetSchema(repositories *repo.Repositories) (graphql.Schema, error) {
	s := `
		type User {
			id: ID!
			email: String!
			createdAt: String!
			updatedAt: String!
		}

		type Query {
			listUsers: [User]!
		}

		type Mutation {
			registerUser(email: String!): User
		}
	`

	schema := graphql.MustParseSchema(s)

	resolver := Resolver{Repo: repositories}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: queryFields}
	rootMutation := graphql.ObjectConfig{Name: "RootMutation", Fields: mutationFields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery), Mutation: graphql.NewObject(rootMutation)}

	schema, err := graphql.NewSchema(schemaConfig)

	return schema, err
}
