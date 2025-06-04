# Bulk teammate creation using for_each
locals {
  developers = {
    "dev1" = {
      email  = "dev1@example.com"
      scopes = ["mail.send", "templates.read"]
    }
    "dev2" = {
      email  = "dev2@example.com"
      scopes = ["mail.send", "templates.read", "stats.read"]
    }
    "dev3" = {
      email  = "dev3@example.com"
      scopes = ["mail.send"]
    }
  }
}

resource "sendgrid_teammate" "developers" {
  for_each = local.developers

  email    = each.value.email
  is_admin = false
  is_sso   = false
  scopes   = each.value.scopes

  # Important: Use higher timeouts when creating multiple resources
  timeouts {
    create = "45m"
    update = "45m"
    delete = "45m"
  }
}
