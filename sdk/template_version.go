package sendgrid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// TemplateVersion is a Sendgrid transactional template version.
type TemplateVersion struct {
	ID                   string    `json:"id,omitempty"`
	TemplateID           string    `json:"template_id,omitempty"`   //nolint:tagliatelle
	UpdatedAt            string    `json:"updated_at,omitempty"`    //nolint:tagliatelle
	ThumbnailURL         string    `json:"thumbnail_url,omitempty"` //nolint:tagliatelle
	Warnings             []Warning `json:"warning,omitempty"`
	Active               int       `json:"active,omitempty"`
	Name                 string    `json:"name,omitempty"`
	HTMLContent          string    `json:"html_content,omitempty"`           //nolint:tagliatelle
	PlainContent         string    `json:"plain_content,omitempty"`          //nolint:tagliatelle
	GeneratePlainContent bool      `json:"generate_plain_content,omitempty"` //nolint:tagliatelle
	Subject              string    `json:"subject,omitempty"`
	Editor               string    `json:"editor,omitempty"`
	TestData             string    `json:"test_data,omitempty"` //nolint:tagliatelle
}

type Warning struct {
	Message string `json:"message,omitempty"`
}

func parseTemplateVersion(respBody string) (*TemplateVersion, RequestError) {
	var body TemplateVersion

	err := json.Unmarshal([]byte(respBody), &body)
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed parsing template version: %w", err),
		}
	}

	return &body, RequestError{StatusCode: http.StatusOK, Err: nil}
}

// CreateTemplateVersion creates a new version of a transactional template and returns it.
func (c *Client) CreateTemplateVersion(ctx context.Context, t TemplateVersion) (*TemplateVersion, RequestError) {
	if t.TemplateID == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrTemplateVersionIDRequired,
		}
	}

	if t.Name == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrTemplateVersionNameRequired,
		}
	}

	if t.Subject == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrTemplateVersionSubjectRequired,
		}
	}

	respBody, statusCode, err := c.Post(ctx, "POST", "/templates/"+t.TemplateID+"/versions", t)
	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed creating template version: %w", err),
		}
	}

	return parseTemplateVersion(respBody)
}

// ReadTemplateVersion retreives a version of a transactional template and returns it.
func (c *Client) ReadTemplateVersion(ctx context.Context, templateID, id string) (*TemplateVersion, RequestError) {
	if templateID == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrTemplateVersionIDRequired,
		}
	}

	if id == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrTemplateIDRequired,
		}
	}

	respBody, statusCode, err := c.Get(ctx, "GET", "/templates/"+templateID+"/versions/"+id)
	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed reading template version: %w", err),
		}
	}

	return parseTemplateVersion(respBody)
}

// UpdateTemplateVersion edits a version of a transactional template and returns it.
func (c *Client) UpdateTemplateVersion(ctx context.Context, t TemplateVersion) (*TemplateVersion, RequestError) {
	if t.ID == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrTemplateVersionIDRequired,
		}
	}

	if t.TemplateID == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrTemplateIDRequired,
		}
	}

	respBody, statusCode, err := c.Post(ctx, "PATCH", "/templates/"+t.TemplateID+"/versions/"+t.ID, t)
	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed updating template version: %w", err),
		}
	}

	return parseTemplateVersion(respBody)
}

// ActivateTemplateVersion activates a version of a transactional template and returns it.
func (c *Client) ActivateTemplateVersion(ctx context.Context, t TemplateVersion) (*TemplateVersion, RequestError) {
	if t.ID == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrTemplateVersionIDRequired,
		}
	}

	if t.TemplateID == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrTemplateIDRequired,
		}
	}

	respBody, statusCode, err := c.Post(ctx, "POST", "/templates/"+t.TemplateID+"/versions/"+t.ID+"/activate", t)
	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed activating template version: %w", err),
		}
	}

	return parseTemplateVersion(respBody)
}

// DeleteTemplateVersion deletes a version of a transactional template.
func (c *Client) DeleteTemplateVersion(ctx context.Context, templateID, id string) (bool, RequestError) {
	if templateID == "" {
		return false, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrTemplateVersionIDRequired,
		}
	}

	if _, statusCode, err := c.Get(ctx, "DELETE", "/templates/"+templateID+"/versions/"+id); statusCode > 299 ||
		err != nil {
		return false, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed deleting template version: %w", err),
		}
	}

	return true, RequestError{StatusCode: http.StatusOK, Err: nil}
}
