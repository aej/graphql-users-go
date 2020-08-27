package http

import (
	"github.com/andyjones11/graphql-users/internal/auth"
	"github.com/andyjones11/graphql-users/internal/graph"
	"github.com/andyjones11/graphql-users/internal/storage"
	"github.com/labstack/echo/v4"
)

// Puts the echo.Context object into the Go context under the key "HttpResponse"
// This enables upstream functions to access the request/response objects
// from within the resolver functions
func GraphqlContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := graph.AddHttpResponseToContext(c.Request().Context(), c)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}

// Checks the request for a user session. If a session if found
// then set the user onto the request context for use upstream
func AuthenticationMiddleware(repo *storage.Repositories) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, err := auth.CheckAuth(repo, c)
			if err != nil {
				ctx := graph.AddAnonymousUserToContext(c.Request().Context())
				c.SetRequest(c.Request().WithContext(ctx))
				return next(c)
			}
			ctx := graph.AddUserToContext(c.Request().Context(), user)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)

		}
	}
}
