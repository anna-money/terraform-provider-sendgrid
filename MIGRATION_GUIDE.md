# Migration Guide: anna-money/sendgrid → arslanbekov/sendgrid

The SendGrid Terraform provider has been transferred from the `anna-money` organization to `arslanbekov` personal GitHub account.

## What Changed

- **GitHub Repository**: `anna-money/terraform-provider-sendgrid` → `arslanbekov/terraform-provider-sendgrid`
- **Terraform Registry**: `anna-money/sendgrid` → `arslanbekov/sendgrid`

## Migration Steps

### 1. Update Provider Source

In your `terraform` configuration block, update the provider source:

```hcl
# OLD - anna-money namespace
terraform {
  required_providers {
    sendgrid = {
      source  = "anna-money/sendgrid"
      version = "~> 2.0"
    }
  }
}

# NEW - arslanbekov namespace
terraform {
  required_providers {
    sendgrid = {
      source  = "arslanbekov/sendgrid"
      version = "~> 2.0"  # Use latest version
    }
  }
}
```

### 2. Reinitialize Terraform

After updating the provider source:

```bash
# Remove existing provider lock
rm .terraform.lock.hcl

# Reinitialize with new provider
terraform init

# Verify the change worked
terraform version
```

### 3. Verify State Compatibility

The provider functionality and resource schemas remain the same, so your existing state should work without changes:

```bash
# Check that your resources are still properly managed
terraform plan
```

## Version Compatibility

- **anna-money/sendgrid v1.0.x - v2.0.x**: All functionality preserved
- **arslanbekov/sendgrid v2.0.x+**: Same features with new namespace

## Benefits of Migration

- ✅ **Active Maintenance**: Continued development and bug fixes
- ✅ **Enhanced Error Handling**: Better validation and error messages
- ✅ **Comprehensive Documentation**: Updated with examples and troubleshooting
- ✅ **Community Support**: Direct access to maintainer

## Support

If you encounter any issues during migration:

- **GitHub Issues**: [arslanbekov/terraform-provider-sendgrid/issues](https://github.com/arslanbekov/terraform-provider-sendgrid/issues)
- **Documentation**: [Provider Documentation](https://registry.terraform.io/providers/arslanbekov/sendgrid/latest/docs)

## Backwards Compatibility

The `anna-money/sendgrid` provider will remain available in Terraform Registry for existing users, but new features and bug fixes will only be available in `arslanbekov/sendgrid`.

## Version 2.1.0 - Pending User Support

### New Features

#### Enhanced Teammate Management

Added support for pending users in teammate management. This resolves the issue where non-SSO teammates would cause Terraform errors on subsequent runs.

**New Field:**

- `user_status` (computed) - Shows "pending" for users who haven't accepted invitations, "active" for confirmed users

**Behavior Changes:**

- Non-SSO teammates now properly handle the invitation workflow
- Terraform operations (read, update, delete) work correctly for both pending and active users
- No more "username with email not found" errors on subsequent terraform runs

**Example:**

```hcl
resource "sendgrid_teammate" "example" {
  email    = "user@example.com"
  is_admin = false
  is_sso   = false
  scopes   = ["mail.send"]
}

# Check if user has accepted invitation
output "user_status" {
  value = sendgrid_teammate.example.user_status
}
```

**Migration Required:** None - this is a backward-compatible enhancement.

**Benefits:**

- Eliminates errors when managing non-SSO teammates
- Provides visibility into invitation status
- Allows proper lifecycle management of pending invitations

### Technical Details

The provider now:

1. Checks active teammates first
2. Falls back to pending invitations if user not found in active list
3. Supports updates and deletes for both pending and active users
4. Provides clear status indication through `user_status` field

This resolves the common issue where Terraform would fail on the second run with:

```shell
Error: request failed: resource not found. It may have been deleted outside of Terraform or the ID is incorrect
Original error: username with email user@example.com not found
```
