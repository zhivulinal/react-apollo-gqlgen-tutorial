package store

import (
	"fmt"
	"net/http"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
)

// Обрабатывает сессию клиента
func (s *Store) SessionHandleClient(w http.ResponseWriter, r *http.Request) *http.Request {

	// Получим контекст
	ctx := r.Context()

	// Создадим сессию
	var sess *model.Session

	// Проверим наличие токена c ClientID
	cookie, err := r.Cookie("_sid")
	if err != nil {

		// Нет ClientID, создадим сессию
		sess = model.NewSession()

	} else {

		// Тут должна быть логика валидации
		// Но нам сейчас удобно видеть действительную запись
		sess = model.NewSessionWithSid(cookie.Value)

		// Клиент имеет ID, соединение по websocket возможно
		sess.SetOnline()
	}

	// Если есть ошибка – устанавливаем новые cookie
	if err != nil {

		// Получим ID клиента
		sid, err2 := sess.GetSid()
		if err2 != nil {
			fmt.Printf(err.Error())
			return r
		}

		// Создадим cookie
		cookie = &http.Cookie{
			Name: "_sid",
			// Сid следует завернуть в токен, например JWT
			Value: sid,
			HttpOnly: true,
			//Secure: true,
		}

		// Установим cookie
		http.SetCookie(w, cookie)
	}

	// Сохраним сессию в контекст и вернем *http.Request
	return r.WithContext(sess.WithContext(ctx))
}
