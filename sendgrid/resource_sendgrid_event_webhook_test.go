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

func TestAccSendgridEventWebhookBasic(t *testing.T) {
	url := "https://example-" + acctest.RandString(10) + ".com/webhook"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSendgridEventWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSendgridEventWebhookConfigBasic(url),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridEventWebhookExists("sendgrid_event_webhook.test"),
					resource.TestCheckResourceAttr("sendgrid_event_webhook.test", "url", url),
					resource.TestCheckResourceAttr("sendgrid_event_webhook.test", "enabled", "true"),
				),
			},
		},
	})
}

func TestAccSendgridEventWebhookWithEvents(t *testing.T) {
	url := "https://events-" + acctest.RandString(10) + ".com/webhook"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSendgridEventWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSendgridEventWebhookConfigWithEvents(url),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridEventWebhookExists("sendgrid_event_webhook.events"),
					resource.TestCheckResourceAttr("sendgrid_event_webhook.events", "url", url),
					resource.TestCheckResourceAttr("sendgrid_event_webhook.events", "enabled", "true"),
					resource.TestCheckResourceAttr("sendgrid_event_webhook.events", "bounce", "true"),
					resource.TestCheckResourceAttr("sendgrid_event_webhook.events", "click", "true"),
					resource.TestCheckResourceAttr("sendgrid_event_webhook.events", "deferred", "false"),
				),
			},
		},
	})
}

func TestAccSendgridEventWebhookUpdate(t *testing.T) {
	url := "https://update-" + acctest.RandString(10) + ".com/webhook"
	urlUpdated := "https://updated-" + acctest.RandString(10) + ".com/webhook"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSendgridEventWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSendgridEventWebhookConfigBasic(url),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridEventWebhookExists("sendgrid_event_webhook.test"),
					resource.TestCheckResourceAttr("sendgrid_event_webhook.test", "url", url),
				),
			},
			{
				Config: testAccCheckSendgridEventWebhookConfigBasic(urlUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridEventWebhookExists("sendgrid_event_webhook.test"),
					resource.TestCheckResourceAttr("sendgrid_event_webhook.test", "url", urlUpdated),
				),
			},
		},
	})
}

func TestAccSendgridEventWebhookWithRateLimiting(t *testing.T) {
	url := "https://rate-limit-" + acctest.RandString(10) + ".com/webhook"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSendgridEventWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSendgridEventWebhookConfigWithTimeouts(url),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridEventWebhookExists("sendgrid_event_webhook.rate_limit"),
					resource.TestCheckResourceAttr("sendgrid_event_webhook.rate_limit", "url", url),
				),
			},
		},
	})
}

func testAccCheckSendgridEventWebhookDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*sendgrid.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sendgrid_event_webhook" {
			continue
		}

		ctx := context.Background()
		_, err := c.ReadEventWebhook(ctx)
		if err.StatusCode != 404 && err.Err == nil {
			return fmt.Errorf("event webhook still exists")
		}
	}

	return nil
}

func testAccCheckSendgridEventWebhookConfigBasic(url string) string {
	return fmt.Sprintf(`
resource "sendgrid_event_webhook" "test" {
	url     = "%s"
	enabled = true
}
`, url)
}

func testAccCheckSendgridEventWebhookConfigWithEvents(url string) string {
	return fmt.Sprintf(`
resource "sendgrid_event_webhook" "events" {
	url      = "%s"
	enabled  = true
	bounce   = true
	click    = true
	deferred = false
	delivered = true
	dropped  = false
}
`, url)
}

func testAccCheckSendgridEventWebhookConfigWithTimeouts(url string) string {
	return fmt.Sprintf(`
resource "sendgrid_event_webhook" "rate_limit" {
	url     = "%s"
	enabled = true

	timeouts {
		create = "30m"
		update = "30m"
		delete = "30m"
	}
}
`, url)
}

func testAccCheckSendgridEventWebhookExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No event webhook ID set")
		}

		c := testAccProvider.Meta().(*sendgrid.Client)
		ctx := context.Background()

		_, err := c.ReadEventWebhook(ctx)
		if err.Err != nil {
			return fmt.Errorf("event webhook not found")
		}

		return nil
	}
}
