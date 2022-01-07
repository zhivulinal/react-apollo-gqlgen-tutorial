package store

import (
	"context"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
)

// Вернет User согласно текущего состояния авторизации
func (s *Store) User(ctx context.Context) (user *model.User, err error) {
	// ...
	return
}