MOCKERY_CMD = mockery --name=HTTP --dir=./internal/client --output=./internal/client/mocks --outpkg=mocks
GOLANGCI_LINT :=  $(GOPATH)/bin/golangci-lint

.PHONY: mocks
mocks:
	@echo "Generating mocks..."
	$(MOCKERY_CMD)
	@echo "Mocks generated successfully."

.PHONY: clean
clean:
	@echo "Cleaning generated mocks..."
	rm -rf ./internal/client/mocks
	@echo "Generated mocks cleaned."

.PHONY: clean-output
clean-output:
	@echo "Cleaning output files..."
	rm -rf ./*.mp3
	@echo "Generated mocks cleaned."

.PHONY: install-mockery
install-mockery:
	@echo "Installing mockery..."
	go install github.com/vektra/mockery/v2@latest
	@echo "Mockery installed successfully."

.PHONY: test
test: install-mockery mocks
	@echo "Running tests..."
	go test ./...
	@echo "Tests completed."

.PHONY: lint
lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run

.PHONY: build
build:
	@echo "Building binary..."
	go build -o bin/chatter ./cmd/chatter
	@echo "Binary built successfully."