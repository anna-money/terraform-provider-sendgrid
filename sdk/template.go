package sendgrid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Template is a Sendgrid transactional template.
type Template struct {
	ID         string            `json:"id,omitempty"`
	Name       string            `json:"name,omitempty"`
	Generation string            `json:"generation,omitempty"`
	UpdatedAt  string            `json:"updated_at,omitempty"` //nolint:tagliatelle
	Versions   []TemplateVersion `json:"versions,omitempty"`
	Warning    Warning           `json:"warning,omitempty"`
}

type Templates struct {
	Result []Template `json:"result"`
}

func parseTemplate(respBody string) (*Template, RequestError) {
	var body Template

	err := json.Unmarshal([]byte(respBody), &body)
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed parsing template: %w", err),
		}
	}

	return &body, RequestError{StatusCode: http.StatusOK, Err: nil}
}

func parseTemplates(respBody string) ([]Template, RequestError) {
	var body Templates

	err := json.Unmarshal([]byte(respBody), &body)
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed parsing template: %w", err),
		}
	}

	return body.Result, RequestError{StatusCode: http.StatusOK, Err: nil}
}

// CreateTemplate creates a transactional template and returns it.
func (c *Client) CreateTemplate(ctx context.Context, name, generation string) (*Template, RequestError) {
	if name == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrTemplateNameRequired,
		}
	}

	if generation == "" {
		generation = "dynamic"
	}

	respBody, statusCode, err := c.Post(ctx, "POST", "/templates", Template{
		Name:       name,
		Generation: generation,
	})
	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed creating template: %w", err),
		}
	}

	return parseTemplate(respBody)
}

// ReadTemplate retreives a transactional template and returns it.
func (c *Client) ReadTemplate(ctx context.Context, id string) (*Template, RequestError) {
	if id == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrTemplateIDRequired,
		}
	}

	respBody, statusCode, err := c.Get(ctx, "GET", "/templates/"+id)
	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed reading template: %w", err),
		}
	}

	return parseTemplate(respBody)
}

func (c *Client) ReadTemplates(ctx context.Context, generation string) ([]Template, RequestError) {
	respBody, statusCode, err := c.Get(ctx, "GET", "/templates?page_size=200&generations="+generation)
	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed reading template: %w", err),
		}
	}

	return parseTemplates(respBody)
}

// UpdateTemplate edits a transactional template and returns it.
// We can't change the "generation" of a transactional template.
func (c *Client) UpdateTemplate(ctx context.Context, id, name string) (*Template, RequestError) {
	if id == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrTemplateIDRequired,
		}
	}

	if name == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrTemplateNameRequired,
		}
	}

	respBody, statusCode, err := c.Post(ctx, "PATCH", "/templates/"+id, Template{
		Name: name,
	})
	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed updating template: %w", err),
		}
	}

	return parseTemplate(respBody)
}

// DeleteTemplate deletes a transactional template.
func (c *Client) DeleteTemplate(ctx context.Context, id string) (bool, RequestError) {
	if id == "" {
		return false, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrTemplateIDRequired,
		}
	}

	if _, statusCode, err := c.Get(ctx, "DELETE", "/templates/"+id); statusCode > 299 || err != nil {
		return false, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed deleting template: %w", err),
		}
	}

	return true, RequestError{StatusCode: http.StatusOK, Err: nil}
}
