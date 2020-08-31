package api

import (
	"github.com/andyjones11/graphql-users/repo"
	"github.com/andyjones11/graphql-users/services/auth"
	userservice "github.com/andyjones11/graphql-users/services/user"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type Resolver struct {
	Repo *repo.Repositories
}

func (r *Resolver) ListUsersResolver() []userservice.User {
	users, _ := r.Repo.User.ListAllUsers()
	return users
}

func (r *Resolver) RegisterUserResolver(email string) userservice.User {

	user, _ := r.Repo.User.CreateUser(email)

	//if err != nil {
	//	if err == repo.UserEmailExists {
	//		return nil, UserEmailExists
	//	}
	//}

	_, auth_err := auth.AuthenticateUser(r.Repo, &user)

	if auth_err != nil {
		panic("Unable to authenticate user")
	}

	return user
}

func InitializeApi(schema graphql.Schema, defaultPort string) {
	h := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: false,
	})

	r := gin.Default()
	r.POST("/graphql", GraphqlHandler(h))
	r.GET("/graphql", GraphqlHandler(h))
	r.Run(defaultPort)
}
