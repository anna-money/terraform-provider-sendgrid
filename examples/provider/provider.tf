terraform {
  required_providers {
    sendgrid = {
      version = "1.0.0"
      source  = "registry.terraform.io/anna-money/sendgrid"
    }
  }
}

provider "sendgrid" {
}
