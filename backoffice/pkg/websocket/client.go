package websocket

import (
	"fmt"
	"sync"
)

type client struct {
	observers 	map[string]*observer
	mu 			sync.Mutex
}

func (c *client) Send(ch interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Необходимо отправить сообщение всем слушателям
	// Пройдемся в цикле и запустим отправку
	//
	// Лучше всего это сделать в отдельной горутине

	// Создадим WaitGroup
	// Про применение описано в этой статье:
	// https://habr.com/ru/company/otus/blog/557312/
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for _, obs := range c.observers {
			obs.Send(ch)
		}
		wg.Done()
	}()
	wg.Wait()
}

// Добавляет слушателя Клиента
func (c *client) Add(sid string, ch interface{}) error {

	// Заблокируем мапу слушателей
	// чтобы безопасно с ней работать работать
	c.mu.Lock()

	// Разблокируем мапу после выхода из функции
	defer c.mu.Unlock()

	// Поищем слушателя
	obs, ok := c.observers[sid]
	if !ok {

		// Слушатель не найден, создадим
		obs = &observer{}

		// Добавим в мапу
		c.observers[sid] = obs
	}

	err := obs.Add(ch)
	if err != nil {
		return err
	}

	return nil
}

// Удаляет слушателя
// Возвращает признак наличия других слушателей
func (c *client) Delete(sid string, ch interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	obs, ok := c.observers[sid]
	if !ok {
		// Обсервер не найден?
		fmt.Println("log panic")
	}

	// Удаляем канал
	if ok = obs.Delete(ch); ok {

		// Если вернулся признак пустоты
		// Удалим слушатель
		delete(c.observers, sid)
	}

	// Посчитаем количество слушателей
	// и вернем результат
	return len(c.observers) == 0
}

func newClient() *client {
	return &client{
		observers: make(map[string]*observer),
	}
}
