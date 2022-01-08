package websocket

import "sync"

type client struct {
	observers map[string]*observer
	mu sync.Mutex
}
