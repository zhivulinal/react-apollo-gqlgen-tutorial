package model

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

// Типы сессии
// 1 - SessPending: ожидает подтверждения – переход по ссылке, ввод кода из СМС
// 2 - SessApproved: подтвержденная сессия, но клиент пока не обнаружен
// 3 - SessInstalled: активная сессия, валидация токеном
type SessType int
const (
	SessPending SessType = iota + 1
	SessApproved
	SessInstalled
)

type Session struct {

	// SecretCode - не должен передаваться Клиенту
	SecretCode 	string
	Type 		SessType

	Sid 		string 	`json:"sid"`
	ClientID 	string 	`json:"clientID"`
	Online 		bool 	`json:"online"`
	AuthToken 	string 	`json:"authToken"`
	UID 		int 	`json:"uid"`
	Method 		string 	`json:"method"`
}


// Создает новую сессию
func NewSession(cid string) *Session {
	return &Session{

		// Идентификатором будет UUID
		// go get github.com/google/uuid
		// или go mod vendor при указании импорта
		Sid: 		uuid.New().String(),
		ClientID: 	cid,
		Type: 		SessPending,
	}
}

// SetSecretCode - добавляет секретный код
// используется для подтверждения сессии по СМС
func (s *Session) SetSecretCode(code string) {
	s.SecretCode = code
}

// GetSecretCode - сверяет секретный код
func (s *Session) CheckSecretCode(code string) bool {
	if s.SecretCode == code {
		return true
	}
	return false
}

// Получить тип сессии
// Типы:
// 1 - SessPending: ожидает подтверждения
// 2 - SessApproved: подтвержденная сессия
// 3 - SessInstalled: активная сессия
func (s *Session) GetType() SessType {
	return s.Type
}

// SessionApproved - подтвердить сессию
// Активация по СМС или email
func (s *Session) SessionApproved() {
	s.Type = SessApproved
}

// SessionInstalled - активировать сессию
// для использования Пользователем
func (s *Session) SessionInstalled() {
	s.Type = SessInstalled
}

func (s *Session) GetSessionID() string {
	return s.Sid
}

func (s *Session) GetClientID() string {
	return s.ClientID
}

// Проверяем активность сессии
func (s *Session) CheckOnline() bool {
	return s.Online
}

func (s *Session) SetOnline() {
	s.Online = true
}

// Необходим при авторизации пользователя
func (s *Session) SetAuthToken(token string) {
	s.AuthToken = token
}

func (s *Session) CheckAuthToken(token string) bool {
	if s.AuthToken == token {
		return true
	}
	return false
}

func (s *Session) SetUserId(uid int) {
	s.UID = uid
}

func (s *Session) GetUserId() (uid int, err error) {
	if s.UID == 0 {
		return 0, fmt.Errorf("user not found")
	}
	return s.UID, nil
}

//Method
func (s *Session) SetMethod(method string) {
	s.Method = method
}

func (s *Session) GetMethod() (method string, err error) {
	if s.Method == "" {
		return "", fmt.Errorf("method not set")
	}
	return s.Method, nil
}

// Состояние авторизации пользователя
func (s *Session) CheckUserAuthorization() bool {
	if s.UID == 0 || s.Method == "" {
		return false
	}
	return true
}

// Получаем сессию из контекста
func SessionFromContext(ctx context.Context) (*Session, error) {
	if meta := ctx.Value(sessionCtxKey{"session"}); meta != nil {
		return meta.(*Session), nil
	}
	return nil, fmt.Errorf("session not found")
}

// Сохраняет сессию в контекст
func (s *Session) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, sessionCtxKey{"session"}, s)
}

// Ключ контекста
type sessionCtxKey struct {
	name string
}
