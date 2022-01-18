package repo

import (
	"fmt"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
	"sync"
)

type UserRepo struct {
	users map[string]*model.User
	mu sync.Mutex
}

func (u *UserRepo) GetByUid(uid int) (*model.User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	for _, usr := range u.users {
		if usr.UID == uid {
			return usr, nil
		}
	}

	return nil, fmt.Errorf("user not found")
}

func (u *UserRepo) GetByUsername(username string) (*model.User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	usr, ok := u.users[username]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	return usr, nil
}

func newUser() *UserRepo {

	// Тестовые пользователи
	data := map[string]*model.User{}
	data["admin"] = &model.User{
		UID: 1,
		Username: "admin",
		Active: true,
		Email: "admin@example.com",
		Phone: "1234567890",
		Method: "phone",
	}
	data["user"] = &model.User{
		UID: 2,
		Username: "user",
		Active: true,
		Email: "user@example.com",
		Phone: "1234567890",
		Method: "email",
	}

	return &UserRepo{
		users: data,
	}
}
