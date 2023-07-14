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

func TestAccSendgridSubuserBasic(t *testing.T) {
	username := "terraform-subuser-" + acctest.RandString(10)
	password := acctest.RandString(10)
	email := username + "@example.org"
	ips := []string{"127.0.0.1", "255.255.255.255"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSendgridSubuserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSendgridSubuserConfigBasic(username, password, email, ips),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridSubuserExists("sendgrid_subuser.new"),
				),
			},
		},
	})
}

func testAccCheckSendgridSubuserDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*sendgrid.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sendgrid_subuser" {
			continue
		}

		SubuserName := rs.Primary.ID
		ctx := context.Background()
		_, requestErr := c.DeleteSubuser(ctx, SubuserName)
		if requestErr.Err != nil {
			return requestErr.Err
		}
	}

	return nil
}

func testAccCheckSendgridSubuserConfigBasic(username, password, email string, ips []string) string {
	return fmt.Sprintf(`
	resource "sendgrid_subuser" "subuser" {
		username = %s
		password = %s
		email    = %s
		ips      = %s
	}
	`, username, password, email, ips)
}

func testAccCheckSendgridSubuserExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SubUserName set")
		}

		return nil
	}
}
