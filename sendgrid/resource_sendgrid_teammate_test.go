package sendgrid_test

import (
	"context"
	"fmt"
	"testing"

	sendgrid "github.com/anna-money/terraform-provider-sendgrid/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSendgridTeammateBasic(t *testing.T) {
	email := "terraform-teammate-" + acctest.RandString(10) + "@example.com"
	scopes := []string{"mail.send", "marketing.read"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSendgridTeammateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSendgridTeammateConfigBasic(email, false, false, scopes),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridTeammateExists("sendgrid_teammate.test"),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "email", email),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "is_admin", "false"),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "is_sso", "false"),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "scopes.#", "2"),
				),
			},
		},
	})
}

func TestAccSendgridTeammateAdmin(t *testing.T) {
	email := "terraform-admin-" + acctest.RandString(10) + "@example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSendgridTeammateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSendgridTeammateConfigAdmin(email, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridTeammateExists("sendgrid_teammate.admin"),
					resource.TestCheckResourceAttr("sendgrid_teammate.admin", "email", email),
					resource.TestCheckResourceAttr("sendgrid_teammate.admin", "is_admin", "true"),
					resource.TestCheckResourceAttr("sendgrid_teammate.admin", "is_sso", "false"),
				),
			},
		},
	})
}

func TestAccSendgridTeammateUpdate(t *testing.T) {
	email := "terraform-update-" + acctest.RandString(10) + "@example.com"
	initialScopes := []string{"mail.send"}
	updatedScopes := []string{"mail.send", "marketing.read", "marketing.send"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSendgridTeammateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSendgridTeammateConfigBasic(email, false, false, initialScopes),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridTeammateExists("sendgrid_teammate.test"),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "scopes.#", "1"),
				),
			},
			{
				Config: testAccCheckSendgridTeammateConfigBasic(email, false, false, updatedScopes),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridTeammateExists("sendgrid_teammate.test"),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "scopes.#", "3"),
				),
			},
		},
	})
}

func TestAccSendgridTeammateWithRateLimiting(t *testing.T) {
	email := "terraform-rate-limit-" + acctest.RandString(10) + "@example.com"
	scopes := []string{"mail.send"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSendgridTeammateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSendgridTeammateConfigWithTimeouts(email, false, false, scopes),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridTeammateExists("sendgrid_teammate.rate_limit"),
					resource.TestCheckResourceAttr("sendgrid_teammate.rate_limit", "email", email),
				),
			},
		},
	})
}

func testAccCheckSendgridTeammateDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*sendgrid.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sendgrid_teammate" {
			continue
		}

		email := rs.Primary.ID
		ctx := context.Background()

		_, err := c.ReadUser(ctx, email)
		if err.StatusCode != 404 {
			return fmt.Errorf("teammate still exists: %s", email)
		}
	}

	return nil
}

func testAccCheckSendgridTeammateConfigBasic(email string, isAdmin, isSso bool, scopes []string) string {
	scopesStr := ""
	if len(scopes) > 0 {
		scopesStr = fmt.Sprintf(`scopes = %q`, scopes)
	}

	return fmt.Sprintf(`
resource "sendgrid_teammate" "test" {
	email    = "%s"
	is_admin = %t
	is_sso   = %t
	%s
}
`, email, isAdmin, isSso, scopesStr)
}

func testAccCheckSendgridTeammateConfigAdmin(email string, isAdmin, isSso bool) string {
	return fmt.Sprintf(`
resource "sendgrid_teammate" "admin" {
	email    = "%s"
	is_admin = %t
	is_sso   = %t
}
`, email, isAdmin, isSso)
}

func testAccCheckSendgridTeammateConfigWithTimeouts(email string, isAdmin, isSso bool, scopes []string) string {
	scopesStr := ""
	if len(scopes) > 0 {
		scopesStr = fmt.Sprintf(`scopes = %q`, scopes)
	}

	return fmt.Sprintf(`
resource "sendgrid_teammate" "rate_limit" {
	email    = "%s"
	is_admin = %t
	is_sso   = %t
	%s

	timeouts {
		create = "30m"
		update = "30m"
		delete = "30m"
	}
}
`, email, isAdmin, isSso, scopesStr)
}

func testAccCheckSendgridTeammateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No teammate email set")
		}

		c := testAccProvider.Meta().(*sendgrid.Client)
		ctx := context.Background()

		_, err := c.ReadUser(ctx, rs.Primary.ID)
		if err.Err != nil {
			return fmt.Errorf("teammate not found: %s", rs.Primary.ID)
		}

		return nil
	}
}
