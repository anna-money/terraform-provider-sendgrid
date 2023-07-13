package sendgrid_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/octoenergy/terraform-provider-sendgrid/sendgrid"
)

var testAccProviders map[string]*schema.Provider

var testAccProvider *schema.Provider

func init() {
	testAccProvider = sendgrid.Provider()
	testAccProviders = map[string]*schema.Provider{
		"sendgrid": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := sendgrid.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	_ = sendgrid.Provider()
}

func testAccPreCheck(t *testing.T) {
	t.Helper()

	if err := os.Getenv("SENDGRID_API_KEY"); err == "" {
		t.Fatal("SENDGRID_API_KEY must be set for acceptance tests")
	}
}
