package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Создадим структуру Jwt
type Jwt struct {
	SecretKey  string
	Issuer     string
	Expiration int64
}

func newJwt(opt Options) *Jwt {
	return &Jwt{
		SecretKey: 	opt.SecretKey,
		Issuer: 	opt.Issuer,
		Expiration: opt.ExpSeconds,
	}
}

type JwtOptions struct {
	SecretKey 	string
	Issuer 		string
	ExpSeconds 	int64
}

func (j *Jwt) Generate(opt Claims) (string, error) {
	if err := opt.ValidateSession(); err != nil {
		return "", err
	}

	// Создадим JwtClaims
	claims := &JwtClaims{}

	// Добавим Claims к JwtClaims
	claims.AddClaims(opt)

	// Инициализация StandardClaims
	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Local().Add(time.Second * time.Duration(j.Expiration)).Unix(),
		Issuer:    j.Issuer,
	}

	// Генерация токена
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return t.SignedString([]byte(j.SecretKey))
}

func (j *Jwt) Validate(opt ValidateOptions) (*Claims, error) {
	if opt.Token == "" {
		return nil, fmt.Errorf("incorrect token")
	}

	// Попробуем получить полезную нагрузку
	token, _ := jwt.ParseWithClaims(
		opt.Token,
		&JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)

	// Полезной нагрузки нет
	// Что-то явно не валидное – вернем ошибку
	if token == nil {
		return nil, fmt.Errorf("token invalid")
	}

	// Пробуем получить Claims
	claims, ok := token.Claims.(*JwtClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Claims должен содержать сессию
	if err := claims.ValidateSession(); err != nil {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Дополнительная валидация из ValidateOptions
	switch {
	case opt.Cid != "":
		if err := claims.ValidateClientID(opt.Cid); err != nil {
			return nil, err
		}
	case opt.AccessToken:
		if err := claims.ValidateAccessToken(); err != nil {
			return nil, err
		}
	}

	// Проверим срок жизни токена
	if j.Expiration > 0 && claims.ExpiresAt < time.Now().Local().Unix() {

		// Если токен протух
		// Вернем полезную нагрузку вместе с ошибкой
		return claims.GetAccessToken(), fmt.Errorf("token is expired")
	}

	return claims.GetClaims(), nil
}

//
type JwtClaims struct {
	Claims
	jwt.StandardClaims
}

func (j *JwtClaims) AddClaims(opt Claims) {
	j.Claims = opt
}

func (j *JwtClaims) GetClaims() *Claims {
	return &j.Claims
}

func (j *JwtClaims) GetAccessToken() *Claims {
	if j.Claims.AccessToken == "" {
		return nil
	}
	return &Claims{
		AccessToken: j.Claims.AccessToken,
	}
}
