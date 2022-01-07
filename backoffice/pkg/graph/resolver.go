package graph

import (
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"
	"net/http"
	"react-apollo-gqlgen-tutorial/backoffice/graph/generated"
	"react-apollo-gqlgen-tutorial/backoffice/pkg/store"
	"time"
)

var (
	mb int64 = 1 << 20
)

type Resolver struct{
	store *store.Store
}

// Создадим функцию NewServer
func NewServer(opt Options) *handler.Server {

	// Переместим создание сервера из cmd/main.go
	srv := handler.New(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &Resolver{
					store: opt.Store,
				},
			},
		),
	)
	srv.AddTransport(transport.MultipartForm{
		MaxMemory:     32 * mb,
		MaxUploadSize: 50 * mb,
	})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		InitFunc: transport.WebsocketInitFunc(func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
			return ctx, nil
		}),
	})
	srv.Use(extension.Introspection{})

	return srv
}

type Options struct {
	Store *store.Store
}
