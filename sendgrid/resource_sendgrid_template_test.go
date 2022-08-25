package sendgrid_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sendgrid "github.com/trois-six/terraform-provider-sendgrid/sdk"
)

func TestAccSendgridTemplateBasic(t *testing.T) {
	name := "terraform-template-" + acctest.RandString(10)
	generation := "dynamic"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSendgridTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSendgridTemplateConfigBasic(name, generation),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridTemplateExists("sendgrid_template.new"),
				),
			},
		},
	})
}

func testAccCheckSendgridTemplateDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*sendgrid.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sendgrid_template" {
			continue
		}

		templateID := rs.Primary.ID

		_, err := c.DeleteTemplate(templateID)
		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckSendgridTemplateConfigBasic(name, generation string) string {
	return fmt.Sprintf(`
	resource "sendgrid_template" "template" {
		name = %s
		generation = %s
	}
	`, name, generation)
}

func testAccCheckSendgridTemplateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No templateID set")
		}

		return nil
	}
}
