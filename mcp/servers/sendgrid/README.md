# SendGrid MCP Server

A Model Context Protocol (MCP) server that provides email sending capabilities using SendGrid.

## Overview

This MCP server exposes a `sendEmail` tool that allows LLM applications to send emails through the SendGrid API. It converts the existing SendGrid email functionality into an MCP-compatible server using the [mcp-go](https://github.com/mark3labs/mcp-go) library.

## Features

- **Email Sending**: Send emails with text and/or HTML content
- **Multiple Recipients**: Support for comma-separated recipient lists
- **Flexible Content**: Support for both plain text and HTML email bodies
- **Validation**: Input validation for required fields
- **Error Handling**: Comprehensive error reporting
- **Secure Configuration**: API key loaded from environment variable
- **Graceful Degradation**: Server starts without API key, returns helpful error when tool is called

## Installation

### Prerequisites

- Go 1.23.4 or later
- SendGrid API key (for sending emails)

### Building

```bash
go build -o sendgrid-server .
```

## Configuration

### SendGrid API Key

You'll need a SendGrid API key to send emails. You can obtain one from the [SendGrid Console](https://app.sendgrid.com/settings/api_keys).

Set the API key as an environment variable:

```bash
export SENDGRID_API_KEY="your-sendgrid-api-key-here"
```

**Note**: The server will start successfully without the API key, but will return a helpful error message when the `sendEmail` tool is called without it. This allows the server to be configured and tested even when the API key is not immediately available.

## Usage

### Running the Server

The server communicates via stdio (standard input/output):

```bash
# The server can start without the API key
./sendgrid-server

# Or with the API key for full functionality
export SENDGRID_API_KEY="your-sendgrid-api-key-here"
./sendgrid-server
```

### Tool: sendEmail

Send an email using SendGrid.

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `from` | string | Yes | Sender email address |
| `fromName` | string | No | Sender name (defaults to 'SendGrid MCP Server' if not provided) |
| `to` | string | Yes | Recipient email address(es), comma-separated for multiple recipients |
| `subject` | string | Yes | Email subject |
| `textBody` | string | No | Plain text email body (optional if htmlBody is provided) |
| `htmlBody` | string | No | HTML email body (optional if textBody is provided) |

#### Notes

- The SendGrid API key is read from the `SENDGRID_API_KEY` environment variable
- If the API key is not set, the tool will return a helpful error message
- Either `textBody` or `htmlBody` (or both) must be provided
- Multiple recipients can be specified by separating email addresses with commas
- The server validates all inputs and provides detailed error messages

#### Example Tool Call

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "sendEmail",
    "arguments": {
      "from": "sender@example.com",
      "fromName": "Your Name",
      "to": "recipient@example.com",
      "subject": "Test Email",
      "textBody": "This is a test email sent via MCP.",
      "htmlBody": "<p>This is a <strong>test email</strong> sent via MCP.</p>"
    }
  }
}
```

## Integration with MCP Clients

This server can be integrated with any MCP-compatible client. Here's an example configuration for common clients:

### Claude Desktop

Add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "sendgrid": {
      "command": "/path/to/sendgrid-server",
      "env": {
        "SENDGRID_API_KEY": "your-sendgrid-api-key-here"
      }
    }
  }
}
```

### Other MCP Clients

Refer to your MCP client's documentation for configuration instructions. The server uses stdio transport and follows the standard MCP protocol. Set the `SENDGRID_API_KEY` environment variable when you want to enable email sending functionality.

## Development

### Project Structure

```
.
├── main.go           # MCP server implementation
├── cmd/
│   └── send.go      # Core SendGrid email functionality
├── go.mod           # Go module dependencies
├── go.sum           # Go module checksums
└── README.md        # This file
```

### Dependencies

- `github.com/mark3labs/mcp-go` - MCP Go implementation
- `github.com/sendgrid/sendgrid-go` - SendGrid Go SDK

### Testing

You can test the server manually using JSON-RPC messages:

1. **Initialize the server (no API key required):**
```bash
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2025-03-26", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0.0"}}}' | ./sendgrid-server
```

2. **List available tools (no API key required):**
```bash
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}' | ./sendgrid-server
```

3. **Test without API key (should return helpful error):**
```bash
echo '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "sendEmail", "arguments": {"from": "test@example.com", "to": "recipient@example.com", "subject": "Test", "textBody": "Hello!"}}}' | ./sendgrid-server
```

4. **Test with API key:**
```bash
export SENDGRID_API_KEY="your-sendgrid-api-key-here"
echo '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "sendEmail", "arguments": {"from": "test@example.com", "to": "recipient@example.com", "subject": "Test", "textBody": "Hello!"}}}' | ./sendgrid-server
```

### Using the Makefile

The project includes a Makefile with helpful targets:

```bash
# Build the server
make build

# Test server initialization and tool listing (no API key required)
make test-init
make test-list

# Test tool call without API key (should fail gracefully)
make test-call-no-key

# Test actual email sending (requires environment variables)
make test-call SENDGRID_API_KEY=your-key FROM_EMAIL=sender@example.com TO_EMAIL=recipient@example.com

# Show all available targets
make help
```

## Error Handling

The server provides detailed error messages for common issues:

- Missing `SENDGRID_API_KEY` environment variable (returned as tool error, not startup failure)
- Missing required parameters
- Invalid email addresses
- SendGrid API errors
- Network connectivity issues

All errors are returned as tool results with `isError: true` to allow LLMs to see and respond to the errors appropriately. The server follows MCP protocol standards and only sends valid JSON-RPC messages over stdout/stderr.

## Security Considerations

- API key is read from environment variable and not exposed in tool parameters
- Server can start without API key for testing and configuration
- Input validation prevents basic injection attacks
- The server runs with minimal privileges
- All external API calls are made through the official SendGrid SDK
- Only valid MCP protocol messages are sent over stdout/stderr

## License

This project uses the same license as the parent repository.

## Contributing

Contributions are welcome! Please ensure that:

1. Code follows Go best practices
2. All functionality is properly tested
3. Documentation is updated for any new features
4. Security considerations are addressed
5. MCP protocol compliance is maintained

## Support

For issues related to:
- **MCP Protocol**: See the [MCP specification](https://modelcontextprotocol.io/)
- **mcp-go Library**: See the [mcp-go repository](https://github.com/mark3labs/mcp-go)
- **SendGrid API**: See the [SendGrid documentation](https://docs.sendgrid.com/)
