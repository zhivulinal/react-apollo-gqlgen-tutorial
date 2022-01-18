package store

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
	"react-apollo-gqlgen-tutorial/backoffice/pkg/cookie"
	"react-apollo-gqlgen-tutorial/backoffice/pkg/token"
	"strconv"
	"time"
)

// Создает cookie с токеном сессии Пользователя
func (s *Store) sessionCreateUserToken(w http.ResponseWriter, r *http.Request, sess *model.Session) {

	sid := sess.GetSessionID()
	cookieUser := cookie.NewUser(w, r)

	// Пользователь был авторизован ранее
	// Нужно создать accessToken, токен и cookie

	// Создадим AccessToken, запишем в него SessionID
	accessToken, err := s.token.AccessToken.Generate(token.Claims{
		Secret: sid,
	})
	if err != nil {
		return
	}

	// Создадим токен Пользователя
	userToken, err := s.token.UserToken.Generate(token.Claims{
		AccessToken: accessToken,
		Session:     sess,
	})
	if err != nil {
		return
	}

	// Установим cookie
	cookieUser.Add(userToken)
}

// Валидирует cookie Пользователя
func (s *Store) sessionValidateUser(w http.ResponseWriter, r *http.Request) *http.Request {
	cookieUser := cookie.NewUser(w, r)

	// Найдем cookie Пользователя
	userToken, ok := cookieUser.Get()
	if !ok {
		return r
	}

	// Валидируем токен
	claims, err := s.token.UserToken.Validate(token.ValidateOptions{
		Token: userToken,

		// Дополнительная валидация
		// Claims должен включать AccessToken
		AccessToken: true,
	})

	if claims == nil {
		return r
	}

	ctx := r.Context()

	// Токен не валиден.
	// Если есть Claims, значит просрочен
	if err != nil {

		// Получим AccessToken
		accessToken, _ := claims.GetAccessToken()

		// Валидируем AccessToken
		cl, err2 := s.token.AccessToken.Validate(token.ValidateOptions{
			Token: accessToken,
		})
		if err2 != nil {

			// Удалим cookie
			cookieUser.Delete()
			return r
		}

		// Извлечем Cid из Claims
		sid, err3 := cl.GetSecret()
		if err3 != nil {

			// Удалим cookie
			cookieUser.Delete()
			return r
		}

		// Получим сессию
		sess, err4 := s.repo.Session.GetInstalled(sid)
		if err4 != nil {

			// Удалим cookie
			cookieUser.Delete()
			return r
		}

		// Перезапишем cookie
		s.sessionCreateUserToken(w, r, sess)

		// Сохраним сессию в контексте
		return r.WithContext(sess.WithContext(ctx))
	}

	sess, err := claims.GetSession()
	if err != nil {

		// Удалим cookie
		cookieUser.Delete()
		return r
	}

	return r.WithContext(sess.WithContext(ctx))
}

// Работает с сессией пользователя
func (s *Store) SessionHandleUser(w http.ResponseWriter, r *http.Request) *http.Request {
	ctx := r.Context()

	// Получим сессию из контекста
	sess, err := model.SessionFromContext(ctx)
	if err != nil {
		return r
	}

	if ok := sess.CheckOnline(); !ok {
		return r
	}

	// Проверим авторизацию Пользователя
	// Могла быть активирована сессия
	// в SessionHandleApproved
	if ok := sess.CheckUserAuthorization(); ok {
		s.sessionCreateUserToken(w, r, sess)
	} else {
		return s.sessionValidateUser(w, r)
	}

	return r.WithContext(ctx)
}

// Обрабатывает авторизацию Пользователя
func (s *Store) SessionHandleApproved(w http.ResponseWriter, r *http.Request) *http.Request {

	// Ищем заголовок Authorization
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		return r
	}

	ctx := r.Context()

	// Валидируем токен
	claims, err := s.token.AuthToken.Validate(token.ValidateOptions{
		Token: authToken,
	})
	if err != nil {
		return r
	}

	// Заберем секрет из Claims
	secret, err := claims.GetSecret()
	if err != nil {
		return r
	}

	// Ищем сессию
	sess, err := s.repo.Session.GetApproved(secret)
	if err != nil {
		return r
	}

	// Сохраним сессию в контекст
	return r.WithContext(sess.WithContext(ctx))
}

