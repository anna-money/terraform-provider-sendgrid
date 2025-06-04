package sendgrid

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestValidateTeammateScopes(t *testing.T) {
	tests := []struct {
		name          string
		scopes        []interface{}
		expectErrors  bool
		errorContains string
	}{
		{
			name:         "valid scopes",
			scopes:       []interface{}{"mail.send", "templates.read", "stats.read"},
			expectErrors: false,
		},
		{
			name:          "invalid scope",
			scopes:        []interface{}{"mail.send", "invalid.scope"},
			expectErrors:  true,
			errorContains: "not valid or assignable",
		},
		{
			name:          "automatic scope 2fa_exempt",
			scopes:        []interface{}{"mail.send", "2fa_exempt"},
			expectErrors:  true,
			errorContains: "set automatically by SendGrid",
		},
		{
			name:          "automatic scope 2fa_required",
			scopes:        []interface{}{"mail.send", "2fa_required"},
			expectErrors:  true,
			errorContains: "set automatically by SendGrid",
		},
		{
			name:          "mix of valid and invalid",
			scopes:        []interface{}{"mail.send", "invalid.scope", "templates.read"},
			expectErrors:  true,
			errorContains: "not valid or assignable",
		},
		{
			name:         "marketing scopes",
			scopes:       []interface{}{"marketing.read", "marketing.automation.read"},
			expectErrors: false,
		},
		{
			name:         "advanced scopes",
			scopes:       []interface{}{"subusers.create", "api_keys.read", "teammates.update"},
			expectErrors: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scopeSet := schema.NewSet(schema.HashString, tt.scopes)
			_, errors := validateTeammateScopes(scopeSet, "scopes")

			if tt.expectErrors {
				if len(errors) == 0 {
					t.Errorf("Expected errors but got none")
				} else if tt.errorContains != "" {
					found := false
					for _, err := range errors {
						if containsString(err.Error(), tt.errorContains) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected error containing '%s', but got: %v", tt.errorContains, errors)
					}
				}
			} else {
				if len(errors) > 0 {
					t.Errorf("Expected no errors but got: %v", errors)
				}
			}
		})
	}
}

func TestSanitizeScopes(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no automatic scopes",
			input:    []string{"mail.send", "templates.read"},
			expected: []string{"mail.send", "templates.read"},
		},
		{
			name:     "with 2fa_exempt",
			input:    []string{"mail.send", "2fa_exempt", "templates.read"},
			expected: []string{"mail.send", "templates.read"},
		},
		{
			name:     "with 2fa_required",
			input:    []string{"mail.send", "2fa_required"},
			expected: []string{"mail.send"},
		},
		{
			name:     "with both automatic scopes",
			input:    []string{"2fa_exempt", "mail.send", "2fa_required", "templates.read"},
			expected: []string{"mail.send", "templates.read"},
		},
		{
			name:     "only automatic scopes",
			input:    []string{"2fa_exempt", "2fa_required"},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeScopes(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d scopes, got %d", len(tt.expected), len(result))
			}

			for i, expected := range tt.expected {
				if i >= len(result) || result[i] != expected {
					t.Errorf("Expected scope %s at position %d, got %s", expected, i, result[i])
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(substr) <= len(s) && s[len(s)-len(substr):] == substr) ||
		(len(substr) <= len(s) && s[:len(substr)] == substr) ||
		containsStringHelper(s, substr))
}

func containsStringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
