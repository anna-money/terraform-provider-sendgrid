package sendgrid

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	// ErrBodyNotNil low error displayed when the prepared body for a POST call
	// to the API is nil.
	ErrBodyNotNil = errors.New("body must not be nil")

	// ErrNameRequired error displayed when an api key name wasn't specified.
	ErrNameRequired = errors.New("a name is required")

	// ErrAPIKeyIDRequired error displayed when an api key ID wasn't specified.
	ErrAPIKeyIDRequired = errors.New("an API Key ID is required")

	// ErrFailedCreatingAPIKey error displayed when the provider can not create an api key.
	ErrFailedCreatingAPIKey = errors.New("failed creating apiKey")

	// ErrFailedDeletingAPIKey error displayed when the provider can not delete an api key.
	ErrFailedDeletingAPIKey = errors.New("failed deleting apiKey")

	// ErrUsernameRequired error displayed when a subUser username wasn't specified.
	ErrUsernameRequired = errors.New("a username is required")

	// ErrEmailRequired error displayed when a subUser email wasn't specified.
	ErrEmailRequired = errors.New("an email is required")

	// ErrPasswordRequired error displayed when a subUser password wasn't specified.
	ErrPasswordRequired = errors.New("a password is required")

	// ErrIPRequired error displayed when at least one IP per subUser wasn't specified.
	ErrIPRequired = errors.New("at least one ip address is required")

	// ErrFailedCreatingSubUser error displayed when the provider can not create a subuser.
	ErrFailedCreatingSubUser = errors.New("failed creating subUser")

	// ErrFailedDeletingSubUser error displayed when the provider can not delete a subuser.
	ErrFailedDeletingSubUser = errors.New("failed deleting subUser")

	// ErrTemplateIDRequired error displayed when a template ID wasn't specified.
	ErrTemplateIDRequired = errors.New("a template ID is required")

	// ErrTemplateNameRequired error displayed when a template name wasn't specified.
	ErrTemplateNameRequired = errors.New("a template name is required")

	// ErrTemplateVersionIDRequired error displayed when a template version ID wasn't specified.
	ErrTemplateVersionIDRequired = errors.New("a template version ID is required")

	// ErrTemplateVersionNameRequired error displayed when a template version ID wasn't specified.
	ErrTemplateVersionNameRequired = errors.New("a template version name is required")

	// ErrTemplateVersionSubjectRequired error displayed when a template version subject wasn't specified.
	ErrTemplateVersionSubjectRequired = errors.New("a template version subject is required")

	ErrFailedCreatingUnsubscribeGroup = errors.New("failed to create unsubscribe list")

	ErrUnsubscribeGroupIDRequired = errors.New("unsubscribe list id is required")

	ErrFailedDeletingUnsubscribeGroup = errors.New("failed deleting unsubscribe list")

	ErrFailedCreatingParseWebhook = errors.New("failed to create parse webhook")

	ErrFailedDeletingParseWebhook = errors.New("failed deleting parse webhook")

	ErrHostnameRequired = errors.New("a hostname is required")

	ErrURLRequired = errors.New("a url is required")

	ErrFailedPatchingEventWebhook = errors.New("failed to patch event webhook")

	ErrFailedCreatingDomainAuthentication = errors.New("failed to create domain authentication")

	ErrDomainAuthenticationIDRequired = errors.New("id for domain authentication is required")

	ErrFailedDeletingDomainAuthentication = errors.New("failed deleting domain authentication")

	ErrLinkBrandingIDRequired = errors.New("link branding id is required")

	ErrFailedDeletingLinkBranding = errors.New("failed to delete link branding")

	ErrFailedCreatingLinkBranding = errors.New("failed to create link branding")

	// ErrSubUserPassword should be empty.
	ErrSubUserPassword = errors.New("new password must be non empty")

	// ErrSSOIntegrationMissingField error displayed when a required SSO integration field is not specified.
	ErrSSOIntegrationMissingField = errors.New("SSO integration field is missing")

	// ErrFailedCreatingSSOIntegration error displayed when an SSO integration creation request fails.
	ErrFailedCreatingSSOIntegration = errors.New("failed to create SSO integration")

	// ErrFailedUpdatingSSOIntegration error displayed when an SSO integration update request fails.
	ErrFailedUpdatingSSOIntegration = errors.New("failed to update SSO integration")

	// ErrSSOCertificateMissingField error displayed when a required SSO certificate field is not specified.
	ErrSSOCertificateMissingField = errors.New("SSO certificate field is missing")

	// ErrFailedCreatingSSOCertificate error displayed when an SSO certificate creation request fails.
	ErrFailedCreatingSSOCertificate = errors.New("failed to create SSO certificate")

	// ErrFailedUpdatingSSOCertificate error displayed when an SSO certificate update request fails.
	ErrFailedUpdatingSSOCertificate = errors.New("failed to update SSO certificate")
)

type APIError struct {
	f interface{} // unknown
}

// RequestError struct permits to embed to return the statucode and the error to the parent function.
type RequestError struct {
	StatusCode int
	Err        error
}

type subUserError struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

type subUserErrors struct {
	Errors []subUserError `json:"errors,omitempty"`
}

// RetryOnRateLimit management of RequestErrors, and launch a retry if needed.
func RetryOnRateLimit(
	ctx context.Context, d *schema.ResourceData, f func() (interface{}, RequestError),
) (interface{}, error) {
	var resp interface{}

	err := resource.RetryContext(
		ctx,
		d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			var requestErr RequestError
			resp, requestErr = f()
			if requestErr.Err != nil {
				if requestErr.StatusCode == http.StatusTooManyRequests {
					return resource.RetryableError(requestErr.Err)
				}

				return resource.NonRetryableError(requestErr.Err)
			}

			return nil
		})
	if err != nil {
		return resp, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

func (e *APIError) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &e.f); err != nil {
		e.f = string(b)
	}
	return nil
}

func (e *APIError) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.f)
}

func (e APIError) Detail() string {
	switch v := e.f.(type) {
	case map[string]interface{}:
		if len(v) == 1 {
			if detail, ok := v["detail"].(string); ok {
				return detail
			}
		}
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (e APIError) Error() string {
	return fmt.Sprintf("sendgrid: %s", e.Detail())
}

// Empty returns true if empty.
func (e APIError) Empty() bool {
	return e.f == nil
}
