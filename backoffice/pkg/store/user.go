package store

import (
	"context"
	"github.com/vektah/gqlparser/v2/gqlerror"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
)

// Вернет User согласно текущего состояния авторизации
func (s *Store) User(ctx context.Context) (*model.User, error) {

	// Получим сессию
	sess, err := model.SessionFromContext(ctx)
	if err != nil {

		// Если сессию не нашли
		return nil, gqlerror.Errorf("internal error")
	}

	uid, err := sess.GetUserId()
	if err != nil {
		return &model.User{}, nil
	}

	user, err := s.repo.User.GetByUid(uid)
	if err != nil {
		return &model.User{}, nil
	}

	return user, nil
}