package store

import (
	"context"
	"fmt"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
)

// Вернет User согласно текущего состояния авторизации
func (s *Store) User(ctx context.Context) (*model.User, error) {

	err := s.SendAuth(ctx)
	fmt.Println("Запрашиваем Auth из метода User")
	fmt.Printf("Ошибка: %v\n", err)

	// ...
	return &model.User{
		Username: "LOLO",
	}, nil
}