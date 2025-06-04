# Installation Guide

This guide covers all methods to install the SendGrid Terraform provider.

## Requirements

- **Terraform**: 0.13+ (recommended: latest stable)
- **Go**: 1.21+ (for building from source)
- **SendGrid Account**: With API key access

## Method 1: Terraform Registry (Recommended)

The easiest way to install the provider:

```hcl
terraform {
  required_version = ">= 0.13"
  required_providers {
    sendgrid = {
      source  = "arslanbekov/sendgrid"
      version = "~> 2.0"
    }
  }
}
```

Then run:

```bash
terraform init
```

## Method 2: Manual Installation

### Download Pre-built Binary

```bash
# Set your platform (linux_amd64, darwin_amd64, windows_amd64, etc.)
PLATFORM="linux_amd64"
VERSION="2.0.0"

# Download
wget "https://github.com/arslanbekov/terraform-provider-sendgrid/releases/download/v${VERSION}/terraform-provider-sendgrid_${VERSION}_${PLATFORM}.zip"

# Extract
unzip "terraform-provider-sendgrid_${VERSION}_${PLATFORM}.zip"

# Install to Terraform plugins directory
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/arslanbekov/sendgrid/${VERSION}/${PLATFORM}/
mv terraform-provider-sendgrid ~/.terraform.d/plugins/registry.terraform.io/arslanbekov/sendgrid/${VERSION}/${PLATFORM}/

# Make executable
chmod +x ~/.terraform.d/plugins/registry.terraform.io/arslanbekov/sendgrid/${VERSION}/${PLATFORM}/terraform-provider-sendgrid
```

### Alternative: Legacy Plugin Directory

For older Terraform versions:

```bash
# Download and extract (same as above)
# Then move to legacy location:
mkdir -p ~/.terraform.d/plugins/
mv terraform-provider-sendgrid ~/.terraform.d/plugins/
chmod +x ~/.terraform.d/plugins/terraform-provider-sendgrid
```

## Method 3: Build from Source

### Prerequisites

```bash
# Install Go 1.19+
go version

# Clone repository
git clone https://github.com/arslanbekov/terraform-provider-sendgrid.git
cd terraform-provider-sendgrid
```

### Build and Install

```bash
# Build
go build -o terraform-provider-sendgrid

# Install to plugins directory
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/arslanbekov/sendgrid/1.1.0/linux_amd64/
mv terraform-provider-sendgrid ~/.terraform.d/plugins/registry.terraform.io/arslanbekov/sendgrid/1.1.0/linux_amd64/
chmod +x ~/.terraform.d/plugins/registry.terraform.io/arslanbekov/sendgrid/1.1.0/linux_amd64/terraform-provider-sendgrid
```

### Development Build

```bash
# For development with hot reload
make install-dev

# Or manually
go build -o terraform-provider-sendgrid
mkdir -p ~/.terraform.d/plugins/
mv terraform-provider-sendgrid ~/.terraform.d/plugins/
```

## Method 4: Docker

Use the provider in a containerized environment:

```dockerfile
FROM hashicorp/terraform:1.5

# Download and install provider
RUN wget https://github.com/arslanbekov/terraform-provider-sendgrid/releases/latest/download/terraform-provider-sendgrid_linux_amd64.zip \
    && unzip terraform-provider-sendgrid_linux_amd64.zip \
    && mkdir -p ~/.terraform.d/plugins/ \
    && mv terraform-provider-sendgrid ~/.terraform.d/plugins/ \
    && chmod +x ~/.terraform.d/plugins/terraform-provider-sendgrid

WORKDIR /workspace
```

## Verification

After installation, verify the provider is available:

```bash
# Initialize a test configuration
cat > test.tf << 'EOF'
terraform {
  required_providers {
    sendgrid = {
      source  = "arslanbekov/sendgrid"
      version = "~> 1.1"
    }
  }
}

provider "sendgrid" {}
EOF

# Initialize (should download/find the provider)
terraform init

# Verify provider is loaded
terraform providers
```

Expected output:

```shell
Providers required by configuration:
.
└── provider[registry.terraform.io/arslanbekov/sendgrid] ~> 2.0
```

## Upgrading

### From Terraform Registry

Update the version constraint and run:

```hcl
terraform {
  required_providers {
    sendgrid = {
      source  = "arslanbekov/sendgrid"
      version = "~> 2.0"  # Updated version
    }
  }
}
```

```bash
terraform init -upgrade
```

### Manual Upgrade

1. Download the new version
2. Replace the binary in the plugins directory
3. Run `terraform init -upgrade`

## Troubleshooting

### Provider Not Found

```bash
# Error: Failed to query available provider packages
```

**Solutions:**

1. Check version constraint syntax
2. Verify internet connection for registry access
3. Try manual installation method

### Permission Denied

```bash
# Error: fork/exec ~/.terraform.d/plugins/terraform-provider-sendgrid: permission denied
```

**Solution:**

```bash
chmod +x ~/.terraform.d/plugins/terraform-provider-sendgrid
```

### Version Conflicts

```bash
# Error: Incompatible provider version
```

**Solutions:**

1. Update Terraform to latest version
2. Adjust version constraints
3. Clear `.terraform` directory and re-initialize

### Plugin Directory Issues

**Check plugin directory structure:**

```bash
ls -la ~/.terraform.d/plugins/
# Should show: terraform-provider-sendgrid (executable)

# For newer Terraform versions:
ls -la ~/.terraform.d/plugins/registry.terraform.io/arslanbekov/sendgrid/
```

## Next Steps

After installation:

1. [Set up Authentication](AUTHENTICATION.md)
2. [Check Examples](EXAMPLES.md)
3. [Browse Resources](RESOURCES.md)

## Need Help?

- [Report Installation Issues](https://github.com/arslanbekov/terraform-provider-sendgrid/issues)
- [Community Discussions](https://github.com/arslanbekov/terraform-provider-sendgrid/discussions)
- [Terraform Provider Documentation](https://www.terraform.io/docs/configuration/providers.html)
