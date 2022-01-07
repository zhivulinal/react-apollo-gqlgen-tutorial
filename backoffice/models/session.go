package model

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

// Создает новую сессию
func NewSession() *Session {
	return &Session{

		// Идентификатором будет UUID
		// go get github.com/google/uuid
		// или go mod vendor при указании импорта
		Sid: uuid.New().String(),
	}
}

// Создаст сессию с существующим идентификатором
func NewSessionWithSid(sid string) *Session {
	return &Session{
		Sid: sid,
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
