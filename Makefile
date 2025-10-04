.PHONY: build test clean install fmt vet lint

# Build the provider
build:
	go build -o terraform-provider-euno .

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f terraform-provider-euno

# Install dependencies
install:
	go mod download
	go mod verify

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Run golangci-lint
lint:
	golangci-lint run

# Run all checks
check: fmt vet lint test

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o terraform-provider-euno-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o terraform-provider-euno-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o terraform-provider-euno-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o terraform-provider-euno-windows-amd64.exe .

# Install provider locally for testing
install-local: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/euno-ai/euno/0.0.1/linux_amd64
	cp terraform-provider-euno ~/.terraform.d/plugins/registry.terraform.io/euno-ai/euno/0.0.1/linux_amd64/

# Test with terraform
test-terraform: install-local
	cd test && terraform init
	cd test && terraform validate
	cd test && terraform plan

# Generate documentation
docs:
	go generate ./...

# Run goreleaser locally
release-local:
	goreleaser release --snapshot --clean

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build the provider"
	@echo "  test           - Run tests"
	@echo "  clean          - Clean build artifacts"
	@echo "  install        - Install dependencies"
	@echo "  fmt            - Format code"
	@echo "  vet            - Run go vet"
	@echo "  lint           - Run golangci-lint"
	@echo "  check          - Run all checks"
	@echo "  build-all      - Build for multiple platforms"
	@echo "  install-local  - Install provider locally for testing"
	@echo "  test-terraform - Test with terraform"
	@echo "  docs           - Generate documentation"
	@echo "  release-local  - Run goreleaser locally"
	@echo "  help           - Show this help"
