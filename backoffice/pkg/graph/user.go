package graph

import (
	"context"
	"fmt"
	"github.com/vektah/gqlparser/v2/gqlerror"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
)

func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	return r.store.User(ctx)
}

func (r *subscriptionResolver) User(ctx context.Context) (<-chan *model.User, error) {
	user := make(chan *model.User)

	// Получим сессию из контекста
	sess, err := model.SessionFromContext(ctx)
	if err != nil {
		return nil, gqlerror.Errorf("internal error")
	}

	fmt.Printf("User. Session: %v\n", sess.Sid)

	return user, nil
}
