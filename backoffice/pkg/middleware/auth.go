package middleware

import (
	"net/http"
	"react-apollo-gqlgen-tutorial/backoffice/pkg/store"
)

func AuthMiddleware(store *store.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Метод из Store, обрабатывает логику авторизации
			r = store.HandleAuthHTTP(w, r)
			next.ServeHTTP(w, r)
		})
	}
}
