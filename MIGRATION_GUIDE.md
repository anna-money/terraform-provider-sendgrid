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
      version = "~> 1.1"
    }
  }
}

# NEW - arslanbekov namespace
terraform {
  required_providers {
    sendgrid = {
      source  = "arslanbekov/sendgrid"
      version = "~> 1.2"  # Use latest version
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

- **anna-money/sendgrid v1.0.x - v1.1.x**: All functionality preserved
- **arslanbekov/sendgrid v1.2.x+**: Same features with new namespace

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
