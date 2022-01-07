package graph

import (
	"context"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
)

func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	return r.store.User(ctx)
}
