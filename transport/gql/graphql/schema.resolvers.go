package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/kasbuunk/microservice/auth/models"
	"github.com/kasbuunk/microservice/transport/gql/graphql/generated"
	"github.com/kasbuunk/microservice/transport/gql/graphql/model"
)

// RegisterUser is the resolver for the registerUser field.
func (r *mutationResolver) RegisterUser(ctx context.Context, input model.RegisterUserInput) (*models.User, error) {
	registeredUser, err := r.Auth.Register(models.EmailAddress(input.Email), models.Password(input.Password))
	if err != nil {
		return nil, fmt.Errorf("registering user: %w", err)
	}
	return &registeredUser, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// ID is the resolver for the id field.
func (r *userResolver) ID(ctx context.Context, obj *models.User) (string, error) {
	return "1023012", nil
	//panic(fmt.Errorf("not implemented"))
}

// Email is the resolver for the email field.
func (r *userResolver) Email(ctx context.Context, obj *models.User) (string, error) {
	return "gjlksdjfl", nil
	//panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
