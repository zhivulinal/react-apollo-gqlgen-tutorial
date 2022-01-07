package graph

import (
	"context"
	"fmt"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
)

func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}
