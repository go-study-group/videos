package zoom

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Client is a zoom client
type Client struct {
	httpClient *http.Client
	Key        string
	Secret     string
	Transport  http.RoundTripper
	Timeout    time.Duration // set to value > 0 to enable a request timeout
	endpoint   string
}

// NewClient returns a new API client
func NewClient(apiKey string, apiSecret string) *Client {
	var uri = url.URL{
		Scheme: "https",
		Host:   apiURI,
		Path:   apiVersion,
	}

	return &Client{
		httpClient: http.DefaultClient,
		Key:        apiKey,
		Secret:     apiSecret,
		endpoint:   uri.String(),
	}
}

// request creates a new GET request at path
func (c *Client) get(method, path string, jsonRes interface{}) error {

	requestURL := c.endpoint + path

	client := c.httpClient
	if c.Timeout > 0 {
		client.Timeout = c.Timeout
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return err
	}

	if err := c.signWithJWT(req); err != nil {
		return err
	}

	// all v1 API endpoints use POST requests
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("expected status code 200, got %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	if jsonRes != nil {
		if json.NewDecoder(resp.Body).Decode(jsonRes); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) signWithJWT(req *http.Request) error {
	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: 15000,
		Issuer:    c.Key,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString(c.Secret)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", signedString))
	return nil
}
