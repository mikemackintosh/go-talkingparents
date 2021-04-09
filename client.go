package tp

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
)

const (
	DATE_FORMAT      = "2006-01-02T15:04:05"
	SENT_MESSAGE int = 1
	VIEW_MESSAGE int = 2

	ENDPOINT_DOMAIN        = "https://talkingparents.com"
	ENDPOINT_LOGIN         = ENDPOINT_DOMAIN + "/login.aspx"
	ENDPOINT_CONVERSATIONS = ENDPOINT_DOMAIN + "/api/Conversations"
	ENDPOINT_THREADS       = ENDPOINT_DOMAIN + "/api/Threads"

	USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"
)

// ErrNoClientConfigured is used when no client is configured.
type ErrNoClientConfigured struct{}

// Error matches the error interface.
func (e ErrNoClientConfigured) Error() string {
	return "need to create a client with tp.NewClient() first"
}

// ErrInvalidUsernameOrPassword is used when an invalid username or password was used during authentication.
type ErrInvalidUsernameOrPassword struct{}

// Error matches the error interface.
func (e ErrInvalidUsernameOrPassword) Error() string {
	return "invalid username or password"
}

// Client is a client struct for tp.
type Client struct {
	httpClient *http.Client
}

// NewClient will generate a new http.Client with the use of a cookiejar, since form-based auth is used.
func NewClient() (*Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: jar,
	}

	return &Client{httpClient: client}, nil
}

// Authenticate takes a username and password and returns true/false | error.
func (c *Client) Authenticate(username, password string) (error) {
	var res *http.Response
	var err error

	// validate if a client was configured or not.
	if c.httpClient == nil {
		return ErrNoClientConfigured{}
	}

	// Create the first grab of the login page.
	if res, err = c.httpClient.Get(ENDPOINT_LOGIN); err != nil {
		return err
	}

	// Read the CSRF tokens and other required fields for successful auth.
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	data := url.Values{}
	doc.Find("input[type='hidden']").Each(func(i int, s *goquery.Selection) {
		if name, ok := s.Attr("name"); ok {
			if val, ok := s.Attr("value"); ok {
				data.Set(name, val)
			}
		}
	})

	// set additional required authentication fields
	data.Set("__EVENTARGUMENT", "")
	data.Set("__EVENTTARGET", "LoginForm$LoginButton")
	data.Set("LoginForm$UserName", username)
	data.Set("LoginForm$Password", password)

	// create the authentication request
	loginRequest, err := http.NewRequest("POST", ENDPOINT_LOGIN, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	// Set the appropriate headers
	loginRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	loginRequest.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	loginRequest.Header.Add("Origin", ENDPOINT_DOMAIN)
	loginRequest.Header.Add("Referer", ENDPOINT_LOGIN)
	loginRequest.Header.Add("User-Agent", USER_AGENT)

	res, err = c.httpClient.Do(loginRequest)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return ErrInvalidUsernameOrPassword{}
	}
	defer res.Body.Close()

	return nil
}
