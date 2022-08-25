package sendgrid

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// UnsubscribeGroup is a Sendgrid - Suppressions - Unsubscribe Group.
type UnsubscribeGroup struct {
	ID           int32  `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	IsDefault    bool   `json:"is_default"` //nolint:tagliatelle
	Unsubscribes int32  `json:"unsubscribes,omitempty"`
}

func parseUnsubscribeGroup(respBody string) (*UnsubscribeGroup, RequestError) {
	var body UnsubscribeGroup
	if err := json.Unmarshal([]byte(respBody), &body); err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed parsing API key: %w", err),
		}
	}

	return &body, RequestError{StatusCode: http.StatusOK, Err: nil}
}

func parseUnsubscribeGroups(respBody string) ([]UnsubscribeGroup, RequestError) {
	var body []UnsubscribeGroup
	if err := json.Unmarshal([]byte(respBody), &body); err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed parsing API key: %w", err),
		}
	}

	return body, RequestError{StatusCode: http.StatusOK, Err: nil}
}

// CreateUnsubscribeGroup creates an UnsubscribeGroup and returns it.
func (c *Client) CreateUnsubscribeGroup(
	name string,
	description string,
	isDefault bool) (*UnsubscribeGroup, RequestError) {
	if name == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrNameRequired,
		}
	}

	respBody, statusCode, err := c.Post("POST", "/asm/groups", UnsubscribeGroup{
		Name:        name,
		Description: description,
		IsDefault:   isDefault,
	})
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed creating API key: %w", err),
		}
	}

	if statusCode >= http.StatusMultipleChoices {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w, status: %d, response: %s", ErrFailedCreatingUnsubscribeGroup, statusCode, respBody),
		}
	}

	return parseUnsubscribeGroup(respBody)
}

// ReadUnsubscribeGroup retreives an UnsubscribeGroup and returns it.
func (c *Client) ReadUnsubscribeGroup(id string) (*UnsubscribeGroup, RequestError) {
	if id == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrUnsubscribeGroupIDRequired,
		}
	}

	respBody, _, err := c.Get("GET", "/asm/groups/"+id)
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return parseUnsubscribeGroup(respBody)
}

// ReadUnsubscribeGroups retrieves all UnsubscribeGroup and returns them.
func (c *Client) ReadUnsubscribeGroups() ([]UnsubscribeGroup, RequestError) {
	respBody, _, err := c.Get("GET", "/asm/groups")
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return parseUnsubscribeGroups(respBody)
}

// UpdateUnsubscribeGroup edits an UnsubscribeGroup and returns it.
func (c *Client) UpdateUnsubscribeGroup(
	id string,
	name string,
	description string,
	isDefault bool) (*UnsubscribeGroup, RequestError) {
	if id == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrUnsubscribeGroupIDRequired,
		}
	}

	t := UnsubscribeGroup{}
	t.Name = name
	t.IsDefault = isDefault

	if len(description) > 0 {
		t.Description = description
	}

	respBody, _, err := c.Post("PATCH", "/asm/groups/"+id, t)
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return parseUnsubscribeGroup(respBody)
}

// DeleteUnsubscribeGroup deletes an UnsubscribeGroup.
func (c *Client) DeleteUnsubscribeGroup(id string) (bool, RequestError) {
	if id == "" {
		return false, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrUnsubscribeGroupIDRequired,
		}
	}

	responseBody, statusCode, err := c.Get("DELETE", "/asm/groups/"+id)
	if err != nil {
		return false, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if statusCode >= http.StatusMultipleChoices && statusCode != http.StatusNotFound { // ignore not found
		return false, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w, status: %d, response: %s", ErrFailedDeletingUnsubscribeGroup, statusCode, responseBody),
		}
	}

	return true, RequestError{StatusCode: http.StatusOK, Err: nil}
}
