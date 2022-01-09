package store

import (
	"context"
	"fmt"
	"net/http"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
	"react-apollo-gqlgen-tutorial/backoffice/pkg/token"
)

// Валидирует сессию слушателя
func (s *Store) ValidateClientSession(ctx context.Context) (sessionID string, err error) {

	// Получим сессию из контекста
	sess, err := model.SessionFromContext(ctx)
	if err != nil {
		return "", fmt.Errorf("internal error")
	}

	if ok := sess.CheckOnline(); !ok {

		// Если клиент не авторизован: SessionID отсутствует
		// Создадим SessionID, и отправим клиенту
		sessionToken, err2 := s.token.SessionID.Generate(token.JwtClaims{
			Sess: sess,
		})

		if err2 != nil {
			fmt.Println(err2)
			return "", fmt.Errorf("internal error")
		}

		return sessionToken, nil
	}

	return "", nil
}

// Валидирует токен слушателя
func (s *Store) ValidateSessionToken(sid string) (*model.Session, error) {

	// Валидируем токен
	// Считаем токен не валидным если нет claims
	if claims, _ := s.token.SessionID.Validate(token.JwtValidateOptions{
		Token: sid,
	}); claims != nil {
		sess := claims.Sess

		// Сессию получили из заголовка: клиент онлайн
		sess.SetOnline()

		// Сохраним сессию в контекст
		return sess, nil
	}

	return nil, fmt.Errorf("invalid session token")
}

// Обрабатывает сессию клиента
func (s *Store) SessionHandleClient(w http.ResponseWriter, r *http.Request) *http.Request {

	// Получим контекст
	ctx := r.Context()

	// Сюда запишем сессию, если сработает кейс
	var sess *model.Session
	var ClientID string

	// Проверим наличие cookie c ClientID
	cookie, err := r.Cookie("_cid")
	if err == nil {
		ClientID = cookie.Value

		// У клиента есть ClientID
		// 1. Проверим наличие заголовка Session-ID
		// 2. Получаем токен и валидируем его
		// 2.1. Токен валидный: сохраним сессию из токена
		// 2.2. Токен протух: сохраним сессию из токена
		// 2.3. Токен Invalid: создадим новую сессии

		// Ищем заголовок Session-ID
		if t := r.Header.Get("Session-ID"); t != "" {

			// Нашли сессию
			if ss, err2 := s.ValidateSessionToken(t); err2 == nil {
				sess = ss
			}
		}

		// Этот метод теперь удален
		//sess = model.NewSessionWithSid(cookie.Value)
	}

	// Если сессии нет: создаем сессию
	if sess == nil {
		sess = model.NewSession()

		if ClientID != "" {
			sess.AddClientID(ClientID)
		}
	}

	// Если есть ошибка при чтении cookie
	if err != nil {

		// Получим ID клиента
		cid, err2 := sess.GetSid()
		if err2 != nil {
			fmt.Printf(err.Error())
			return r
		}

		// Создадим cookie
		cookie = &http.Cookie{
			Name: "_cid",
			Value: cid,
			HttpOnly: true,
			//Secure: true,
		}

		// Установим cookie
		http.SetCookie(w, cookie)
	}

	// Сохраним сессию в контекст и вернем *http.Request
	return r.WithContext(sess.WithContext(ctx))
}
