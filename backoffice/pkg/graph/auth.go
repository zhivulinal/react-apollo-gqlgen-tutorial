package graph

import (
	"context"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
)

func (r *queryResolver) Auth(ctx context.Context) (*model.Auth, error) {
	return r.store.Auth(ctx)
}

func (r *subscriptionResolver) Auth(ctx context.Context) (<-chan *model.Auth, error) {
	return r.store.AuthCreateWebsocket(ctx)
}

func (r *mutationResolver) Authorization(ctx context.Context, login string) (*model.Auth, error) {
	return r.store.AuthorizeForUsername(ctx, login)
}

func (r *mutationResolver) SmsCode(ctx context.Context, code string) (*model.Auth, error) {
	return r.store.AuthSMSApprove(ctx, code)
}
