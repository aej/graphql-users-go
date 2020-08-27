package graph

import (
	"context"
	"errors"
	"github.com/andyjones11/graphql-users/internal/auth"
	"github.com/andyjones11/graphql-users/internal/storage"
	"github.com/labstack/echo/v4"
)

var userContextKey = "user"
var httpResponseContextKey = "HttpResponse"

func AddAnonymousUserToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userContextKey, auth.AnonymousUser{})
}

func AddUserToContext(ctx context.Context, user storage.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// Fetch the User from the context
func UserFromContext(ctx context.Context) (storage.User, error) {
	user, ok := ctx.Value(userContextKey).(storage.User)

	if ok {
		return user, nil
	}
	return storage.User{}, errors.New("user not found")
}

// Add the HttpResponse object to the context under a defined key.
// The HttpResponse object has access to methods and attributes on
// the http request/response cycle. In this case it is the echo.Context
// object - this can easily be switched out depending on the flavour
// of http library used.
func AddHttpResponseToContext(ctx context.Context, c echo.Context) context.Context {
	return context.WithValue(ctx, httpResponseContextKey, c)
}

// Fetch the HttpResponse object from the context.
func HttpResponseFromContext(ctx context.Context) echo.Context {
	httpResponse, _ := ctx.Value(httpResponseContextKey).(echo.Context)

	return httpResponse
}
