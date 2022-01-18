package model

func (a *Auth) UserAuthorize() {
	a.Authorized = true
}

func (a *Auth) AddSessionId(sid string) {
	a.SessionID = sid
}

func (a *Auth) AddToken(token string) {
	a.Token = token
}

func (a *Auth) AddMethod(method string) {
	a.Method = method
}

func (a *Auth) SetMethod(method string) {
	a.Method = method
}
