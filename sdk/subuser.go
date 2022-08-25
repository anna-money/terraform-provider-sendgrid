package sendgrid

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type creditAllocation struct {
	Type string `json:"credit_allocation,omitempty"` //nolint:tagliatelle
}

// SubUser is a Sendgrid SubUser.
type SubUser struct {
	ID                 int              `json:"id,omitempty"`
	UserID             int              `json:"user_id,omitempty"` //nolint:tagliatelle
	UserName           string           `json:"username,omitempty"`
	Password           string           `json:"password,omitempty"`
	ConfirmPassword    string           `json:"confirm_password,omitempty"` //nolint:tagliatelle
	Email              string           `json:"email,omitempty"`
	IPs                []string         `json:"ips,omitempty"`
	Disabled           bool             `json:"disabled,omitempty"`
	SignupSessionToken string           `json:"signup_session_token,omitempty"` //nolint:tagliatelle
	AuthorizationToken string           `json:"authorization_token,omitempty"`  //nolint:tagliatelle
	CreditAllocation   creditAllocation `json:"credit_allocation,omitempty"`    //nolint:tagliatelle
}

type UpdateSubUserPassword struct {
	NewPassword string `json:"new_password"` //nolint:tagliatelle
	OldPassword string `json:"old_password"` //nolint:tagliatelle
}

func parseSubUser(respBody string) (*SubUser, RequestError) {
	var body SubUser
	if err := json.Unmarshal([]byte(respBody), &body); err != nil {
		log.Printf("[DEBUG] [parseSubUser] failed parsing subUser, response body: %s", respBody)

		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return &body, RequestError{StatusCode: http.StatusOK, Err: nil}
}

func parseSubUsers(respBody string) ([]SubUser, RequestError) {
	var body []SubUser
	if err := json.Unmarshal([]byte(respBody), &body); err != nil {
		log.Printf("[DEBUG] [parseSubUsers] failed parsing subUsers, response body: %s", respBody)

		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return body, RequestError{StatusCode: http.StatusOK, Err: nil}
}

// CreateSubuser creates a subuser and returns it.
func (c *Client) CreateSubuser(username, email, password string, ips []string) (*SubUser, RequestError) {
	if username == "" {
		return nil, RequestError{StatusCode: http.StatusNotAcceptable, Err: ErrUsernameRequired}
	}

	if email == "" {
		return nil, RequestError{StatusCode: http.StatusNotAcceptable, Err: ErrEmailRequired}
	}

	if password == "" {
		return nil, RequestError{StatusCode: http.StatusNotAcceptable, Err: ErrPasswordRequired}
	}

	if len(ips) < 1 {
		return nil, RequestError{StatusCode: http.StatusNotAcceptable, Err: ErrIPRequired}
	}

	respBody, statusCode, err := c.Post("POST", "/subusers", SubUser{
		UserName:        username,
		Email:           email,
		Password:        password,
		ConfirmPassword: password,
		IPs:             ips,
	})
	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed creating subUser: %w", err),
		}
	}

	if statusCode >= http.StatusMultipleChoices {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w, status: %d, response: %s", ErrFailedCreatingSubUser, statusCode, respBody),
		}
	}

	return parseSubUser(respBody)
}

// ReadSubUser retreives a subuser and returns it.
func (c *Client) ReadSubUser(username string) ([]SubUser, RequestError) {
	if username == "" {
		return nil, RequestError{StatusCode: http.StatusNotAcceptable, Err: ErrUsernameRequired}
	}

	endpoint := "/subusers?username=" + url.QueryEscape(username)

	respBody, statusCode, err := c.Get("GET", endpoint)
	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed reading subUser: %w", err),
		}
	}

	return parseSubUsers(respBody)
}

// UpdateSubuser enables/disables a subuser.
func (c *Client) UpdateSubuser(username string, disabled bool) (bool, RequestError) {
	if username == "" {
		return false, RequestError{StatusCode: http.StatusNotAcceptable, Err: ErrUsernameRequired}
	}

	respBody, statusCode, err := c.Post("PATCH", "/subusers/"+username, SubUser{
		Disabled: disabled,
	})
	if err != nil {
		return false, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed updating subUser: %w", err),
		}
	}

	var body subUserErrors
	if err = json.Unmarshal([]byte(respBody), &body); err != nil {
		return false, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed updating subUser: %w", err),
		}
	}

	return len(body.Errors) == 0, RequestError{StatusCode: http.StatusOK, Err: nil}
}

func (c *Client) UpdateSubuserIPs(username string, ips []string) RequestError {
	if username == "" {
		return RequestError{StatusCode: http.StatusNotAcceptable, Err: ErrUsernameRequired}
	}

	_, statusCode, err := c.Post("PUT", "/subusers/"+username+"/ips", ips)
	if err != nil {
		return RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed updating subUser Ips: %w", err),
		}
	}

	return RequestError{StatusCode: http.StatusOK, Err: nil}
}

// DeleteSubuser deletes a subuser.
func (c *Client) DeleteSubuser(username string) (bool, RequestError) {
	if username == "" {
		return false, RequestError{StatusCode: http.StatusNotAcceptable, Err: ErrUsernameRequired}
	}

	respBody, statusCode, err := c.Get("DELETE", "/subusers/"+username)
	if err != nil {
		return false, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed deleting subUser: %w", err),
		}
	}

	if statusCode >= http.StatusMultipleChoices && statusCode != http.StatusNotFound { // ignore not found
		return false, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w: statusCode: %d, respBody: %s", err, statusCode, respBody),
		}
	}

	return true, RequestError{StatusCode: http.StatusOK, Err: nil}
}

func (c *Client) UpdateSubuserPassword(username string, oldPassword string, newPassword string) RequestError {
	if newPassword == "" {
		return RequestError{StatusCode: http.StatusBadRequest, Err: ErrSubUserPassword}
	}

	origOnBehalfOf := c.OnBehalfOf
	c.OnBehalfOf = username
	_, statusCode, err := c.Post("PUT", "/user/password", UpdateSubUserPassword{
		NewPassword: newPassword,
		OldPassword: oldPassword,
	})
	c.OnBehalfOf = origOnBehalfOf

	if err != nil {
		return RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed updating subUser password: %w", err),
		}
	}

	if statusCode >= http.StatusMultipleChoices {
		return RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w: statusCode: %d", err, statusCode),
		}
	}

	return RequestError{StatusCode: http.StatusOK, Err: nil}
}
