package store

import (
	"context"
	"net/http"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
)

// Вызывается в AuthMiddleware
// Обрабатывает HTTP заголовки
// Проводит авторизации клиента и пользователя
func (s *Store) HandleAuthHTTP(w http.ResponseWriter, r *http.Request) *http.Request {

	// Обработаем сессию клиента
	r = s.SessionHandleClient(w, r)

	return r
}

// Авторизовывает websocket
// Создает сессию
// Обрабатывает подключение и создает канал
func (s *Store) AuthWebsocket(ctx context.Context) (out <-chan *model.Auth, err error) {

	return
}

// Возвращает состояние Auth исходя из текущего контекста
func (s *Store) Auth(ctx context.Context) (auth *model.Auth, err error) {

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