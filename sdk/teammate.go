package sendgrid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
		PendingID      string   `json:"pending_id,omitempty"`
		Token          string   `json:"token,omitempty"`
		Email          string   `json:"email,omitempty"`
		IsAdmin        bool     `json:"is_admin,omitempty"`
		IsReadOnly     bool     `json:"is_read_only,omitempty"`
		ExpirationDate int      `json:"expiration_date,omitempty"`
		Scopes         []string `json:"scopes,omitempty"`
	} `json:"result"`
}

func parseUser(respBody string) (*User, RequestError) {
	var body User

	err := json.Unmarshal([]byte(respBody), &body)
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed parsing teammate: %w", err),
		}
	}

	return &body, RequestError{StatusCode: http.StatusOK, Err: nil}
}

func (c *Client) GetUsernameByEmail(ctx context.Context, email string) (string, RequestError) {
	respBody, statusCode, err := c.Get(ctx, "GET", "/teammates?limit=10000")
	if err != nil {
		return "", RequestError{
			StatusCode: statusCode,
			Err:        err,
		}
	}

	users := &Users{}

	decoder := json.NewDecoder(bytes.NewReader([]byte(respBody)))
	err = decoder.Decode(users)
	if err != nil {
		return "", RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	for _, user := range users.Result {
		if user.Email == email && user.Username != "" {
			return user.Username, RequestError{StatusCode: http.StatusOK, Err: nil}
		}
	}
	return "", RequestError{
		StatusCode: http.StatusNotFound,
		Err:        fmt.Errorf("username with email %s not found", email),
	}
}

func (c *Client) CreateUser(ctx context.Context, email string, scopes []string, isAdmin bool) (*User, RequestError) {
	respBody, statusCode, err := c.Post(ctx, "POST", "/teammates", User{
		Email:   email,
		IsAdmin: isAdmin,
		Scopes:  scopes,
	})
	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        err,
		}
	}

	return parseUser(respBody)
}

func (c *Client) CreateSSOUser(ctx context.Context, firstName, lastName, email string, scopes []string, isAdmin bool) (*User, RequestError) {
	respBody, statusCode, err := c.Post(ctx, "POST", "/sso/teammates", User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		IsAdmin:   isAdmin,
		Scopes:    scopes,
	})
	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        err,
		}
	}

	return parseUser(respBody)
}

func (c *Client) ReadUser(ctx context.Context, email string) (*User, RequestError) {
	username, requestErr := c.GetUsernameByEmail(ctx, email)
	if requestErr.Err != nil {
		// If user not found in active teammates, check pending invitations
		if requestErr.StatusCode == http.StatusNotFound {
			pendingUser, pendingErr := c.ReadPendingUser(ctx, email)
			if pendingErr.Err != nil {
				// User not found in either active or pending
				return nil, RequestError{
					StatusCode: http.StatusNotFound,
					Err:        fmt.Errorf("user with email %s not found in active teammates or pending invitations. Original active error: %v. Pending error: %v", email, requestErr.Err, pendingErr.Err),
				}
			}
			return pendingUser, RequestError{StatusCode: http.StatusOK, Err: nil}
		}
		return nil, requestErr
	}

	respBody, statusCode, err := c.Get(ctx, "GET", "/teammates/"+username)
	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        err,
		}
	}

	var u User
	err = json.Unmarshal([]byte(respBody), &u)
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return &u, RequestError{StatusCode: http.StatusOK, Err: nil}
}

func (c *Client) UpdateUser(ctx context.Context, email string, scopes []string, isAdmin bool) (*User, RequestError) {
	username, requestErr := c.GetUsernameByEmail(ctx, email)
	if requestErr.Err != nil {
		return nil, requestErr
	}

	respBody, statusCode, err := c.Post(ctx, "PATCH", "/teammates/"+username, User{
		IsAdmin: isAdmin,
		Scopes:  scopes,
	})
	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        err,
		}
	}

	return parseUser(respBody)
}

