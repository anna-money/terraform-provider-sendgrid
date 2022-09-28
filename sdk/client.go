package sendgrid

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
)

const (
	defaultBaseURL = "https://api.sendgrid.com/v3/"
	userAgent      = "go-sendgrid"

	// https://docs.sendgrid.com/v2-api/using_the_web_api#rate-limits
	headerRateLimit     = "X-Ratelimit-Limit"
	headerRateRemaining = "X-Ratelimit-Remaining"
	headerRateReset     = "X-Ratelimit-Reset"
)

var errNonNilContext = errors.New("context must be non-nil")

// Client is a Sendgrid client.
type Client struct {
	client    *http.Client
	BaseURL   *url.URL
	UserAgent string

	apiKey     string
	host       string
	OnBehalfOf string
}

type ErrorResponse struct {
	Response *http.Response
	Detail   string `json:"detail"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf(
		"%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Detail)
}

func (r *ErrorResponse) Is(target error) bool {
	v, ok := target.(*ErrorResponse)
	if !ok {
		return false
	}
	if r.Detail != v.Detail ||
		!matchHTTPResponse(r.Response, v.Response) {
		return false
	}
	return true
}

type RateLimitError struct {
	Rate     Rate
	Response *http.Response
	Detail   string
}

func (r *RateLimitError) Error() string {
	return fmt.Sprintf(
		"%v %v: %d %v %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Detail, fmt.Sprintf("[rate reset in %v]", time.Until(r.Rate.Reset)))
}
func (r *RateLimitError) Is(target error) bool {
	v, ok := target.(*RateLimitError)
	if !ok {
		return false
	}
	return r.Rate == v.Rate &&
		r.Detail == v.Detail &&
		matchHTTPResponse(r.Response, v.Response)
}

// matchHTTPResponse compares two http.Response objects. Currently, only StatusCode is checked.
func matchHTTPResponse(r1, r2 *http.Response) bool {
	if r1 == nil && r2 == nil {
		return true
	}
	if r1 != nil && r2 != nil {
		return r1.StatusCode == r2.StatusCode
	}
	return false
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
	if err != nil {
		return "", resp.StatusCode, fmt.Errorf("failed getting resource: %w", err)
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
	if err != nil {
		return "", resp.StatusCode, fmt.Errorf("failed posting resource: %w", err)
	}

	return resp.Body, resp.StatusCode, nil
}

// ParseRate parses the rate limit headers.
func ParseRate(r *http.Response) Rate {
	var rate Rate
	if limit := r.Header.Get(headerRateLimit); limit != "" {
		rate.Limit, _ = strconv.Atoi(limit)
	}
	if remaining := r.Header.Get(headerRateRemaining); remaining != "" {
		rate.Remaining, _ = strconv.Atoi(remaining)
	}
	if reset := r.Header.Get(headerRateReset); reset != "" {
		if v, _ := strconv.ParseInt(reset, 10, 64); v != 0 {
			rate.Reset = time.Unix(v, 0).UTC()
		}
	}

	return rate
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		apiError := new(APIError)
		json.Unmarshal(data, apiError)
		if apiError.Empty() {
			errorResponse.Detail = strings.TrimSpace(string(data))
		} else {
			errorResponse.Detail = apiError.Detail()
		}
	}
	// Re-populate error response body.
	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	switch {
	case r.StatusCode == http.StatusTooManyRequests &&
		r.Header.Get(headerRateRemaining) == "0":
		return &RateLimitError{
			Rate:     ParseRate(r),
			Response: errorResponse.Response,
			Detail:   errorResponse.Detail,
		}
	}

	return errorResponse
}

func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	response.Rate = ParseRate(r)
	return response
}

func (c *Client) BareDo(ctx context.Context, req *http.Request) (*Response, error) {
	if ctx == nil {
		return nil, errNonNilContext
	}

	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return nil, err
		}
	}

	response := newResponse(resp)
	err = CheckResponse(resp)
	return response, err
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.BareDo(ctx, req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		dec := json.NewDecoder(resp.Body)
		dec.UseNumber()
		decErr := dec.Decode(v)
		if decErr == io.EOF {
			decErr = nil
		}
		if decErr != nil {
			err = decErr
		}
	}
	return resp, err
}

func (c *Client) NewRequest(method, urlRef string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}

	u, err := c.BaseURL.Parse(urlRef)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}
