package repo

import (
	"fmt"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
	"sync"
)

type SessionRepo struct {
	sessions map[string]*model.Session
	mu sync.Mutex
}

// Получить сессию ожидающую подтверждение по email или СМС
func(s *SessionRepo) GetPending(sid string) (*model.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.sessions[sid]
	if !ok || sess.Type != model.SessPending {
		return nil, fmt.Errorf("session not found")
	}
	sess.Type = model.SessApproved
	return sess, nil
}

// Получить сессию ожидающую запрос Клиента
func(s *SessionRepo) GetApproved(sid string) (*model.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.sessions[sid]
	if !ok || sess.Type != model.SessApproved {
		return nil, fmt.Errorf("session not found")
	}
	sess.Type = model.SessInstalled
	return sess, nil
}

// Получить активную Пользовательскую сессию
func(s *SessionRepo) GetInstalled(sid string) (*model.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.sessions[sid]
	if !ok || sess.Type != model.SessInstalled {
		return nil, fmt.Errorf("session not found")
	}
	return sess, nil
}

func(s *SessionRepo) GetByCode(code string) (*model.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, sess := range s.sessions {
		if sess.SecretCode == code {
			sess.Type = model.SessApproved
			sess.SecretCode = ""
			return sess, nil
		}
	}
	return nil, fmt.Errorf("session not found")
}

func(s *SessionRepo) Add(sess *model.Session) (sid string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[sess.Sid] = sess
	return sess.Sid
}

func newSession() *SessionRepo {
	return &SessionRepo{
		sessions: make(map[string]*model.Session),
	}
}
