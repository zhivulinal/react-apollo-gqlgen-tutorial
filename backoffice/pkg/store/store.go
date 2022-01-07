package store

type Store struct {
	token TokenOptions
}

func NewStore(opt Options) *Store {
	return &Store{
		token: opt.Token,
	}
}

type Options struct {
	Token TokenOptions
}

type TokenOptions struct {}