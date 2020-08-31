package api

import (
	"context"
	"github.com/andyjones11/graphql-users/graph"
	"github.com/andyjones11/graphql-users/repo"
	"github.com/andyjones11/graphql-users/services/auth"
	"github.com/gin-gonic/gin"
)

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "HttpResponse", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

type AnonymousUser struct{}

// Checks the request for a user session. If a session if found
// then set the user onto the request context for use upstream
func AuthenticationMiddleware(repo *repo.Repositories) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := auth.CheckAuth(repo, c)
		if err != nil {
			// user is not authenticated so just call next
			ctx := context.WithValue(c.Request.Context(), "user", AnonymousUser{})
			c.Request = c.Request.WithContext(ctx)
			c.Next()
		}
		ctx := context.WithValue(c.Request.Context(), "user", user)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func InitializeApi(resolver *graph.Resolver) {
	r := gin.Default()
	r.Use(GinContextToContextMiddleware())
	r.Use(AuthenticationMiddleware(resolver.Repos))
	r.POST("/query", GraphqlHandler(resolver))
	r.GET("/graphql", PlaygroundHandler())
	r.Run("localhost:8080")
}
