package sendgrid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

var (
	UserName string
)

type Users struct {
	Result []struct {
		Email     string   `json:"email,omitempty"`
		Username  string   `json:"username,omitempty"`
		FirstName string   `json:"first_name,omitempty"`
		LastName  string   `json:"last_name,omitempty"`
		IsAdmin   bool     `json:"is_admin,omitempty"`
		IsSso     bool     `json:"is_sso,omitempty"`
		UserType  string   `json:"user_type,omitempty"`
		Scopes    []string `json:"scopes,omitempty"`
	} `json:"result"`
}

type User struct {
	Username  string   `json:"username,omitempty"`
	Email     string   `json:"email,omitempty"`
	FirstName string   `json:"first_name,omitempty"`
	LastName  string   `json:"last_name,omitempty"`
	Address   string   `json:"address,omitempty"`
	Address2  string   `json:"address2,omitempty"`
	City      string   `json:"city,omitempty"`
	State     string   `json:"state,omitempty"`
	Zip       string   `json:"zip,omitempty"`
	Country   string   `json:"country,omitempty"`
	Company   string   `json:"company,omitempty"`
	Website   string   `json:"website,omitempty"`
	Phone     string   `json:"phone,omitempty"`
	IsAdmin   bool     `json:"is_admin,omitempty"`
	IsSSO     bool     `json:"is_sso,omitempty"`
	UserType  string   `json:"user_type,omitempty"`
	Scopes    []string `json:"scopes,omitempty"`
}

func parseUser(respBody string) (*User, error) {
	var body User

	err := json.Unmarshal([]byte(respBody), &body)
	if err != nil {
		return nil, fmt.Errorf("failed parsing teammate: %w", err)
	}

	return &body, nil
}

func (c *Client) GetUsernameByEmail(ctx context.Context, email string) (string, error) {
	respBody, _, err := c.Get(ctx, "GET", "/teammates")
	if err != nil {
		return "", err
	}

	users := &Users{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(respBody)))
	err = decoder.Decode(users)
	if err != nil {
		return "", err
	}

	for _, u := range users.Result {
		if u.Email == email && u.Username != "" {
			return u.Username, nil
		}
	}
	return "", fmt.Errorf("username with email %s not found", email)
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
	username, err := c.GetUsernameByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	respBody, _, err := c.Get(ctx, "GET", "/teammates/"+username)
	if err != nil {
		return nil, err
	}

	var u User
	err = json.Unmarshal([]byte(respBody), &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (c *Client) UpdateUser(ctx context.Context, email string, scopes []string, isAdmin bool) (*User, error) {
	username, err := c.GetUsernameByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	respBody, _, err := c.Post(ctx, "PATCH", "/teammates/"+username, User{
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
	username, err := c.GetUsernameByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if _, statusCode, err := c.Get(ctx, "DELETE", "/teammates/"+username); statusCode > 299 || err != nil {
		return false, fmt.Errorf("failed deleting user: %w", err)
	}

	return true, nil
}
