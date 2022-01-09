package websocket

import (
	"fmt"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
)

type observer struct {
	auth chan 	*model.Auth
	user chan 	*model.User
}

// Отправляет сообщение
func (o *observer) Send(ch interface{}) {

	// Получим тип из интерфейса
	switch ch.(type) {
	case *model.Auth:

		// Валидируем канал
		if o.auth == nil {
			fmt.Println("Auth sending error")
			return
		}

		// Отправляем сообщение
		o.auth <- ch.(*model.Auth)
	case *model.User:

		// Валидируем канал
		if o.user == nil {
			fmt.Println("User sending error")
			return
		}

		// Отправляем сообщение
		o.user <- ch.(*model.User)
	default:
		fmt.Println("unknown message type")
	}
}

func (o *observer) Add(ch interface{}) error {

	// Получим тип из интерфейса
	switch ch.(type) {
	case chan *model.Auth:
		o.auth = ch.(chan *model.Auth)
		return nil
	case chan *model.User:
		o.user = ch.(chan *model.User)
		return nil
	default:
		return fmt.Errorf("observer: unknown type")
	}
}

// Удаляет наблюдателя,
// если вернет true - можно удалить слушатель
func (o *observer) Delete(ch interface{}) bool {

	// Получим тип из интерфейса
	switch ch.(type) {
	case chan *model.Auth:
		o.auth = nil
	case chan *model.User:
		o.user = nil
	}

	return o.checkEmpty()
}

// Вернет истину если нет слушателей
func (o *observer) checkEmpty() bool {
	switch {
	case o.auth != nil:
		return false
	case o.user != nil:
		return false
	}
	return true
}
