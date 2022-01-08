package websocket

import (
	model "react-apollo-gqlgen-tutorial/backoffice/models"
	"sync"
)

type Websocket struct {
	clients map[string]*client

	// Защищаем мапу
	mu sync.Mutex
}

func (w *Websocket) Add(sess *model.Session) error {

	// Заблокируем мапу clients чтобы безопасно с ней работать работать
	w.mu.Lock()

	// Разблокируем мапу после выхода из функции
	defer w.mu.Unlock()

	return nil
}

func (w *Websocket) Delete() {

}

func (w *Websocket) Send() {

}

func NewWebsocket() *Websocket {
	return &Websocket{}
}
