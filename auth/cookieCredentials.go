package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"

	"github.com/palantir/stacktrace"
)

// Implementation of the API Credentials interface based on cookies authentication
type CookieCredentials struct {
	Client               *http.Client
	PersistentParameters *PersistentParameters
	RuntimeParameters    *RuntimeParameters
}

// Struct that holds the persistent parameters relative to an user
type PersistentParameters struct {
	UserId string `json:"userId"`
}

// Struct that contains all the parameters that changes
type RuntimeParameters struct {
	AtToken string
}

// Structure that is serialized in JSON to store the user credentials
type AuthFile struct {
	// Cookies to make the requests
	Cookies []*http.Cookie `json:"cookies"`

	// Parameters used by the requests (user id, etc ...)
	PersistentParameters *PersistentParameters `json:"persistantParameters"`
}

// Create a new CookieCredentials given a slice of cookies and the PersistentParameters
func NewCookieCredentials(cookies []*http.Cookie, parameters *PersistentParameters) *CookieCredentials {
	// Create a cookie jar for the client
	jar, _ := cookiejar.New(nil)
	cookiesUrl, _ := url.Parse("https://photos.google.com")

	// Add cookies in the jar
	jar.SetCookies(cookiesUrl, cookies)

	return &CookieCredentials{
		Client: &http.Client{
			Jar: jar,
		},
		PersistentParameters: parameters,
		RuntimeParameters:    &RuntimeParameters{},
	}
}

// Restore an CookieCredentials object from a JSON
func NewCookieCredentialsFromJson(inputReader io.Reader) (*CookieCredentials, error) {
	// Parse AuthFile
	authFile := AuthFile{}
	if err := json.NewDecoder(inputReader).Decode(&authFile); err != nil {
		return nil, fmt.Errorf("auth: Can't read the JSON AuthFile (%v)", err)
	}

	credentials := NewCookieCredentials(authFile.Cookies, authFile.PersistentParameters)
	if invalid, err := credentials.Validate(); err != nil {
		if invalid {
			return nil, stacktrace.Propagate(err, "invalid cookie credentials")
		}
		return nil, stacktrace.Propagate(err, "failed validating credentials")
	}

	// Get a new API token using the TokenScraper from the api package
	_, err := NewAtTokenScraper(*credentials).ScrapeNewAtToken()
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed scraping AtToken")
	}

	return credentials, nil
}

// Restore an CookieCredentials object from a JSON file
func NewCookieCredentialsFromFile(fileName string) (*CookieCredentials, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("auth: Can't open %v", fileName)
	}
	defer file.Close()

	return NewCookieCredentialsFromJson(file)
}

// ValidateCredentials checks validity/format of authentication cookies.
// err may contain invalidity details, or explain what went wrong in the execution of the function
// invalid only gets switched to true if the cookies are invalid for sure, not if the function goes wrong
func (c *CookieCredentials) Validate() (invalid bool, err error) {
	// To check if the cookies are valid, make a request to the Google Photos Login and check if we're redirected
	req, err := http.NewRequest("GET", "https://photos.google.com/login", nil)
	if err != nil {
		return false, err
	}
	res, err := c.Client.Do(req)
	if err != nil {
		return false, stacktrace.Propagate(err, "request to check cookies validity failed")
	}
	if res.Request.URL.String() != "https://photos.google.com/" {
		return true, fmt.Errorf("Google didn't redirect to Photos Homepage after accessing the Login page")
	}

	return false, nil
}

// Serialize the CookieCredentials object into a JSON object, to be restored in the future using
// NewCookieCredentialsFromJson
func (c *CookieCredentials) Serialize(out io.Writer) error {
	cookiesUrl, _ := url.Parse("https://photos.google.com")
	cookies := c.Client.Jar.Cookies(cookiesUrl)

	for _, cookie := range cookies {
		if cookie.Name == "OTZ" {
			cookie.Domain = "photos.google.com"
		} else {
			cookie.Domain = ".google.com"
		}
		cookie.Path = "/"
	}

	return json.NewEncoder(out).Encode(&AuthFile{
		Cookies:              cookies,
		PersistentParameters: c.PersistentParameters,
	})
}

// SerializeToFile serializes the CookieCredentials object into a JSON file, to be restored in the future using
// NewCookieCredentialsFromJsonFile
func (c *CookieCredentials) SerializeToFile(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("auth: Can't create the file %v (%v)", fileName, err)
	}
	defer file.Close()

	return c.Serialize(file)
}
