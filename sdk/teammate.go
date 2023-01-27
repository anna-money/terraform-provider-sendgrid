package sendgrid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

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

type Users struct {
	Result []User `json:"result"`
}

type PendingUser struct {
	Result []struct {
		Token          string   `json:"token,omitempty"`
		Email          string   `json:"email,omitempty"`
		IsAdmin        bool     `json:"is_admin,omitempty"`
		IsReadOnly     bool     `json:"is_read_only,omitempty"`
		ExpirationDate int      `json:"expiration_date,omitempty"`
		Scopes         []string `json:"scopes,omitempty"`
	} `json:"result"`
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
	respBody, _, err := c.Get(ctx, "GET", "/teammates?limit=10000")
	if err != nil {
		return "", err
	}

	users := &Users{}

	decoder := json.NewDecoder(bytes.NewReader([]byte(respBody)))
	err = decoder.Decode(users)
	if err != nil {
		return "", err
	}

	for _, user := range users.Result {
		if user.Email == email && user.Username != "" {
			return user.Username, nil
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

func (c *Client) CreateSSOUser(ctx context.Context, firstName, lastName, email string, scopes []string, isAdmin bool) (*User, error) {
	respBody, _, err := c.Post(ctx, "POST", "/sso/teammates", User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		IsAdmin:   isAdmin,
		Scopes:    scopes,
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
		IsAdmin: isAdmin,
		Scopes:  scopes,
	})
	if err != nil {
		return nil, err
	}

	return parseUser(respBody)
}

func (c *Client) UpdateSSOUser(ctx context.Context, firstName, lastName, email string, scopes []string, isAdmin bool) (*User, error) {
	username, err := c.GetUsernameByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	respBody, _, err := c.Post(ctx, "PATCH", "/sso/teammates/"+username, User{
		FirstName: firstName,
		LastName:  lastName,
		IsAdmin:   isAdmin,
		Scopes:    scopes,
	})
	if err != nil {
		return nil, err
	}

	return parseUser(respBody)
}

func (c *Client) DeleteUser(ctx context.Context, email string) (bool, error) {
	username, err := c.GetUsernameByEmail(ctx, email)
	if err != nil {
		tokenInvite, err := c.GetPendingUserToken(ctx, email)
		if _, statusCode, err := c.Get(ctx, "DELETE", "/teammates/pending/"+tokenInvite); statusCode > 299 || err != nil {
			return false, fmt.Errorf("failed deleting user: %w", err)
		}
		return false, err
	}

	if _, statusCode, err := c.Get(ctx, "DELETE", "/teammates/"+username); statusCode > 299 || err != nil {
		return false, fmt.Errorf("failed deleting user: %w", err)
	}

	return true, nil
}

func (c *Client) GetPendingUserToken(ctx context.Context, email string) (string, error) {
	respBody, _, err := c.Get(ctx, "GET", "/teammates/pending")
	if err != nil {
		return "", err
	}

	pendingUsers := &PendingUser{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(respBody)))
	err = decoder.Decode(pendingUsers)
	if err != nil {
		return "", err
	}

	for _, user := range pendingUsers.Result {
		if user.Email == email {
			return user.Token, nil
		}
	}
	return "", fmt.Errorf("pending user with email %s not found", email)
}
