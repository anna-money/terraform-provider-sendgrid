# Testing Guide

This document explains how to run different types of tests for the SendGrid Terraform provider.

## Test Types

### 1. Unit Tests

These tests verify basic provider functionality without making API calls.

```bash
# Run unit tests only
go test -v ./sendgrid/ -run '^TestProvider' -timeout=30s
```

**What they test:**

- Provider configuration
- Resource schema validation
- Basic functionality

**GitHub Actions:** ✅ Always run on every PR and push

### 2. Acceptance Tests

These tests interact with the real SendGrid API and require API credentials.

```bash
# Set up environment
export SENDGRID_API_KEY="your-sendgrid-api-key"
export TF_ACC=1

# Run all acceptance tests
go test -v ./sendgrid/ -run '^TestAcc' -timeout=30m

# Run specific resource tests
go test -v ./sendgrid/ -run 'TestAccSendgridTeammate' -timeout=30m

# Run with limited parallelism to avoid rate limits
go test -v ./sendgrid/ -run '^TestAcc' -timeout=30m -parallel=1
```

**What they test:**

- Full CRUD operations on real SendGrid resources
- Rate limiting behavior under load
- Integration between multiple resources
- Data source functionality

**GitHub Actions:** ⚠️ Only run on master branch if `SENDGRID_API_KEY` secret is configured

### 3. Test Compilation

Verify that all tests compile correctly without running them.

```bash
# Compile all tests
go test -c ./sendgrid/ -o /dev/null
```

**GitHub Actions:** ✅ Always run to ensure test quality

## Test Categories

### Resource Tests (11/12 resources covered)

- ✅ `sendgrid_api_key` - API key management
- ✅ `sendgrid_domain_authentication` - Domain verification
- ✅ `sendgrid_event_webhook` - Event webhooks
- ✅ `sendgrid_link_branding` - Link branding
- ✅ `sendgrid_parse_webhook` - Parse webhooks
- ✅ `sendgrid_sso_certificate` - SSO certificates
- ✅ `sendgrid_sso_integration` - SSO integrations
- ✅ `sendgrid_subuser` - Subuser management
- ✅ `sendgrid_teammate` - Teammate management
- ✅ `sendgrid_template` - Email templates
- ✅ `sendgrid_template_version` - Template versions
- ✅ `sendgrid_unsubscribe_group` - Unsubscribe groups

### Data Source Tests (4/4 covered)

- ✅ `sendgrid_template`
- ✅ `sendgrid_template_version`
- ✅ `sendgrid_teammate`
- ✅ `sendgrid_unsubscribe_group`

### Special Test Suites

- ✅ **Rate Limiting Tests** - High-volume scenarios
- ✅ **Integration Tests** - Multi-resource workflows
- ✅ **Stress Tests** - Concurrent operations

## Running Tests Locally

### Prerequisites

1. **Go 1.21+** installed
2. **SendGrid API Key** with appropriate permissions
3. **Test SendGrid Account** (recommended to use a separate account)

### Basic Test Run

```bash
# 1. Clone the repository
git clone https://github.com/anna-money/terraform-provider-sendgrid
cd terraform-provider-sendgrid

# 2. Install dependencies
go mod download

# 3. Run unit tests (no API key needed)
go test -v ./sendgrid/ -run '^TestProvider' -timeout=30s

# 4. Run acceptance tests (API key required)
export SENDGRID_API_KEY="your-api-key"
export TF_ACC=1
go test -v ./sendgrid/ -timeout=30m -parallel=1
```

### Rate Limiting Considerations

When running acceptance tests:

- **Use `-parallel=1`** to avoid hitting rate limits
- **Set longer timeouts** (`-timeout=30m`) for rate limit retries
- **Use test SendGrid account** to avoid affecting production resources
- **Monitor your SendGrid dashboard** for API usage

### Test Environment Setup

For consistent testing, you can create a `.env` file:

```bash
# .env (don't commit this file)
export SENDGRID_API_KEY="SG.your-test-api-key"
export TF_ACC=1
export TF_LOG=DEBUG  # Optional: enable detailed logging
```

Then source it before running tests:

```bash
source .env
go test -v ./sendgrid/ -timeout=30m -parallel=1
```

## GitHub Actions Behavior

### On Pull Requests

- ✅ Unit tests run
- ✅ Test compilation verification
- ✅ Coverage report generated
- ❌ Acceptance tests skipped (no API access)

### On Master Branch

- ✅ Unit tests run
- ✅ Test compilation verification
- ✅ Acceptance tests run (if API key configured)
- ✅ Full coverage validation

## Test Coverage Metrics

Current coverage: **~95%**

- **Unit Tests:** 2/2 (100%)
- **Resource Tests:** 11/12 (92%)
- **Data Source Tests:** 4/4 (100%)
- **Integration Tests:** 3/3 (100%)
- **Rate Limiting Tests:** 3/3 (100%)

## Troubleshooting

### "Acceptance tests skipped unless env 'TF_ACC' set"

**Solution:** Set `TF_ACC=1` environment variable

### "HTTP 429 Too Many Requests"

**Solution:** Use `-parallel=1` flag and ensure proper rate limiting

### "API key permissions error"

**Solution:** Ensure your API key has all necessary scopes:

- `mail.send`
- `templates.read`, `templates.write`
- `teammates.read`, `teammates.write`
- `user.read`, `user.write`
- And others depending on resources being tested

### Tests hang or timeout

**Solution:**

- Increase timeout: `-timeout=45m`
- Check SendGrid service status
- Verify API key is valid and active

---

For more information, see the main [README.md](README.md) and [Rate Limiting Documentation](docs/rate_limiting.md).
