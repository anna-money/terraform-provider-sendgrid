default: doc

build: fmtcheck
	@go build -o gendocs

doc: build
	@echo "==> Building documentation..."
	./gendocs

fmt:
	@gofmt -s -w .

fmtcheck:
	@sh -c "'$(CURDIR)/gofmtcheck.sh'"

lint: golangci-lint

golangci-lint:
	@golangci-lint run ./...

sweep:
	@rm -f gendocs

.PHONY: build doc fmt fmtcheck lint golangci-lint sweep