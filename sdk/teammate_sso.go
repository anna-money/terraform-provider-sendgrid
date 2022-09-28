package sendgrid

import (
	"context"
	"encoding/json"
	"fmt"
)

type SsoUser struct {
	FirstName  string   `json:"first_name"`
	LastName   string   `json:"last_name"`
	Email      string   `json:"email"`
	IsAdmin    bool     `json:"is_admin"`
	IsReadOnly bool     `json:"is_read_only"`
	Scopes     []string `json:"scopes"`
}

func parseUserSSO(respBody string) (*SsoUser, error) {
	var body SsoUser

	err := json.Unmarshal([]byte(respBody), &body)
	if err != nil {
		return nil, fmt.Errorf("failed parsing template: %w", err)
	}

	return &body, nil
}

func (c *Client) CreateUserSSO(ctx context.Context, first_name string, last_name string, email string, scopes []string, is_read_only bool, isAdmin bool) (*SsoUser, error) {
	respBody, _, err := c.Post(ctx, "POST", "/teammates", SsoUser{
		FirstName:  first_name,
		LastName:   last_name,
		Email:      email,
		IsReadOnly: is_read_only,
		IsAdmin:    isAdmin,
		Scopes:     scopes,
	})

	if err != nil {
		return nil, err
	}

	return parseUserSSO(respBody)
}

func (c *Client) ReadUserSSO(ctx context.Context, email string) (*SsoUser, error) {
	respBody, _, err := c.Get(ctx, "GET", "/teammates/"+email)
	if err != nil {
		return nil, err
	}

	return parseUserSSO(respBody)
}

func (c *Client) UpdateUserSSO(ctx context.Context, first_name string, last_name string, email string, scopes []string, is_read_only bool, isAdmin bool) (*SsoUser, error) {
	respBody, _, err := c.Post(ctx, "PATCH", "/teammates/"+email, SsoUser{
		FirstName:  first_name,
		LastName:   last_name,
		Email:      email,
		IsReadOnly: is_read_only,
		IsAdmin:    isAdmin,
		Scopes:     scopes,
	})

	if err != nil {
		return nil, err
	}

	return parseUserSSO(respBody)
}

func (c *Client) DeleteUserSSO(ctx context.Context, email string) (bool, error) {

	if _, statusCode, err := c.Get(ctx, "DELETE", "/teammates/"+email); statusCode > 299 || err != nil {
		return false, fmt.Errorf("failed deleting user: %w", err)
	}

	return true, nil
}
