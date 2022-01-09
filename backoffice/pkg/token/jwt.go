package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
	"time"
)

// Создадим структуру Jwt
type Jwt struct {
	SecretKey  string
	Issuer     string
	Expiration int64
}

// Опции при генерации токена
type JwtClaims struct {

	// Нужен для обновления токена
	AccessToken string

	// Прикрепим сессию
	Sess *model.Session

	jwt.StandardClaims
}

// Генерация токена
func (j *Jwt) Generate(opt JwtClaims) (token string, err error) {

	// Получим Claims
	claims := &opt

	// Инициализация StandardClaims
	//
	// Здесь "подключаются" все настройки
	// Необходимые для валидации токена
	//
	// Указываются при инициализации структуры Jwt
	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Local().Add(time.Second * time.Duration(j.Expiration)).Unix(),
		Issuer:    j.Issuer,
	}

	// Генерация токена
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return t.SignedString([]byte(j.SecretKey))
}

// Опции при валидации токена
type JwtValidateOptions struct {
	Token 	string
}

// Валидация токена
func (j *Jwt) Validate(opt JwtValidateOptions) (claims *JwtClaims, err error) {

	// Попробуем получить полезную нагрузку
	token, err := jwt.ParseWithClaims(
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

	// Получим Claims
	claims, ok := token.Claims.(*JwtClaims)
	if !ok {
		return nil, fmt.Errorf("error token claims")
	}

	// Проверим срок жизни токена
	if j.Expiration > 0 && claims.ExpiresAt < time.Now().Local().Unix() {

		// Если токен протух
		// Вернем полезную нагрузку вместе с ошибкой
		//
		// Для дальнейшей валидации:
		// claims будет содержать AccessToken
		return claims, fmt.Errorf("token is expired")
	}

	return claims, nil
}

// Опции структуры Jwt
type JwtOptions struct {
	SecretKey 	string
	Issuer 		string
	ExpSeconds 	int64
}

func NewJwt(opt JwtOptions) *Jwt {
	return &Jwt{
		SecretKey: 	opt.SecretKey,
		Issuer: 	opt.Issuer,
		Expiration: opt.ExpSeconds,
	}
}
