package store

import (
	"context"
	"fmt"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"net/http"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
	"react-apollo-gqlgen-tutorial/backoffice/pkg/token"
)

// Метод осуществляющий активацию сессии по url
func (s *Store) AuthAcceptSessionHandleHTTP(w http.ResponseWriter, r *http.Request) {

	// получим код из url
	code := r.URL.Query().Get("code")
	if code == "" {
		return
	}

	if sess, err := s.SessionValidateURLCode(code); err == nil {

		// Создадим контекст
		ctx := context.Background()

		// Сохраним сессию в контекст
		ctx = sess.WithContext(ctx)

		// Отправим сообщение об авторизации по websocket
		_ = s.SendAuth(ctx)
	}
}

// Подтверждение авторизации из СМС сообщения
func (s *Store) AuthSMSApprove(ctx context.Context, code string) (*model.Auth, error) {
	if sess, err := s.SessionValidateSMSCode(code); err == nil {

		// Сохраним сессию в контекст
		ctx = sess.WithContext(ctx)

		// Отправим сообщение об авторизации по websocket
		_ = s.SendAuth(ctx)
	}

	return s.Auth(ctx)
}

// Авторизация по Username
func (s *Store) AuthorizeForUsername(ctx context.Context, login string) (*model.Auth, error) {

	// Получим сессию
	sess, err := model.SessionFromContext(ctx)
	if err != nil {
		// Если сессию не нашли
		return nil, gqlerror.Errorf("internal error")
	}

	// Проверим авторизацию Клиента
	if ok := sess.CheckOnline(); !ok {
		return nil, gqlerror.Errorf("internal error")
	}

	// Получим Пользователя по username
	user, err := s.repo.User.GetByUsername(login)
	if err != nil {
		return nil, gqlerror.Errorf("incorrect username")
	}

	authMethod 	:= user.GetMethod()
	uid 		:= user.GetUserID()

	// Добавим метод и UID к сессии
	sess.SetMethod(authMethod)
	sess.SetUserId(uid)
	// Сохраним сессию в контекст
	ctx = sess.WithContext(ctx)

	// Получим Auth
	auth, err := s.Auth(ctx)
	if err != nil {
		return nil, gqlerror.Errorf("incorrect username")
	}

	// Создадим код
	switch authMethod {

	// Авторизация по телефону
	case "phone":

		// - создадим СМС код
		// - получим токен авторизации Пользователя
		code, authToken, err2 := s.SessionWithSmsCode(ctx)
		if err2 != nil {
			fmt.Println(err2)
			return nil, gqlerror.Errorf("internal error")
		}

		// Добавим токен в model.Auth
		auth.AddToken(authToken)

		// Логика отправки сообщения
		fmt.Printf("Сообщение отправлено на номер: %v\n", user.Phone)
		fmt.Printf("СМС код: %v\n", code)

	// Авторизация по URL
	case "email":

		// - создадим URL код
		// - получим токен авторизации Пользователя
		code, authToken, err2 := s.SessionCreateURLCode(ctx)
		if err2 != nil {
			fmt.Println(err2)
			return nil, gqlerror.Errorf("internal error")
		}

		// Добавим токен в model.Auth
		auth.AddToken(authToken)

		fmt.Printf("Сообщение отправлено на: %v\n", user.Email)
		fmt.Printf("Url авторизации: http://localhost:2000/accept_auth?code=%v\n", code)
	default:
		fmt.Println("user method error")
		return nil, gqlerror.Errorf("internal error")
	}

	return auth, nil
}

// Вызывается в AuthMiddleware
// Обрабатывает HTTP заголовки
// Проводит авторизации клиента и пользователя
func (s *Store) HandleAuthHTTP(w http.ResponseWriter, r *http.Request) *http.Request {

	r = s.SessionHandleClient(w, r)
	r = s.SessionHandleApproved(w, r)
	r = s.SessionHandleUser(w, r)

	return r
}

// Инициирует отправку текущего состояния Auth клиенту
func (s *Store) SendAuth(ctx context.Context) error {

	// Получим текущее состояние
	auth, err := s.Auth(ctx)
	if err != nil {
		return err
	}

	if err = s.websocket.Send(ctx, auth); err != nil {
		return err
	}

	return nil
}

// Авторизовывает websocket
func (s *Store) AuthWebsocket(ctx context.Context) (<-chan *model.Auth, error) {

	// Создаем канал в который будем писать сообщения
	ch := make(chan *model.Auth)

	// Подключим канал к менеджеру websocket
	err := s.websocket.NewObserver(ctx, ch)
	if err != nil {
		return nil, gqlerror.Errorf("internal error")
	}

	// Вернем канал
	return ch, nil
}

// Возвращает текущее состояние Auth
func (s *Store) Auth(ctx context.Context) (*model.Auth, error) {
	auth := &model.Auth{}

	// Получим сессию
	sess, err := model.SessionFromContext(ctx)
	if err != nil {
		// Если сессию не нашли
		return nil, gqlerror.Errorf("internal error")
	}

	// Добавим SessionID к Auth
	// Проверим авторизацию Клиента
	if ok := sess.CheckOnline(); !ok {

		// Если нет авторизации создадим токен сессии
		if tkn, err2 := s.token.ClientSession.Generate(token.Claims{
			Session: sess,
		}); err2 != nil {
			return nil, gqlerror.Errorf("internal error")
		} else {
			auth.AddSessionId(tkn)
		}
	}

	// Проверим авторизацию Пользователя
	if ok := sess.CheckUserAuthorization(); ok {
		method, _ := sess.GetMethod()
		auth.SetMethod(method)
		auth.UserAuthorize()
	}

	// Вернем текущее состояние
	return auth, nil
}
