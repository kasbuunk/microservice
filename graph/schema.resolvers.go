package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/kasbuunk/microservice/graph/generated"
	"github.com/kasbuunk/microservice/graph/model"
	"github.com/kasbuunk/microservice/user"
)

func (r *mutationResolver) RegisterUser(ctx context.Context, input model.RegisterUserInput) (*user.User, error) {
	usr, _ := user.New(input.Email, input.Password)

	usr, err := r.UserRepository.Save(usr)
	if err != nil {
		return &usr, fmt.Errorf("saving user: %w", err)
	}
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Users(ctx context.Context) ([]*user.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) ID(ctx context.Context, obj *user.User) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Email(ctx context.Context, obj *user.User) (string, error) {
	panic(fmt.Errorf("not implemented"))
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
