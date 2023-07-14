TEST?=$$(go list ./... | grep -v /sdk$$)
GOFMT_FILES?=$$(find . -name '*.go')
PGK_NAME=sendgrid

default: build

build: fmtcheck
	go install
	$(MAKE) --directory=scripts doc

test: fmtcheck
	@go test $(TEST) $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test ./$(PKG_NAME) -v $(TESTARGS) -timeout 1m

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
	tfplugindocs generate --rendered-provider-name 'SendGrid provider' --provider-name octoenergy/sendgrid

release:
	goreleaser release --rm-dist

.PHONY: build test testacc fmt fmtcheck lint golangci-lint sweep test-release release
