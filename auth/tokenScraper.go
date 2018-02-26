package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

const (
	GooglePhotoUrl = "https://photos.google.com/"
)

// AtTokenScraper used to scape tokens to upload images
type AtTokenScraper struct {
	credentials CookieCredentials
}

type ApiTokenContainer struct {
	Token string `json:"SNlM0e"`
}

// Create a new scraper for the at token. This token is user-dependent, so you need to create a new token scraper
// for each Credentials object.
func NewAtTokenScraper(credentials CookieCredentials) *AtTokenScraper {
	return &AtTokenScraper{
		credentials: credentials,
	}
}

// Use this method to get a new at token. The method makes an http request to Google and uses the user credentials
func (ts *AtTokenScraper) ScrapeNewAtToken() (string, error) {
	req, err := http.NewRequest("GET", GooglePhotoUrl, nil)
	if err != nil {
		return "", fmt.Errorf("Can't create the request to get the Google Photos homepage (%v)", err)
	}

	// Make the request
	res, err := ts.credentials.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Can't complete the request to get the Google Photos homepage (%v)", err)
	}

	// Parse the request response as a HTML page
	t := html.NewTokenizer(res.Body)
	var script string
	found := false
	for !found {
		tt := t.Next()

		switch {
		case tt == html.ErrorToken: // End of html document
			return "", errors.New("Can't find the script tag with the token in the response")

		case tt == html.StartTagToken && t.Token().Data == "script": // We need the first script tag
			t.Next()

			// Get the script string
			script = t.Token().Data
			found = true
		}
	}

	// The script assigns an object to the global window object. We are going to parse the script as a JSON
	// so we need to get rid of the assignment code
	equalsIndex := strings.Index(script, "=")
	start := equalsIndex + 1
	end := len(script) - 1
	script = script[start:end]

	// Parse the json
	object := ApiTokenContainer{}
	if err = json.NewDecoder(strings.NewReader(script)).Decode(&object); err != nil {
		return "", fmt.Errorf("can't parse the JSON object that contains the at token (%v)", err)
	}

	ts.credentials.RuntimeParameters.AtToken = object.Token
	return object.Token, nil
}
