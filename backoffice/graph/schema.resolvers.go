package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"react-apollo-gqlgen-tutorial/backoffice/graph/generated"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
)

func (r *mutationResolver) Authorization(ctx context.Context, login string) (*model.Auth, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SmsCode(ctx context.Context, code string) (*model.Auth, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Auth(ctx context.Context) (*model.Auth, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *subscriptionResolver) Auth(ctx context.Context) (<-chan *model.Auth, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
