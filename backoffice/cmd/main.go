package main

import (
	"github.com/99designs/gqlgen/graphql/playground"
	"log"
	"net/http"
	"os"
	"react-apollo-gqlgen-tutorial/backoffice/pkg/graph"
	"github.com/gorilla/mux"
	"react-apollo-gqlgen-tutorial/backoffice/pkg/middleware"
)

var (
	defaultPort = "2000"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Создадим GraphQL сервер
	srv := graph.NewServer(graph.Options{})

	// Создадим роутер
	router := mux.NewRouter()

	// Подключим CORS middleware
	router.Use(middleware.CorsMiddleware())

	router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	router.Handle("/graphql", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
