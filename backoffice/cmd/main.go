package main

import (
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"react-apollo-gqlgen-tutorial/backoffice/pkg/graph"
	"react-apollo-gqlgen-tutorial/backoffice/pkg/middleware"
	"react-apollo-gqlgen-tutorial/backoffice/pkg/token"
	st "react-apollo-gqlgen-tutorial/backoffice/pkg/store"
)

var (
	defaultPort = "2000"

	clientSessionToken = token.NewJwt(token.Options{
		SecretKey: "super-secret-001",
	})
	userToken = token.NewJwt(token.Options{
		SecretKey: "super-secret-002",
		ExpSeconds: 5,
	})
	authToken = token.NewBase64(token.Options{
		SecretKey: "super-secret-003",
	})
	urlToken = token.NewBase64(token.Options{
		SecretKey: "super-secret-004",
		ExpSeconds: 5 * 60,
	})
	accessToken = token.NewBase64(token.Options{
		SecretKey: "super-secret-003",
	})
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Создадим Store
	store := st.NewStore(st.Options{
		Token: st.TokenOptions{
			ClientSession: 	clientSessionToken,
			AuthToken: 		authToken,
			UrlToken: 		urlToken,
			UserToken: 		userToken,
			AccessToken:	accessToken,
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

	// Путь для подтверждения сессии по URL
	router.HandleFunc("/accept_auth", func(w http.ResponseWriter, r *http.Request) {
		store.AuthAcceptSessionHandleHTTP(w, r)
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
