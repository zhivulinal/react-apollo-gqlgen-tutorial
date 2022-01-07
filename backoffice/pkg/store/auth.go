package store

import (
	"context"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
)

// Возвращает состояние Auth исходя из текущего контекста
func (s *Store) Auth(ctx context.Context) (auth *model.Auth, err error) {
	// ...
	return
}

// Авторизовывает websocket, обрабатывает подключение и создает канал
func (s *Store) AuthCreateWebsocket(ctx context.Context) (out <-chan *model.Auth, err error) {
	// ...
	return
}

// Авторизация по Username
func (s *Store) AuthorizeForUsername(ctx context.Context, login string) (auth *model.Auth, err error) {
	// ...
	return
}

// Подтверждение авторизации из СМС сообщения
func (s *Store) AuthSMSApprove(ctx context.Context, code string) (auth *model.Auth, err error) {
	// ...
	return
}