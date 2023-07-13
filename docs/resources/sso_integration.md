# sendgrid_sso_integration

Provide a resource to manage SSO integrations.

**Note** To finalize the integration, a user must click through the 'enable integration'
workflow once after supplying all required fields including an SSO certificate via `aws_sso_certificate`.

## Example Usage

```hcl

	resource "sendgrid_sso_integration" "sso" {
		name    = "IdP"
		enabled = false

		signin_url  = "https://idp.com/signin"
		signout_url = "https://idp.com/signout"
		entity_id   = "https://idp.com/12345"
	}

```

## Argument Reference

The following arguments are supported:

* `enabled` - (Required) Indicates if the integration is enabled.
* `name` - (Required) The name of the integration.
* `entity_id` - (Optional) An identifier provided by your IdP to identify Twilio SendGrid in the SAML interaction.
					This is called the 'SAML Issuer ID' in the Twilio SendGrid UI.
* `signin_url` - (Optional) The IdP's SAML POST endpoint. This endpoint should receive requests
					and initiate an SSO login flow. This is called the 'Embed Link' in the Twilio SendGrid UI.
* `signout_url` - (Optional) This URL is relevant only for an IdP-initiated authentication flow.
					If a user authenticates from their IdP, this URL will return them to their IdP when logging out.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `audience_url` - The URL where your IdP should POST its SAML response.
					This is the Twilio SendGrid URL that is responsible for receiving and parsing a SAML assertion.
					This is the same URL as the Single Sign-On URL when using SendGrid.
* `completed_integration` - Indicates if the integration is complete.
* `single_signon_url` - The URL where your IdP should POST its SAML response.
					This is the Twilio SendGrid URL that is responsible for receiving and parsing a SAML assertion.
					This is the same URL as the Audience URL when using SendGrid.


## Import

A SSO integration can be imported, e.g.
```sh
$ terraform import sendgrid_sso_integration.sso <integration-id>
```
