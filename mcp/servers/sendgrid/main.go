package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/obot-platform/tools/sendgrid/cmd"
)

var sendGridAPIKey string

func main() {
	// Read SendGrid API key from environment variable
	sendGridAPIKey = os.Getenv("SENDGRID_API_KEY")
	// Note: We don't exit here if the API key is missing - we'll handle it in the tool handler

	// Create MCP server
	s := server.NewMCPServer(
		"sendgrid-server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Define the sendEmail tool
	sendEmailTool := mcp.NewTool(
		"sendEmail",
		mcp.WithDescription("Send an email using SendGrid"),
		mcp.WithString("from", mcp.Description("Sender email address"), mcp.Required()),
		mcp.WithString("fromName", mcp.Description("Sender name (optional, defaults to 'SendGrid MCP Server')")),
		mcp.WithString("to", mcp.Description("Recipient email address(es), comma-separated for multiple recipients"), mcp.Required()),
		mcp.WithString("subject", mcp.Description("Email subject"), mcp.Required()),
		mcp.WithString("textBody", mcp.Description("Plain text email body (optional if htmlBody is provided)")),
		mcp.WithString("htmlBody", mcp.Description("HTML email body (optional if textBody is provided)")),
	)

	// Add the tool with its handler
	s.AddTool(sendEmailTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Check if API key is available
		if sendGridAPIKey == "" {
			return mcp.NewToolResultError("SendGrid API key not configured. Please set the SENDGRID_API_KEY environment variable."), nil
		}

		// Parse arguments using the request's built-in methods
		from := req.GetString("from", "")
		fromName := req.GetString("fromName", "")
		to := req.GetString("to", "")
		subject := req.GetString("subject", "")
		textBody := req.GetString("textBody", "")
		htmlBody := req.GetString("htmlBody", "")

		// Call the existing Send function with the API key from environment
		result, err := cmd.Send(ctx, sendGridAPIKey, from, fromName, to, subject, textBody, htmlBody)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to send email: %v", err)), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	// Start the server using stdio transport
	if err := server.ServeStdio(s); err != nil {
		// This is the only place we can exit with an error, but it should be rare
		// and only happen if there's a fundamental issue with the MCP server itself
		os.Exit(1)
	}
}
