package http

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/andyjones11/graphql-users/internal/graph"
	"github.com/andyjones11/graphql-users/internal/graph/generated"
	"github.com/labstack/echo/v4"
	"net/http"
)

func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

func GraphqlHandler(resolver *graph.Resolver) echo.HandlerFunc {
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	return func(c echo.Context) error {
		request := c.Request()
		response := c.Response()
		h.ServeHTTP(response, request)
		return nil
	}
}

func PlaygroundHandler() echo.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c echo.Context) error {
		request := c.Request()
		response := c.Response()
		h.ServeHTTP(response, request)
		return nil
	}
}
