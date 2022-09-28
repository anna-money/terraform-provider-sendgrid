package sendgrid

import (
	"context"
	"fmt"
)

type User struct {
	Email   string   `json:"email,omitempty"`
	IsAdmin bool     `json:"disabled,omitempty"`
	Scopes  []string `json:"scopes,omitempty"`
}

//type SsoUser struct {
//	FirstName  string   `json:"first_name"`
//	LastName   string   `json:"last_name"`
//	Email      string   `json:"email"`
//	IsAdmin    string   `json:"is_admin"`
//	IsReadOnly string   `json:"is_read_only"`
//	Scopes     []string `json:"scopes"`
//}

func (c *Client) InviteTeammate(ctx context.Context, email string, scopes []string, isAdmin bool) (*User, error) {
	_, err := c.NewRequest("POST", "/teammates", User{
		Email:   email,
		IsAdmin: isAdmin,
		Scopes:  scopes,
	})
	if err != nil {
		return nil, err
	}

	return c.ReadTeammate(ctx, email)
}

func (c *Client) ReadTeammate(ctx context.Context, email string) (*User, error) {
	u := fmt.Sprintf("/teammates/%s/", email)
	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	user := new(User)
	_, err = c.Do(ctx, req, user)
	if err != nil {
		return user, err
	}

	return user, nil
}

//func (c *Client) UpdateUser(ctx context.Context, username string, disabled bool) (bool, RequestError) {
//	if username == "" {
//		return false, RequestError{StatusCode: http.StatusNotAcceptable, Err: ErrUsernameRequired}
//	}
//
//	respBody, statusCode, err := c.Post(ctx, "PATCH", "/subusers/"+username, SubUser{
//		Disabled: disabled,
//	})
//	if err != nil {
//		return false, RequestError{
//			StatusCode: statusCode,
//			Err:        fmt.Errorf("failed updating user: %w", err),
//		}
//	}
//
//	var body subUserErrors
//	if err = json.Unmarshal([]byte(respBody), &body); err != nil {
//		return false, RequestError{
//			StatusCode: http.StatusInternalServerError,
//			Err:        fmt.Errorf("failed updating user: %w", err),
//		}
//	}
//
//	return len(body.Errors) == 0, RequestError{StatusCode: http.StatusOK, Err: nil}
//}
