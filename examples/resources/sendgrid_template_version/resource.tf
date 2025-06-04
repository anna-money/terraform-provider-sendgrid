# Create a template first
resource "sendgrid_template" "welcome_email" {
  name       = "Welcome Email Template"
  generation = "dynamic"
}

# Basic template version with HTML content
resource "sendgrid_template_version" "welcome_v1" {
  name                   = "Welcome Email v1"
  template_id            = sendgrid_template.welcome_email.id
  active                 = 1
  subject                = "Welcome to {{company_name}}!"
  html_content           = <<EOF
<!DOCTYPE html>
<html>
<head>
    <title>Welcome!</title>
</head>
<body>
    <h1>Welcome {{first_name}}!</h1>
    <p>Thank you for joining {{company_name}}. We're excited to have you!</p>
    <p>Best regards,<br>The {{company_name}} Team</p>
</body>
</html>
EOF
  generate_plain_content = true
  editor                 = "code"
  test_data = jsonencode({
    first_name   = "John"
    company_name = "Acme Corp"
  })
}
