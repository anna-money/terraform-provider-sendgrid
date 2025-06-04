# Code Coverage with Codecov

This project uses [Codecov](https://codecov.io) to track and report code coverage metrics. Coverage reports are automatically generated and uploaded during CI/CD pipeline execution.

## Overview

- **Unit Tests Coverage**: Tracks coverage for provider and basic functionality tests
- **Acceptance Tests Coverage**: Tracks coverage for full integration tests (when API key is available)
- **Automatic Reporting**: Coverage reports are uploaded to Codecov on every push and pull request

## GitHub Actions Integration

Coverage is collected and uploaded automatically in the following scenarios:

### Unit Tests (Always Run)

- Runs on every push and pull request
- Covers basic provider functionality tests
- Uses flag: `unittests`

### Acceptance Tests (Master Branch Only)

- Runs only on pushes to master branch
- Requires `SENDGRID_API_KEY` secret to be configured
- Uses flag: `acceptancetests`

## Local Development

### Running Tests with Coverage

```bash
# Run unit tests with coverage
make test-coverage

# Run acceptance tests with coverage (requires SENDGRID_API_KEY)
export SENDGRID_API_KEY="your-api-key"
make testacc-coverage

# Generate HTML coverage report
make coverage-report

# Show total coverage percentage
make coverage-total

# Clean coverage files
make clean-coverage
```

### Manual Coverage Commands

```bash
# Unit tests coverage
go test ./sendgrid/ -run '^TestProvider' -timeout=30s -coverprofile=coverage.txt -covermode=atomic

# Acceptance tests coverage
TF_ACC=1 go test ./sendgrid/ -run '^TestAcc' -timeout=30m -parallel=1 -coverprofile=coverage-acceptance.txt -covermode=atomic

# Generate HTML report
go tool cover -html=coverage.txt -o coverage.html

# Show coverage by function
go tool cover -func=coverage.txt
```

## Configuration

### Codecov Settings (`codecov.yml`)

- **Target Coverage**: 80% for project, 70% for patches
- **Precision**: 2 decimal places
- **Ignored Files**: Test files, examples, tools, documentation
- **Flags**: Separate tracking for unit and acceptance tests

### GitHub Secrets Required

- `CODECOV_TOKEN`: Token for uploading coverage reports
- `SENDGRID_API_KEY`: API key for running acceptance tests (optional)

## Coverage Targets

- **Project Target**: 80% overall coverage
- **Patch Target**: 70% coverage for new changes
- **Current Status**: Check the badge in README.md

## Best Practices

1. **Write Tests First**: Ensure new features have corresponding tests
2. **Check Coverage Locally**: Use `make coverage-report` before committing
3. **Review Coverage Reports**: Check which lines are not covered
4. **Maintain Quality**: Aim for high coverage without sacrificing test quality

## Troubleshooting

### Coverage Not Uploading

- Verify `CODECOV_TOKEN` is set in GitHub secrets
- Check GitHub Actions logs for upload errors
- Ensure coverage files are generated correctly

### Low Coverage Warnings

- Review uncovered code paths
- Add tests for missing scenarios
- Consider if uncovered code is necessary

### Acceptance Test Coverage Missing

- Verify `SENDGRID_API_KEY` is configured
- Check if tests are running on master branch
- Review acceptance test execution logs

## Links

- [Codecov Dashboard](https://codecov.io/gh/arslanbekov/terraform-provider-sendgrid)
- [Coverage Badge](https://codecov.io/gh/arslanbekov/terraform-provider-sendgrid/branch/master/graph/badge.svg)
- [GitHub Actions](https://github.com/arslanbekov/terraform-provider-sendgrid/actions)
