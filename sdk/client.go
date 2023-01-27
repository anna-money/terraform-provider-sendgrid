package sendgrid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
)

const (
	defaultBaseURL = "https://api.sendgrid.com/v3/"
)

// Client is a Sendgrid client.
type Client struct {
	BaseURL   *url.URL
	UserAgent string

	apiKey     string
	host       string
	OnBehalfOf string
}

type Response struct {
	*http.Response

	// For APIs that support cursor pagination, the following field will be populated
	// to point to the next page if more results are available.
	// Set ListCursorParams.Cursor to this value when calling the endpoint again.
	Cursor string

	Rate Rate
}

type Rate struct {
	// The maximum number of requests allowed within the window.
	Limit int

	// The number of requests this caller has left on this endpoint within the current window
	Remaining int

	// The time when the next rate limit window begins and the count resets, measured in UTC seconds from epoch
	Reset time.Time
}

// NewClient creates a Sendgrid Client.
func NewClient(apiKey, host, onBehalfOf string) *Client {
	if host == "" {
		host = defaultBaseURL
	}

	return &Client{
		apiKey:     apiKey,
		host:       host,
		OnBehalfOf: onBehalfOf,
	}
}

func bodyToJSON(body interface{}) ([]byte, error) {
	if body == nil {
		return nil, ErrBodyNotNil
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("body could not be jsonified: %w", err)
	}

	return jsonBody, nil
}

// Get gets a resource from Sendgrid.
func (c *Client) Get(ctx context.Context, method rest.Method, endpoint string) (string, int, error) {
	var req rest.Request
	if c.OnBehalfOf != "" {
		req = sendgrid.GetRequestSubuser(c.apiKey, endpoint, c.host, c.OnBehalfOf)
	} else {
		req = sendgrid.GetRequest(c.apiKey, endpoint, c.host)
	}

	req.Method = method

	resp, err := sendgrid.API(req)
	if err != nil || resp.StatusCode >= 400 {
		return "", resp.StatusCode, fmt.Errorf("api response: HTTP %d: %s, err: %v", resp.StatusCode, resp.Body, err)
	}

	return resp.Body, resp.StatusCode, nil
}

// Post posts a resource to Sendgrid.
func (c *Client) Post(ctx context.Context, method rest.Method, endpoint string, body interface{}) (string, int, error) {
	var err error

	var req rest.Request

	if c.OnBehalfOf != "" {
		req = sendgrid.GetRequestSubuser(c.apiKey, endpoint, c.host, c.OnBehalfOf)
	} else {
		req = sendgrid.GetRequest(c.apiKey, endpoint, c.host)
	}

	req.Method = method

	if body != nil {
		req.Body, err = bodyToJSON(body)
	}

	if err != nil {
		return "", 0, fmt.Errorf("failed preparing request body: %w", err)
	}

	resp, err := sendgrid.API(req)

	if resp.StatusCode >= 400 {
		return "", resp.StatusCode, fmt.Errorf("api response: HTTP %d: %s", resp.StatusCode, resp.Body)
	}
	if err != nil {
		return "", resp.StatusCode, fmt.Errorf("api send post error: %v", err)
	}

	return resp.Body, resp.StatusCode, nil
}
