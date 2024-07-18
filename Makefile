MOCKERY_CMD = mockery --name=HttpClient --dir=./client --output=./client/mocks --outpkg=mocks

# Define the target to generate mocks
.PHONY: mocks
mocks:
	@echo "Generating mocks..."
	$(MOCKERY_CMD)
	@echo "Mocks generated successfully."

.PHONY: clean
clean:
	@echo "Cleaning generated mocks..."
	rm -rf ./client/mocks
	@echo "Generated mocks cleaned."

.PHONY: clean-output
clean:
	@echo "Cleaning generated mocks..."
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

.PHONY: build
build:
	@echo "Building binary..."
	go build -o chatter
	@echo "Binary built successfully."