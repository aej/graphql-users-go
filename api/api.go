package api

import (
	"github.com/andyjones11/graphql-users/repo"
	"github.com/andyjones11/graphql-users/services/auth"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type Resolver struct {
	Repo *repo.Repositories
}

func (r *Resolver) ListUsersResolver(p graphql.ResolveParams) (interface{}, error) {
	users, _ := r.Repo.User.ListAllUsers()
	return users, nil
}

func (r *Resolver) RegisterUserResolver(p graphql.ResolveParams) (interface{}, error) {

	user, err := r.Repo.User.CreateUser(p.Args["email"].(string))

	if err != nil {
		if err == repo.UserEmailExists {
			return nil, UserEmailExists
		}
	}

	_, auth_err := auth.AuthenticateUser(r.Repo, &user)

	if auth_err != nil {
		panic("Unable to authenticate user")
	}

	return user, err
}

//func AuthMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		c.Next()
//		auth_cookie := c.Request.Context().Value("auth-cookie")
//		if auth_cookie != nil {
//			c.SetCookie("_id", auth_cookie, 3600, "/", "localhost", false, true)
//		}
//	}
//}

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
