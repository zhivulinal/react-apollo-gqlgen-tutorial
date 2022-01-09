package model

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

// Привязывает ClientID к сессии
func (s *Session) AddClientID(cid string) {
	s.ClientID = cid
}

// Получает сессию из контекста
func SessionFromContext(ctx context.Context) (*Session, error) {
	if meta := ctx.Value(sessionCtxKey{"session"}); meta != nil {
		return meta.(*Session), nil
	}
	return nil, fmt.Errorf("meta: not found")
}

// Подтверждает активность клиента
func (s *Session) CheckOnline() bool {
	return s.Online
}

// Создает новую сессию
func NewSession() *Session {
	return &Session{

		// Идентификатором будет UUID
		// go get github.com/google/uuid
		// или go mod vendor при указании импорта
		Sid: 		uuid.New().String(),
		ClientID: 	uuid.New().String(),
	}
}

// Получить идентификатор сессии
func (s *Session) GetSid() (sid string, err error) {
	if s.Sid == "" {
		return "", fmt.Errorf("session: not found")
	}
	return s.Sid, nil
}

// Подтверждает активность клиента
func (s *Session) SetOnline() {
	s.Online = true
	return
}

// Сохраняет сессию в контекст
func (s *Session) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, sessionCtxKey{"session"}, s)
}

// Ключ контекста
type sessionCtxKey struct {
	name string
}
