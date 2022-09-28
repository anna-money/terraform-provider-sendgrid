package sendgrid

import (
	"context"
	"encoding/json"
	"fmt"
)

type User struct {
	Email   string   `json:"email,omitempty"`
	IsAdmin bool     `json:"disabled,omitempty"`
	Scopes  []string `json:"scopes,omitempty"`
}

func parseUser(respBody string) (*User, error) {
	var body User

	err := json.Unmarshal([]byte(respBody), &body)
	if err != nil {
		return nil, fmt.Errorf("failed parsing template: %w", err)
	}

	return &body, nil
}

func (c *Client) CreateUser(ctx context.Context, email string, scopes []string, isAdmin bool) (*User, error) {
	respBody, _, err := c.Post(ctx, "POST", "/teammates", User{
		Email:   email,
		IsAdmin: isAdmin,
		Scopes:  scopes,
	})

	if err != nil {
		return nil, err
	}

	return parseUser(respBody)
}

func (c *Client) ReadUser(ctx context.Context, email string) (*User, error) {
	respBody, _, err := c.Get(ctx, "GET", "/teammates/"+email)
	if err != nil {
		return nil, err
	}

	return parseUser(respBody)
}

func (c *Client) UpdateUser(ctx context.Context, email string, scopes []string, isAdmin bool) (*User, error) {
	respBody, _, err := c.Post(ctx, "PATCH", "/teammates/"+email, User{
		Email:   email,
		IsAdmin: isAdmin,
		Scopes:  scopes,
	})

	if err != nil {
		return nil, err
	}

	return parseUser(respBody)
}

func (c *Client) DeleteUser(ctx context.Context, email string) (bool, error) {

	if _, statusCode, err := c.Get(ctx, "DELETE", "/teammates/"+email); statusCode > 299 || err != nil {
		return false, fmt.Errorf("failed deleting user: %w", err)
	}

	return true, nil
}
