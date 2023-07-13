# sendgrid_sso_certificate

Provide a resource to manage SSO certificates.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `integration_id` - (Required) An ID that matches an existing SSO integration.
* `public_certificate` - (Required) This public certificate allows SendGrid to verify that
					SAML requests it receives are signed by an IdP that it recognizes.


## Import

An SSO certificate can be imported, e.g.
```sh
$ terraform import sendgrid_sso_certificate.cert <certificate-id>
```
