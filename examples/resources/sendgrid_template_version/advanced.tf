# Advanced template version with custom plain content
resource "sendgrid_template" "newsletter" {
  name       = "Monthly Newsletter"
  generation = "dynamic"
}

resource "sendgrid_template_version" "newsletter_v2" {
  name        = "Newsletter v2.0 - Rich HTML"
  template_id = sendgrid_template.newsletter.id
  active      = 1
  subject     = "{{month}} Newsletter - {{company_name}}"

  html_content = <<EOF
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{month}} Newsletter</title>
    <style>
        .container { max-width: 600px; margin: 0 auto; font-family: Arial, sans-serif; }
        .header { background-color: #007bff; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; }
        .footer { background-color: #f8f9fa; padding: 15px; text-align: center; font-size: 12px; }
        .button { background-color: #28a745; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{company_name}} Newsletter</h1>
            <p>{{month}} Edition</p>
        </div>
        <div class="content">
            <h2>Hello {{customer_name}}!</h2>
            <p>{{newsletter_intro}}</p>

            {{#each articles}}
            <h3>{{title}}</h3>
            <p>{{summary}}</p>
            <a href="{{url}}" class="button">Read More</a>
            <hr>
            {{/each}}
        </div>
        <div class="footer">
            <p>© {{year}} {{company_name}}. All rights reserved.</p>
            <p><a href="{{unsubscribe_url}}">Unsubscribe</a></p>
        </div>
    </div>
</body>
</html>
EOF

  # Custom plain text content instead of auto-generated
  plain_content = <<EOF
{{company_name}} Newsletter - {{month}} Edition

Hello {{customer_name}}!

{{newsletter_intro}}

{{#each articles}}
{{title}}
{{summary}}
Read more: {{url}}

{{/each}}

© {{year}} {{company_name}}. All rights reserved.
Unsubscribe: {{unsubscribe_url}}
EOF

  generate_plain_content = false # Use custom plain content
  editor                 = "design"

  test_data = jsonencode({
    month            = "January"
    company_name     = "Tech Weekly"
    customer_name    = "Jane Doe"
    year             = "2024"
    newsletter_intro = "Welcome to our monthly roundup of the latest tech news and insights!"
    unsubscribe_url  = "https://example.com/unsubscribe"
    articles = [
      {
        title   = "AI Revolution in 2024"
        summary = "Exploring the latest developments in artificial intelligence."
        url     = "https://example.com/article1"
      },
      {
        title   = "Cloud Computing Trends"
        summary = "What's next for cloud infrastructure and services."
        url     = "https://example.com/article2"
      }
    ]
  })
}
