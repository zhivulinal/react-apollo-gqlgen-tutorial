package store

import (
	"react-apollo-gqlgen-tutorial/backoffice/pkg/token"
	"react-apollo-gqlgen-tutorial/backoffice/pkg/websocket"
)

type Store struct {
	token 		TokenOptions
	websocket 	*websocket.Websocket
}

func NewStore(opt Options) *Store {
	return &Store{
		token: 		opt.Token,
		websocket: 	websocket.New(),
	}
}

type Options struct {
	Token TokenOptions
}

type TokenOptions struct {
	SessionID *token.Jwt
}