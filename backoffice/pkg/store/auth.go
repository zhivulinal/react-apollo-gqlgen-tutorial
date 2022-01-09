package store

import (
	"context"
	"fmt"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"net/http"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
	"time"
)

// Вызывается в AuthMiddleware
// Обрабатывает HTTP заголовки
// Проводит авторизации клиента и пользователя
func (s *Store) HandleAuthHTTP(w http.ResponseWriter, r *http.Request) *http.Request {

	// Обработаем сессию клиента
	r = s.SessionHandleClient(w, r)

	return r
}

// Инициирует отправку текущего состояния Auth клиенту
func (s *Store) SendAuth(ctx context.Context) error {

	// Получим текущее состояние
	auth, err := s.Auth(ctx)
	if err != nil {
		return err
	}

	// Todo: удалить!!!
	// Чтобы увидеть результат изменений
	// Нужно что нибудь рандомное
	auth.Method = time.Now().String()

	if err = s.websocket.Send(ctx, auth); err != nil {
		return err
	}

	return nil
}

// Авторизовывает websocket
func (s *Store) AuthWebsocket(ctx context.Context) (<-chan *model.Auth, error) {

	// Получим текущее состояние авторизации
	auth, err := s.Auth(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, gqlerror.Errorf("internal error")
	}

	// Создаем канал в который будем писать сообщения
	ch := make(chan *model.Auth)

	// Подключим канал к менеджеру websocket
	err = s.websocket.NewObserver(ctx, ch)
	if err != nil {
		fmt.Println(err)
		return nil, gqlerror.Errorf("internal error")
	}

	// Нужно вернуть текущее состояние
	go func() {
		ch <- auth
	}()

	// Вернем канал
	return ch, nil
}

// Возвращает состояние Auth исходя из текущего контекста
func (s *Store) Auth(ctx context.Context) (*model.Auth, error) {

	// создадим модель
	auth := &model.Auth{}

	// Проверим сессию
	sid, err := s.ValidateClientSession(ctx)
	if err != nil {
		return nil, gqlerror.Errorf("internal error")
	}

	// Если есть sid – добавим его к Auth
	if sid != "" {
		auth.AddSessionId(sid)
	}

	// Отправим текущее состояние
	return auth, nil
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