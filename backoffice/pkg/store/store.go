package store

import (
	"react-apollo-gqlgen-tutorial/backoffice/pkg/repo"
	"react-apollo-gqlgen-tutorial/backoffice/pkg/token"
	"react-apollo-gqlgen-tutorial/backoffice/pkg/websocket"
)

type Store struct {
	token 		TokenOptions
	websocket 	*websocket.Websocket
	repo 		*repo.Repo
}

func NewStore(opt Options) *Store {
	return &Store{
		token:     opt.Token,
		websocket: websocket.New(),
		repo:      repo.New(repo.Options{}),
	}
}

type Options struct {
	Token TokenOptions
}

type TokenOptions struct {
	ClientSession 	token.Token
	AuthToken 		token.Token
	UrlToken 		token.Token
	UserToken 		token.Token
	AccessToken 	token.Token
}