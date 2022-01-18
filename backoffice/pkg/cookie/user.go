package cookie

import (
	"net/http"
)

type User struct {
	response 	http.ResponseWriter
	request 	*http.Request
	name 		string
}

func NewUser(w http.ResponseWriter, r *http.Request) *User {
	return &User{
		response: w,
		request: r,
		name: "session",
	}
}

func (u *User) Delete() {
	cookie := &http.Cookie{
		MaxAge: -1,
		Name: u.name,
		Value: "",
		HttpOnly: true,
		//Secure: true,
	}
	http.SetCookie(u.response, cookie)
}

func (u *User) Add(value string) {
	cookie := &http.Cookie{
		Name: u.name,
		Value: value,
		HttpOnly: true,
		//Secure: true,
	}
	http.SetCookie(u.response, cookie)
}

func (u *User) Get() (string, bool) {
	cookie, err := u.request.Cookie(u.name)
	if err != nil {
		return "", false
	}
	return cookie.Value, true
}