func (c *Client) UpdateSSOUser(ctx context.Context, firstName, lastName, email string, scopes []string, isAdmin bool) (*User, RequestError) {
	username, requestErr := c.GetUsernameByEmail(ctx, email)
	if requestErr.Err != nil {
		// If user not found in active teammates, they might be pending
		// Pending users cannot be updated, so return an error
		if requestErr.StatusCode == http.StatusNotFound {
			return nil, RequestError{
				StatusCode: http.StatusNotFound,
				Err:        fmt.Errorf("user %s not found in active teammates - they may be pending and cannot be updated", email),
			}
		}
		return nil, requestErr
	}

	respBody, statusCode, err := c.Post(ctx, "PATCH", "/sso/teammates/"+username, User{
		FirstName: firstName,
		LastName:  lastName,
		IsAdmin:   isAdmin,
		Scopes:    scopes,
	})
	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        err,
		}
	}

	return parseUser(respBody)
}

func (c *Client) DeleteUser(ctx context.Context, email string) (bool, RequestError) {
	username, requestErr := c.GetUsernameByEmail(ctx, email)
	if requestErr.Err != nil {
		tokenInvite, tokenErr := c.GetPendingUserToken(ctx, email)
		if tokenErr.Err != nil {
			return false, tokenErr
		}

		if _, statusCode, err := c.Get(ctx, "DELETE", "/teammates/pending/"+tokenInvite); statusCode > 299 || err != nil {
			return false, RequestError{
				StatusCode: statusCode,
				Err:        fmt.Errorf("failed deleting user: %w", err),
			}
		}
		return false, RequestError{StatusCode: http.StatusOK, Err: nil}
	}

	if _, statusCode, err := c.Get(ctx, "DELETE", "/teammates/"+username); statusCode > 299 || err != nil {
		return false, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed deleting user: %w", err),
		}
	}

	return true, RequestError{StatusCode: http.StatusOK, Err: nil}
}

func (c *Client) GetPendingUserToken(ctx context.Context, email string) (string, RequestError) {
	respBody, statusCode, err := c.Get(ctx, "GET", "/teammates/pending?limit=200")
	if err != nil {
		return "", RequestError{
			StatusCode: statusCode,
			Err:        err,
		}
	}

	pendingUsers := &PendingUser{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(respBody)))
	err = decoder.Decode(pendingUsers)
	if err != nil {
		return "", RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	for _, user := range pendingUsers.Result {
		if user.Email == email {
			// SendGrid API returns token field, not pending_id
			if user.Token != "" {
				return user.Token, RequestError{StatusCode: http.StatusOK, Err: nil}
			}
			// Fallback to pending_id if token is empty (though this seems unlikely based on API response)
			if user.PendingID != "" {
				return user.PendingID, RequestError{StatusCode: http.StatusOK, Err: nil}
			}
		}
	}
	return "", RequestError{
		StatusCode: http.StatusNotFound,
		Err:        fmt.Errorf("pending user with email %s not found", email),
	}
}

// ReadPendingUser reads a pending user invitation by email
func (c *Client) ReadPendingUser(ctx context.Context, email string) (*User, RequestError) {
	respBody, statusCode, err := c.Get(ctx, "GET", "/teammates/pending?limit=10000")
	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed to get pending users: %w", err),
		}
	}

	pendingUsers := &PendingUser{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(respBody)))
	err = decoder.Decode(pendingUsers)
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed to decode pending users response: %w", err),
		}
	}

	// Debug: log all pending users with more details
	var pendingDetails []string
	for _, pendingUser := range pendingUsers.Result {
		detail := fmt.Sprintf("email=%s, pending_id=%s, token=%s, expiration=%d",
			pendingUser.Email, pendingUser.PendingID, pendingUser.Token, pendingUser.ExpirationDate)
		pendingDetails = append(pendingDetails, detail)

		if pendingUser.Email == email {
			// Convert pending user to User struct
			user := &User{
				Email:   pendingUser.Email,
				IsAdmin: pendingUser.IsAdmin,
				Scopes:  pendingUser.Scopes,
				// Mark as pending by setting a special user type
				UserType: "pending",
			}
			return user, RequestError{StatusCode: http.StatusOK, Err: nil}
		}
	}

	return nil, RequestError{
		StatusCode: http.StatusNotFound,
		Err:        fmt.Errorf("pending user with email %s not found. Available pending users: %v. This may mean the user has already accepted the invitation or the invitation has expired", email, pendingDetails),
	}
}
