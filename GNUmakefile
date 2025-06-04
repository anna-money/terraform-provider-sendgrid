TEST?=$$(go list ./... | grep -v /sdk$$)
GOFMT_FILES?=$$(find . -name '*.go')
PKG_NAME=sendgrid

default: build

build: fmtcheck
	go install
	$(MAKE) --directory=scripts doc

test: fmtcheck
	@go test $(TEST) $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test ./$(PKG_NAME) -v $(TESTARGS) -timeout 1m

# Coverage targets
test-coverage: fmtcheck
	@echo "==> Running unit tests with coverage..."
	@go test ./$(PKG_NAME) -run '^TestProvider' -timeout=30s -coverprofile=coverage.txt -covermode=atomic

testacc-coverage: fmtcheck
	@echo "==> Running acceptance tests with coverage..."
	TF_ACC=1 go test ./$(PKG_NAME) -run '^TestAcc' -timeout=30m -parallel=1 -coverprofile=coverage-acceptance.txt -covermode=atomic

coverage-report: test-coverage
	@echo "==> Generating coverage report..."
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"

coverage-total: test-coverage
	@echo "==> Total coverage:"
	@go tool cover -func=coverage.txt | grep total

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./$(PKG_NAME)
	$(MAKE) --directory=scripts $@

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint: golangci-lint

golangci-lint:
	@echo "==> Checking source code against golangci-lint..."
	@golangci-lint run ./$(PKG_NAME)/...
	@$(MAKE) --directory=scripts $@

sweep:
	@rm -rf "$(CURDIR)/dist"
	@$(MAKE) --directory=scripts $@

test-release:
	goreleaser --snapshot --skip-publish --rm-dist

doc:
	tfplugindocs generate --rendered-provider-name 'SendGrid provider' --provider-name sendgrid

docs: doc

release:
	goreleaser release --rm-dist

clean-coverage:
	@echo "==> Cleaning coverage files..."
	@rm -f coverage.txt coverage-acceptance.txt coverage.html

.PHONY: build test testacc test-coverage testacc-coverage coverage-report coverage-total fmt fmtcheck lint golangci-lint sweep test-release doc docs release clean-coverage
