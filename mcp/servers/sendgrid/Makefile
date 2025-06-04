.PHONY: build clean test test-init test-list test-call-no-key test-call install

# Build the sendgrid MCP server
build:
	go build -o bin/sendgrid-server .

# Clean build artifacts
clean:
	rm -f bin

# Run go tests
test:
	go test ./...

# Test server initialization (no API key required)
test-init: build
	@echo "Testing server initialization..."
	@echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2025-03-26", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0.0"}}}' | timeout 5s ./bin/sendgrid-server

# Test tool listing (no API key required)
test-list: build
	@echo "Testing tool listing..."
	@echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}' | timeout 5s ./bin/sendgrid-server

# Test tool call without API key (should return error)
test-call-no-key: build
	@echo "Testing email sending without API key (should fail gracefully)..."
	@echo '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "sendEmail", "arguments": {"from": "test@example.com", "to": "recipient@example.com", "subject": "Test", "textBody": "Hello!"}}}' | timeout 5s ./bin/sendgrid-server

# Test tool call with API key (requires SENDGRID_API_KEY, FROM_EMAIL, TO_EMAIL environment variables)
test-call: build
	@if [ -z "$(SENDGRID_API_KEY)" ] || [ -z "$(FROM_EMAIL)" ] || [ -z "$(TO_EMAIL)" ]; then \
		echo "Error: Please set SENDGRID_API_KEY, FROM_EMAIL, and TO_EMAIL environment variables"; \
		echo "Example: make test-call SENDGRID_API_KEY=your-key FROM_EMAIL=sender@example.com TO_EMAIL=recipient@example.com"; \
		exit 1; \
	fi
	@echo "Testing email sending with API key..."
	@SENDGRID_API_KEY=$(SENDGRID_API_KEY) echo '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "sendEmail", "arguments": {"from": "$(FROM_EMAIL)", "to": "$(TO_EMAIL)", "subject": "MCP Test Email", "textBody": "This is a test email sent via the SendGrid MCP server."}}}' | timeout 10s ./bin/sendgrid-server

# Install the server to a system location
install: build
	sudo cp bin/sendgrid-server /usr/local/bin/

# Show help
help:
	@echo "Available targets:"
	@echo "  build           - Build the sendgrid MCP server"
	@echo "  clean           - Clean build artifacts"
	@echo "  test            - Run go tests"
	@echo "  test-init       - Test server initialization"
	@echo "  test-list       - Test tool listing"
	@echo "  test-call-no-key - Test tool call without API key (should fail gracefully)"
	@echo "  test-call       - Test email sending (requires env vars)"
	@echo "  install         - Install server to /usr/local/bin"
	@echo "  help            - Show this help message"
