package websocket

import (
	"fmt"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
)

type observer struct {
	auth chan 	*model.Auth
	user chan 	*model.User
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
