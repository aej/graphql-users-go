package http

import (
	"github.com/andyjones11/graphql-users/internal/graph"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitializeApi(resolver *graph.Resolver) {
	e := echo.New()

	// Echo middleware
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())

	// Custom middleware
	e.Use(GraphqlContextMiddleware)
	e.Use(AuthenticationMiddleware(resolver.Repos))

	e.GET("/health", HealthCheck)
	e.POST("/query", GraphqlHandler(resolver))
	e.GET("/graphql", PlaygroundHandler())

	e.Logger.Fatal(e.Start(":8000"))
}
