package store

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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

// Авторизовывает websocket
// Создает сессию
// Обрабатывает подключение и создает канал
//
// Каждый клиент вызывавший данный
// метод – является уникальным
func (r *Store) AuthWebsocket(ctx context.Context) (<-chan *model.Auth, error) {

	// Получим сессию из контекста
	sess, err := model.SessionFromContext(ctx)
	if err != nil {

		// Если произошла ошибка то не стоит здесь
		// ее отправлять дальше.
		//
		// Ее нужно логировать и вернуть на фронт
		// что-то более обобщенное
		fmt.Printf("Auth subscriptionResolver. %v", err)

		return nil, gqlerror.Errorf("internal error")
	}

	// Проверим инициатора запроса.
	// Если запрос поступил по вебсокет и от Клиента
	// ранее не имеющего ClientID – не обрабатываем его
	if ok := sess.CheckOnline(); !ok {

		// Если клиент не имеет авторизации
		return nil, gqlerror.Errorf("unauthorized")
	}

	// Подключившийся клиент – уникален
	// Создадим websocket ID
	wsid := uuid.New().String()

	// Создаем канал в который будем писать сообщения
	in := make(chan *model.Auth)

	// Выведем в терминал сообщение при подключении
	fmt.Printf("WS connect. ID: %v\n", wsid)

	// Обработаем остановку соединения
	go func() {

		// Чтобы узнать об отключении websocket
		// достаточно слушать сигнал из контекста
		<- ctx.Done()
		fmt.Printf("WS disconnect. ID: %v\n", wsid)
	}()

	// Тестовая публикация сообщения
	go func() {

		// Сразу опубликуем сообщение
		in <- &model.Auth{
			ClientID: time.Now().String(),
		}

		// Небольшая задержка и отправим следующее
		time.Sleep(time.Second * 2)
		in <- &model.Auth{
			ClientID: time.Now().String(),
		}
	}()

	// Вернем канал
	return in, nil
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