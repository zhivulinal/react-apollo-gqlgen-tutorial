package token

func NewJwt(opt Options) Manager {
	return newJwt(opt)
}

func NewBase64(opt Options) Manager {
	return newBase64(opt)
}

type Manager Token

// Объявим интерфейс токена
type Token interface {
	Generate(opt Claims) (token string, err error)
	Validate(opt ValidateOptions) (claims *Claims, err error)
}

type ValidateOptions struct {
	Token 	string
	Cid 	string
	AccessToken bool
}


// SecretKey: 16 символов для Base64
type Options struct {
	SecretKey  	string
	Issuer     	string
	ExpSeconds 	int64
}

// Валидация SecretKey для Base64
// Нужно ровно 16 байт
func (o *Options) ValidateSecretKeyFormBase64() bool {
	return len([]byte(o.SecretKey)) == 16
}
