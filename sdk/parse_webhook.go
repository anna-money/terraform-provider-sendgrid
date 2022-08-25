package sendgrid

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ParseWebhook is a Sendgrid inbound parse settings.
type ParseWebhook struct {
	Hostname  string `json:"hostname,omitempty"`
	URL       string `json:"url,omitempty"`
	SpamCheck bool   `json:"spam_check"` //nolint:tagliatelle
	SendRaw   bool   `json:"send_raw"`   //nolint:tagliatelle
}

func parseParseWebhook(respBody string) (*ParseWebhook, RequestError) {
	var body ParseWebhook
	if err := json.Unmarshal([]byte(respBody), &body); err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed parsing inbound parse: %w", err),
		}
	}

	return &body, RequestError{StatusCode: http.StatusOK, Err: nil}
}

// CreateParseWebhook creates an ParseWebhook and returns it.
func (c *Client) CreateParseWebhook(
	hostname string,
	url string,
	spamCheck bool,
	sendRaw bool) (*ParseWebhook, RequestError) {
	if hostname == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrHostnameRequired,
		}
	}

	if url == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrURLRequired,
		}
	}

	respBody, statusCode, err := c.Post("POST", "/user/webhooks/parse/settings", ParseWebhook{
		Hostname:  hostname,
		URL:       url,
		SpamCheck: spamCheck,
		SendRaw:   sendRaw,
	})
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed creating inbound parse: %w", err),
		}
	}

	if statusCode >= http.StatusMultipleChoices {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w, status: %d, response: %s", ErrFailedCreatingParseWebhook, statusCode, respBody),
		}
	}

	return parseParseWebhook(respBody)
}

// ReadParseWebhook retreives an ParseWebhook and returns it.
func (c *Client) ReadParseWebhook(hostname string) (*ParseWebhook, RequestError) {
	if hostname == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrHostnameRequired,
		}
	}

	respBody, _, err := c.Get("GET", "/user/webhooks/parse/settings/"+hostname)
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return parseParseWebhook(respBody)
}

// UpdateParseWebhook edits an ParseWebhook and returns it.
func (c *Client) UpdateParseWebhook(hostname string, spamCheck bool, sendRaw bool) RequestError {
	if hostname == "" {
		return RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrHostnameRequired,
		}
	}

	t := ParseWebhook{}
	t.SpamCheck = spamCheck
	t.SendRaw = sendRaw

	_, _, err := c.Post("PUT", "/user/webhooks/parse/settings/"+hostname, t)
	if err != nil {
		return RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return RequestError{
		StatusCode: http.StatusOK,
	}
}

// DeleteParseWebhook deletes an ParseWebhook.
func (c *Client) DeleteParseWebhook(hostname string) (bool, RequestError) {
	if hostname == "" {
		return false, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrHostnameRequired,
		}
	}

	responseBody, statusCode, err := c.Get("DELETE", "/user/webhooks/parse/settings/"+hostname)
	if err != nil {
		return false, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if statusCode >= http.StatusMultipleChoices && statusCode != http.StatusNotFound { // ignore not found
		return false, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w, status: %d, response: %s", ErrFailedDeletingParseWebhook, statusCode, responseBody),
		}
	}

	return true, RequestError{StatusCode: http.StatusOK, Err: nil}
}
