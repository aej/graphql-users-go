package graph

import "github.com/andyjones11/graphql-users/repo"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	Repo *repo.Repositories
}
