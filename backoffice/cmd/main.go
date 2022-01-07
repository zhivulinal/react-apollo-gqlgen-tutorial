package main

import (
	"github.com/99designs/gqlgen/graphql/playground"
	"log"
	"net/http"
	"os"
	"react-apollo-gqlgen-tutorial/backoffice/pkg/graph"
	"github.com/gorilla/mux"
	"react-apollo-gqlgen-tutorial/backoffice/pkg/middleware"
	st "react-apollo-gqlgen-tutorial/backoffice/pkg/store"
)

var (
	defaultPort = "2000"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Создадим Store
	store := st.NewStore(st.Options{})

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
