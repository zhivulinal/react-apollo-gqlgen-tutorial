package token

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// О шифровании Base64
// https://stackoverflow.com/questions/25290596/go-aes128-without-iv
type Base64 struct {
	Key []byte
	Expiration int64
}

func newBase64(opt Options) *Base64 {
	if ok := opt.ValidateSecretKeyFormBase64(); !ok {
		panic("invalid secret key")
	}

	return &Base64{
		Key: []byte(opt.SecretKey),
		Expiration: opt.ExpSeconds,
	}
}

func (b *Base64) Generate(opt Claims) (string, error) {

	// Создадим Claims
	claims := &Bas64Claims{
		ExpSeconds: b.Expiration,
	}
	plainText, err := claims.Create(b.Expiration, opt)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(b.Key)
	if err != nil {
		panic(err)
	}
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func (b *Base64) Validate(opt ValidateOptions) (*Claims, error) {
	if opt.Token == "" {
		return nil, fmt.Errorf("incorrect token")
	}

	cipherText, err := base64.URLEncoding.DecodeString(opt.Token)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(b.Key)
	if err != nil {
		return nil, err
	}
	if len(cipherText) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	// Создадим Claims
	claims := &Bas64Claims{}

	return claims.Get(fmt.Sprintf("%s", cipherText))
}

func (j *Bas64Claims) Get(claims string) (*Claims, error) {
	err := json.Unmarshal([]byte(claims), j)
	if err != nil {
		return nil, err
	}

	if j.ExpSeconds > 0 && j.ExpSeconds < time.Now().Local().Unix() {
		return nil, fmt.Errorf("token expired")
	}

	cl := &Claims{}

	if j.Secret != "" {
		cl.Secret = j.Secret
	}

	return cl, nil
}

//
type Bas64Claims struct {
	Secret 	string
	ExpSeconds 	int64
}

func (b *Bas64Claims) Create(expiration int64, opt Claims) ([]byte, error) {
	if opt.Secret != "" {
		b.Secret = opt.Secret
	}

	if b.ExpSeconds > 0 {
		b.ExpSeconds = time.Now().Local().Add(time.Second * time.Duration(expiration)).Unix()
	}

	return json.Marshal(b)
}