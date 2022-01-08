package websocket

import model "react-apollo-gqlgen-tutorial/backoffice/models"

type observer struct {
	auth chan *model.Auth
	user chan *model.User
}
