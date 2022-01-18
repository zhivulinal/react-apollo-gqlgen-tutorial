package token

import (
	"fmt"
	model "react-apollo-gqlgen-tutorial/backoffice/models"
)

type Claims struct {
	Secret 		string
	AccessToken string
	Session 	*model.Session
}

func (c *Claims) ValidateClientID(cid string) error {
	if c.Session == nil {
		return fmt.Errorf("session not found")
	}
	if c.Session.GetClientID() != cid {
		return fmt.Errorf("incorrect client ID")
	}
	return nil
}

func (c *Claims) ValidateSession() error {
	if c.Session == nil {
		return fmt.Errorf("token invalid")
	}
	return nil
}

func (c *Claims) ValidateAccessToken() error {
	if c.AccessToken == "" {
		return fmt.Errorf("access_token not found")
	}
	return nil
}

func (c *Claims) GetAccessToken() (string, error) {
	if c.AccessToken == "" {
		return "", fmt.Errorf("access_token not found")
	}
	return c.AccessToken, nil
}

func (c *Claims) GetSecret() (string, error) {
	if c.Secret == "" {
		return "", fmt.Errorf("secret not found")
	}
	return c.Secret, nil
}

func (c *Claims) GetSession() (*model.Session, error) {
	if c.Session == nil {
		return nil, fmt.Errorf("session not found")
	}
	return c.Session, nil
}
