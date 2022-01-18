package repo

type Repo struct {
	User 	*UserRepo
	Session *SessionRepo
}

func New(opt Options) *Repo {
	return &Repo{
		User: 		newUser(),
		Session: 	newSession(),
	}
}

type Options struct {}
