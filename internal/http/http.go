package http

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	client *http.Client
}

// New returns a new instance of Client
func New(client *http.Client) Client {
	return Client{
		client: client,
	}
}

// DoRequest sends GET request ot provided URL and returns response body
func (t Client) DoRequest(u *url.URL) ([]byte, error) {
	resp, err := t.client.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("http request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request %s failed with status_code %d", u.String(), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error: %v", err)
	}

	return body, nil
}
