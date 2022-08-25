/*
Provide a resource to manage SSO certificates.
Example Usage
```hcl
resource "sendgrid_sso_integration" "sso" {
	name    = "IdP"
	enabled = true

	signin_url  = "https://idp.com/signin"
	signout_url = "https://idp.com/signout"
	entity_id   = "https://idp.com/12345"
}

resource "sendgrid_sso_certificate" "cert" {
	integration_id = sendgrid_sso_integration.sso.id
	public_certificate = <<EOF
-----BEGIN CERTIFICATE-----
...
EOF
}
```
Import
An SSO certificate can be imported, e.g.
```sh
$ terraform import sendgrid_sso_certificate.cert <certificate-id>
```
*/
package sendgrid

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sendgrid "github.com/trois-six/terraform-provider-sendgrid/sdk"
)

func resourceSendgridSSOCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSendgridSSOCertificateCreate,
		ReadContext:   resourceSendgridSSOCertificateRead,
		UpdateContext: resourceSendgridSSOCertificateUpdate,
		DeleteContext: resourceSendgridSSOCertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"public_certificate": {
				Type: schema.TypeString,
				Description: `This public certificate allows SendGrid to verify that
					SAML requests it receives are signed by an IdP that it recognizes.`,
				Required: true,
			},
			"integration_id": {
				Type:        schema.TypeString,
				Description: "An ID that matches an existing SSO integration.",
				Required:    true,
			},
		},
	}
}

func resourceSendgridSSOCertificateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	publicCertificate := d.Get("public_certificate").(string)
	integrationID := d.Get("integration_id").(string)

	apiKeyStruct, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return c.CreateSSOCertificate(publicCertificate, integrationID)
	})
	if err != nil {
		return diag.FromErr(err)
	}

	certificate := apiKeyStruct.(*sendgrid.SSOCertificate)

	d.SetId(fmt.Sprint(certificate.ID))

	return resourceSendgridSSOCertificateRead(ctx, d, m)
}

func resourceSendgridSSOCertificateRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	certificate, requestErr := c.ReadSSOCertificate(d.Id())

	if requestErr.Err != nil {
		return diag.FromErr(requestErr.Err)
	}

	//nolint:errcheck
	d.Set("public_certificate", certificate.PublicCertificate)
	//nolint:errcheck
	d.Set("integration_id", certificate.IntegrationID)

	return nil
}

func resourceSendgridSSOCertificateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	id := d.Id()
	publicCertificate := d.Get("public_certificate").(string)
	integrationID := d.Get("integration_id").(string)

	_, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return c.UpdateSSOCertificate(id, publicCertificate, integrationID)
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSendgridSSOCertificateRead(ctx, d, m)
}

func resourceSendgridSSOCertificateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	_, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return c.DeleteSSOCertificate(fmt.Sprint(d.Id()))
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
