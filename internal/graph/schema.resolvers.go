package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/andyjones11/graphql-users/internal/auth"
	"github.com/andyjones11/graphql-users/internal/graph/generated"
	"github.com/andyjones11/graphql-users/internal/graph/model"
	userservice "github.com/andyjones11/graphql-users/internal/user"
	"github.com/andyjones11/graphql-users/internal/validators"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func (r *mutationResolver) RegisterUser(ctx context.Context, input model.RegisterUserInput) (*model.User, error) {
	if _, err := UserFromContext(ctx); err == nil {
		return &model.User{}, AuthenticatedDisallowed
	}

	validations := validation.ValidateStruct(
		&input,
		validation.Field(&input.Email, is.Email),
		validation.Field(&input.Password, validation.By(validators.MinLength(8))),
	)
	if validations != nil {
		errs, _ := validations.(validation.Errors)

		for key, val := range errs {
			graphql.AddErrorf(ctx, "%s: %s", key, val)
		}

		return &model.User{}, nil
	}

	dbUser, err := userservice.CreateUser(r.Repos.User, input.Email, input.Password, input.FullName)

	return DbUserToGqlUser(&dbUser), err
}

func (r *mutationResolver) LoginUser(ctx context.Context, input model.LoginInput) (*model.User, error) {
	if _, err := UserFromContext(ctx); err == nil {
		return &model.User{}, AuthenticatedDisallowed
	}

	user, err := auth.ValidateUserCredentials(r.Repos, input.Email, input.Password)

	if err != nil {
		if errors.Is(err, auth.InvalidCredentials) {
			return &model.User{}, InvalidEmailPasswordCombination
		}
		if errors.Is(err, auth.UserUnconfirmed) {
			return &model.User{}, UnconfirmedEmail
		}
		return &model.User{}, ServerError
	}

	httpResponse := HttpResponseFromContext(ctx)

	if _, err := auth.AuthenticateUser(r.Repos, user, httpResponse); err != nil {
		return DbUserToGqlUser(&user), ServerError
	}

	return DbUserToGqlUser(&user), nil
}

func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	user, err := UserFromContext(ctx)

	if err != nil {
		return false, AuthenticationRequired
	}

	httpResponse := HttpResponseFromContext(ctx)

	auth.DeauthenticateUser(r.Repos, user, httpResponse)

	return true, nil
}

func (r *mutationResolver) ConfirmEmail(ctx context.Context, input model.ConfirmEmailInput) (bool, error) {
	if _, err := UserFromContext(ctx); err == nil {
		return false, AuthenticatedDisallowed
	}

	if err := userservice.ConfirmUserEmail(r.Repos.User, input.Token); err != nil {
		return false, nil
	}
	return true, nil
}

func (r *mutationResolver) RequestResetPassword(ctx context.Context, input model.RequestResetPasswordInput) (bool, error) {
	if _, err := UserFromContext(ctx); err == nil {
		return false, AuthenticatedDisallowed
	}

	err := userservice.RequestResetPassword(r.Repos.User, input.Email)

	// Return true even in the error-case
	if err != nil {
		return true, nil
	}
	return true, nil
}

func (r *mutationResolver) ValidateResetPasswordToken(ctx context.Context, input model.ValidateResetPasswordTokenInput) (*model.User, error) {
	if _, err := UserFromContext(ctx); err == nil {
		return &model.User{}, AuthenticatedDisallowed
	}

	user, err := userservice.ValidatePasswordResetToken(r.Repos.User, input.Token)
	if err != nil {
		return DbUserToGqlUser(&user), errors.New("invalid password reset token")
	}
	return DbUserToGqlUser(&user), nil
}

func (r *mutationResolver) ResetPassword(ctx context.Context, input model.ResetPasswordInput) (bool, error) {
	if _, err := UserFromContext(ctx); err == nil {
		return false, AuthenticatedDisallowed
	}

	result, err := userservice.ResetPassword(r.Repos.User, input.Password, input.Token)

	if err != nil {
		return false, errors.New("invalid token")
	}

	return result, nil
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	user, err := UserFromContext(ctx)

	if err != nil {
		return &model.User{}, AuthenticationRequired
	}

	return DbUserToGqlUser(&user), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
