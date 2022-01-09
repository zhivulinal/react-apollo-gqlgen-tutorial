package main

import (
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"react-apollo-gqlgen-tutorial/backoffice/pkg/graph"
	"react-apollo-gqlgen-tutorial/backoffice/pkg/middleware"
	st "react-apollo-gqlgen-tutorial/backoffice/pkg/store"
	"react-apollo-gqlgen-tutorial/backoffice/pkg/token"
)

var (
	defaultPort = "2000"
)

func main() {
	tokenSessID := token.NewJwt(token.JwtOptions{
		SecretKey: "ololo",
		Issuer: "example.com",
		ExpSeconds: 0,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Создадим Store
	store := st.NewStore(st.Options{
		Token: st.TokenOptions{
			SessionID: tokenSessID,
		},
	})

	// Создадим GraphQL сервер
	srv := graph.NewServer(graph.Options{

		// Подключим стор в qraphql сервер
		Store: store,
	})

	// Создадим роутер
	router := mux.NewRouter()

	// Подключим CORS middleware
	router.Use(middleware.CorsMiddleware())

	// Подключим Auth middleware и передадим store в качестве параметра
	router.Use(middleware.AuthMiddleware(store))

	router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	router.Handle("/graphql", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
