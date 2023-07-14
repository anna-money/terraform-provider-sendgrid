# sendgrid_template_version

Provide a resource to manage a version of template.

## Example Usage

```hcl

	resource "sendgrid_template" "template" {
		name       = "my-template"
		generation = "dynamic"
	}

	resource "sendgrid_template_version" "template_version" {
		name                   = "my-template-version"
		template_id            = sendgrid_template.template.id
		active                 = 1
		html_content           = "<%body%>"
		generate_plain_content = true
		subject                = "subject"
	}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the transactional template version, max length: 100.
* `subject` - (Required) Subject of the new transactional template version, max length: 255.
* `template_id` - (Required) ID of the transactional template.
* `active` - (Optional) Set the version as the active version associated with the template. Only one version of a template can be active. The first version created for a template will automatically be set to Active. Allowed values: 0, 1.
* `editor` - (Optional) The editor used in the UI, allowed values: code (default), design.
* `generate_plain_content` - (Optional) If true (default), plain_content is always generated from html_content. If false, plain_content is not altered.
* `html_content` - (Optional) The HTML content of the version, maximum of 1048576 bytes allowed.
* `plain_content` - (Optional) Text/plain content of the transactional template version, maximum of 1048576 bytes allowed.
* `test_data` - (Optional) For dynamic templates only, the mock json data that will be used for template preview and test sends.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `thumbnail_url` - A thumbnail preview of the template's html content.
* `updated_at` - The date and time that this transactional template version was updated.


## Import

A template version can be imported, e.g.
```hcl
$ terraform import sendgrid_template_version.template_version templateID/templateVersionId
```
