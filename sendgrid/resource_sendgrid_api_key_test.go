package sendgrid_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sendgrid "github.com/octoenergy/terraform-provider-sendgrid/sdk"
)

func TestAccSendgridAPIKeyBasic(t *testing.T) {
	name := "terraform-api-key-" + acctest.RandString(10)
	scopes := []string{"mail.send", "sender_verification_eligible"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSendgridAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSendgridAPIKeyConfigBasic(name, scopes),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridAPIKeyExists("sendgrid_api_key.new"),
				),
			},
		},
	})
}

func testAccCheckSendgridAPIKeyDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*sendgrid.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sendgrid_api_key" {
			continue
		}

		apiKeyID := rs.Primary.ID

		ctx := context.Background()
		_, err := c.DeleteAPIKey(ctx, apiKeyID)
		if err.Err != nil {
			return err.Err
		}
	}

	return nil
}

func testAccCheckSendgridAPIKeyConfigBasic(name string, scopes []string) string {
	return fmt.Sprintf(`
	resource "sendgrid_api_key" "api_key" {
		name = %s
		scopes = %s
	}
	`, name, scopes)
}

func testAccCheckSendgridAPIKeyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No apiKeyID set")
		}

		return nil
	}
}
