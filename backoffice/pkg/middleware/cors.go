package middleware

import "net/http"

func CorsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Разрешаем подключаться только конкретному хосту
			// Заголовко всегда должен возвращать адрес действительного хоста
			// Устанавливая звездочку "*" - создаем уязвимость в безопасности
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

			// Разрешаем принимать файлы cookie только от http://localhost:3000
			//
			// Подробнее про сочетания
			// Access-Control-Allow-Origin и Access-Control-Allow-Credentials
			// В этой таблице:
			// https://fetch.spec.whatwg.org/#cors-protocol-and-credentials
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Разрешаем использовать методы
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

			// Блокируем возможность CSRF
			// Access-Control-Allow-Headers = Content-Type
			// Подробнее здесь:
			// https://fetch.spec.whatwg.org/#concept-header
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

			if r.Method == "OPTIONS" {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}