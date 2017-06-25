package auth

import (
	"net/http"
	"os"
	"net/http/cookiejar"
	"encoding/json"
	"net/url"
	"gopkg.in/headzoo/surf.v1/errors"
	"io"
)

type CookieCredentials struct {
	client *http.Client
	APIToken string
	EnableToken int
}

func NewCookieCredentialsFromJson(in io.Reader) (*CookieCredentials, error) {
	cookieJar, _ := cookiejar.New(nil)

	cookiesUrl, err := url.Parse("https://photos.google.com")
	if err != nil {
		return nil, err
	}

	cookies := []*http.Cookie{}
	json.NewDecoder(in).Decode(&cookies)
	cookieJar.SetCookies(cookiesUrl, cookies)

	return &CookieCredentials{
		client: &http.Client{
			Jar: cookieJar,
		},
	}, nil
}

func NewCookieCredentialsFromFile(fileName string) (*CookieCredentials, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.New("Can't open cookie file")
	}
	defer file.Close()

	return NewCookieCredentialsFromJson(file)
}

func (c *CookieCredentials) GetClient() *http.Client {
	return c.client
}

func (c *CookieCredentials) GetAPIToken() string {
	return c.APIToken
}

func (c *CookieCredentials) SetAPIToken(token string) {
	c.APIToken = token
}

func (c *CookieCredentials) GetEnableNumber() int {
	return c.EnableToken
}

func (c *CookieCredentials) SetEnableNumber(number int) {
	c.EnableToken = number
}