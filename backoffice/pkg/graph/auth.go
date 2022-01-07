package graph

import (
	"context"
	"fmt"
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

func (r *subscriptionResolver) Auth(ctx context.Context) (<-chan *model.Auth, error) {
	panic(fmt.Errorf("not implemented"))
}
