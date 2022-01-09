package websocket

import (
	"context"
	"fmt"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
	"sync"
)

type Websocket struct {
	clients map[string]*client

	// Защищаем мапу
	mu sync.Mutex
}

func (w *Websocket) Send(ctx context.Context, ch interface{}) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Получим сессию из контекста
	sess, err := model.SessionFromContext(ctx)
	if err != nil {
		return err
	}

	// Получим ClientID
	cid := sess.ClientID

	// Найдем клиента
	cli, ok := w.clients[cid]
	if !ok {
		return fmt.Errorf("client not found")
	}

	// Отправляем сообщение
	cli.Send(ch)

	return nil
}

// Создает Клиента
func (w *Websocket) NewObserver(ctx context.Context, ch interface{}) error {

	// Заблокируем мапу clients
	// чтобы безопасно с ней работать работать
	w.mu.Lock()

	// Разблокируем мапу после выхода из функции
	defer w.mu.Unlock()

	// Получим сессию из контекста
	sess, err := model.SessionFromContext(ctx)
	if err != nil {
		return err
	}

	// Получим ClientID и SessionID
	cid := sess.ClientID
	sid := sess.Sid

	// Найдем, или создадим клиента
	cli, ok := w.clients[cid]
	if !ok {

		// Клиент не найден, создадим
		cli = newClient()

		// Добавим в мапу
		w.clients[cid] = cli
	}

	// Добавим слушателя клиента
	err = cli.Add(sid, ch)
	if err != nil {
		return err
	}

	// Клиент отписывается – удаляем слушатель
	go func() {
		<- ctx.Done()

		cli.Delete(sid, ch)
	}()

	return nil
}

func New() *Websocket {
	return &Websocket{
		clients: make(map[string]*client),
	}
}
