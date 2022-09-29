package sendgrid

import (
	"context"
	"encoding/json"
	"fmt"
)

// Template is a Sendgrid transactional template.
type Template struct {
	ID         string            `json:"id,omitempty"`
	Name       string            `json:"name,omitempty"`
	Generation string            `json:"generation,omitempty"`
	UpdatedAt  string            `json:"updated_at,omitempty"` //nolint:tagliatelle
	Versions   []TemplateVersion `json:"versions,omitempty"`
	Warnings   []string          `json:"warnings,omitempty"`
}

type Templates struct {
	Result []Template `json:"result"`
}

func parseTemplate(respBody string) (*Template, error) {
	var body Template

	err := json.Unmarshal([]byte(respBody), &body)
	if err != nil {
		return nil, fmt.Errorf("failed parsing template: %w", err)
	}

	return &body, nil
}

func parseTemplates(respBody string) ([]Template, error) {
	var body Templates

	err := json.Unmarshal([]byte(respBody), &body)
	if err != nil {
		return nil, fmt.Errorf("failed parsing template: %w", err)
	}

	return body.Result, nil
}

// CreateTemplate creates a transactional template and returns it.
func (c *Client) CreateTemplate(ctx context.Context, name, generation string) (*Template, error) {
	if name == "" {
		return nil, ErrTemplateNameRequired
	}

	if generation == "" {
		generation = "dynamic"
	}

	respBody, _, err := c.Post(ctx, "POST", "/templates", Template{
		Name:       name,
		Generation: generation,
	})
	fmt.Printf("ERROR from tmpl: %v", err)
	if err != nil {
		return nil, fmt.Errorf("failed creating template: %w", err)
	}

	return parseTemplate(respBody)
}

// ReadTemplate retreives a transactional template and returns it.
func (c *Client) ReadTemplate(ctx context.Context, id string) (*Template, error) {
	if id == "" {
		return nil, ErrTemplateIDRequired
	}

	respBody, _, err := c.Get(ctx, "GET", "/templates/"+id)
	if err != nil {
		return nil, fmt.Errorf("failed reading template: %w", err)
	}

	return parseTemplate(respBody)
}

func (c *Client) ReadTemplates(ctx context.Context, generation string) ([]Template, error) {
	respBody, _, err := c.Get(ctx, "GET", "/templates?page_size=200&generations="+generation)
	if err != nil {
		return nil, fmt.Errorf("failed reading template: %w", err)
	}

	return parseTemplates(respBody)
}

// UpdateTemplate edits a transactional template and returns it.
// We can't change the "generation" of a transactional template.
func (c *Client) UpdateTemplate(ctx context.Context, id, name string) (*Template, error) {
	if id == "" {
		return nil, ErrTemplateIDRequired
	}

	if name == "" {
		return nil, ErrTemplateNameRequired
	}

	respBody, _, err := c.Post(ctx, "PATCH", "/templates/"+id, Template{
		Name: name,
	},
	)
	if err != nil {
		return nil, fmt.Errorf("failed updating template: %w", err)
	}

	return parseTemplate(respBody)
}

// DeleteTemplate deletes a transactional template.
func (c *Client) DeleteTemplate(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, ErrTemplateIDRequired
	}

	if _, statusCode, err := c.Get(ctx, "DELETE", "/templates/"+id); statusCode > 299 || err != nil {
		return false, fmt.Errorf("failed deleting template: %w", err)
	}

	return true, nil
}
