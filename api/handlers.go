package api

import (
	"github.com/andyjones11/graphql-users/graph/generated"
	"github.com/gin-gonic/gin"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/andyjones11/graphql-users/graph"
	"github.com/99designs/gqlgen/graphql/handler"

)

func GraphqlHandler(resolver *graph.Resolver) gin.HandlerFunc {
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func PlaygroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
