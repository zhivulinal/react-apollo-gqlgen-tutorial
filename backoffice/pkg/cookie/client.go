package cookie

import (
	"github.com/google/uuid"
	"net/http"
)

type Client struct {
	response 	http.ResponseWriter
	request 	*http.Request
	name 		string
}

func NewClient(w http.ResponseWriter, r *http.Request) *Client {
	return &Client{
		response: w,
		request: r,
		name: "_cid",
	}
}

func (c *Client) setCookie() string {
	cid := uuid.New().String()

	// Создадим cookie
	cookie := &http.Cookie{
		Name: c.name,
		Value: cid,
		HttpOnly: true,
		//Secure: true,
	}

	// Установим cookie
	http.SetCookie(c.response, cookie)

	return cid
}

func (c *Client) Get() (clientID string) {

	// Попробуем получить ClientID из cookie
	cookie, err := c.request.Cookie(c.name)
	if err != nil {

		// Создадим ClientID
		return c.setCookie()
	}

	return cookie.Value
}


