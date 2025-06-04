package sendgrid_test

import (
	"context"
	"fmt"
	"testing"

	sendgrid "github.com/arslanbekov/terraform-provider-sendgrid/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSendgridDomainAuthenticationBasic(t *testing.T) {
	domain := "test-" + acctest.RandString(10) + ".example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSendgridDomainAuthenticationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSendgridDomainAuthenticationConfigBasic(domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridDomainAuthenticationExists("sendgrid_domain_authentication.test"),
					resource.TestCheckResourceAttr("sendgrid_domain_authentication.test", "domain", domain),
					resource.TestCheckResourceAttr("sendgrid_domain_authentication.test", "default", "false"),
					resource.TestCheckResourceAttr("sendgrid_domain_authentication.test", "automatic_security", "false"),
				),
			},
		},
	})
}

func TestAccSendgridDomainAuthenticationCustomIPs(t *testing.T) {
	domain := "custom-" + acctest.RandString(10) + ".example.com"
	customIPs := []string{"192.168.1.1", "192.168.1.2"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSendgridDomainAuthenticationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSendgridDomainAuthenticationConfigCustomIPs(domain, customIPs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridDomainAuthenticationExists("sendgrid_domain_authentication.custom"),
					resource.TestCheckResourceAttr("sendgrid_domain_authentication.custom", "domain", domain),
					resource.TestCheckResourceAttr("sendgrid_domain_authentication.custom", "custom_ips.#", "2"),
				),
			},
		},
	})
}

func TestAccSendgridDomainAuthenticationWithRateLimiting(t *testing.T) {
	domain := "rate-limit-" + acctest.RandString(10) + ".example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSendgridDomainAuthenticationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSendgridDomainAuthenticationConfigWithTimeouts(domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridDomainAuthenticationExists("sendgrid_domain_authentication.rate_limit"),
					resource.TestCheckResourceAttr("sendgrid_domain_authentication.rate_limit", "domain", domain),
				),
			},
		},
	})
}

func testAccCheckSendgridDomainAuthenticationDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*sendgrid.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sendgrid_domain_authentication" {
			continue
		}

		domainID := rs.Primary.ID
		ctx := context.Background()

		_, err := c.ReadDomainAuthentication(ctx, domainID)
		if err.StatusCode != 404 {
			return fmt.Errorf("domain authentication still exists: %s", domainID)
		}
	}

	return nil
}

func testAccCheckSendgridDomainAuthenticationConfigBasic(domain string) string {
	return fmt.Sprintf(`
resource "sendgrid_domain_authentication" "test" {
	domain             = "%s"
	default            = false
	automatic_security = false
}
`, domain)
}

func testAccCheckSendgridDomainAuthenticationConfigCustomIPs(domain string, customIPs []string) string {
	return fmt.Sprintf(`
resource "sendgrid_domain_authentication" "custom" {
	domain             = "%s"
	default            = false
	automatic_security = false
	custom_ips         = %q
}
`, domain, customIPs)
}

func testAccCheckSendgridDomainAuthenticationConfigWithTimeouts(domain string) string {
	return fmt.Sprintf(`
resource "sendgrid_domain_authentication" "rate_limit" {
	domain             = "%s"
	default            = false
	automatic_security = false

	timeouts {
		create = "30m"
		update = "30m"
		delete = "30m"
	}
}
`, domain)
}

func testAccCheckSendgridDomainAuthenticationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No domain authentication ID set")
		}

		c := testAccProvider.Meta().(*sendgrid.Client)
		ctx := context.Background()

		_, err := c.ReadDomainAuthentication(ctx, rs.Primary.ID)
		if err.Err != nil {
			return fmt.Errorf("domain authentication not found: %s", rs.Primary.ID)
		}

		return nil
	}
}
