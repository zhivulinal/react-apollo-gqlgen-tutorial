package model

func (u *User) GetMethod() string {
	return u.Method
}

func (u *User) GetUserID() int {
	return u.UID
}