// Валидация кода из СМС
func (s *Store) SessionValidateSMSCode(code string) (*model.Session, error) {
	sess, err := s.repo.Session.GetByCode(code)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

// Валидация URL кода
func (s *Store) SessionValidateURLCode(code string) (*model.Session, error) {

	// Код является токеном
	// Валидируем токен, извлекаем SessionID
	claims, err := s.token.UrlToken.Validate(token.ValidateOptions{
		Token: code,
	})
	if err != nil {
		return nil, err
	}

	// Заберем секрет из Claims
	secret, err := claims.GetSecret()
	if err != nil {
		return nil, err
	}

	// Получим сессию
	sess, err := s.repo.Session.GetPending(secret)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

// Создает код для активации по Url
func (s *Store) SessionCreateURLCode(ctx context.Context) (url, authToken string, err error) {

	// Получим сессию из контекста
	sess, err := model.SessionFromContext(ctx)
	if err != nil {
		return "", "", err
	}

	// Получим SessionID
	sid := sess.GetSessionID()

	// Создадим URL токен, запишем в него SessionID
	url, err = s.token.UrlToken.Generate(token.Claims{
		Secret: sid,
	})

	// Создадим токен авторизации
	authToken, err = s.token.AuthToken.Generate(token.Claims{
		Secret: sid,
	})
	if err != nil {
		return "", "", err
	}

	// Добавим authToken к сессии
	sess.SetAuthToken(authToken)

	// Сохраним сессию в репозитории
	s.repo.Session.Add(sess)

	return url, authToken, nil
}

// Создает код для активации по Смс
func (s *Store) SessionWithSmsCode(ctx context.Context) (code, authToken string, err error) {

	// Получим сессию из контекста
	sess, err := model.SessionFromContext(ctx)
	if err != nil {
		return "", "", err
	}

	// Создадим и добавим СМС-код к сессии
	code = createSmsCode()
	sess.SetSecretCode(code)

	// Создадим токен авторизации
	sid := sess.GetSessionID()
	authToken, err = s.token.AuthToken.Generate(token.Claims{
		Secret: sid,
	})
	if err != nil {
		return "", "", err
	}

	// Добавим authToken к сессии
	sess.SetAuthToken(authToken)

	// Сохраним сессию в репозитории
	s.repo.Session.Add(sess)

	return code, authToken, nil
}

// Валидирует токен слушателя
// Окно браузера
func (s *Store) ValidateSessionToken(tkn string) (*model.Session, error) {
	if claims, err := s.token.ClientSession.Validate(token.ValidateOptions{
		Token: tkn,
	}); err == nil {
		sess := claims.Session

		// Сессию получили из заголовка: клиент онлайн
		sess.SetOnline()

		return sess, nil
	}

	return nil, fmt.Errorf("invalid session token")
}

// Обрабатывает сессию клиента
func (s *Store) SessionHandleClient(w http.ResponseWriter, r *http.Request) *http.Request {

	// Получим контекст
	ctx := r.Context()

	// Получим ClientID
	cookieCid := cookie.NewClient(w, r)
	cid := cookieCid.Get()

	var sess *model.Session

	// Ищем заголовок Session-ID
	if tkn := r.Header.Get("Session-ID"); tkn != "" {

		// Нашли сессию
		if _sess, err2 := s.ValidateSessionToken(tkn); err2 == nil {
			sess = _sess
		}
	}

	if sess == nil {
		sess = model.NewSession(cid)
	}

	// Сохраним сессию в контекст и вернем *http.Request
	return r.WithContext(sess.WithContext(ctx))
}

func createSmsCode() string {
	rand.Seed(time.Now().UnixNano())

	form := func(s string) string {
		if len(s) == 2 {
			s += "0"
		} else if len(s) == 1 {
			s += "00"
		}
		return s
	}

	return form(strconv.Itoa(rand.Intn(1000))) + form(strconv.Itoa(rand.Intn(1000)))
}
