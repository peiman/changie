package main

import (
	"context"
	"log"
	"os"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/peiman/changie/internal/mcp"
)

const (
	serverName    = "changie-mcp"
	serverVersion = "v1.0.0"
)

func main() {
	// Create MCP server
	server := mcpsdk.NewServer(&mcpsdk.Implementation{
		Name:    serverName,
		Version: serverVersion,
	}, nil)

	// Register changie_bump_version tool
	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "changie_bump_version",
		Description: "Bump semantic version (major, minor, or patch). Creates changelog entry, commits changes, and optionally pushes to remote.",
	}, mcp.BumpVersion)

	// Register changie_add_changelog tool
	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "changie_add_changelog",
		Description: "Add a new entry to the changelog under a specific section (added, changed, deprecated, removed, fixed, security).",
	}, mcp.AddChangelog)

	// Register changie_init tool
	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "changie_init",
		Description: "Initialize a new changie project by creating the CHANGELOG.md file and configuration.",
	}, mcp.Init)

	// Register changie_get_version tool
	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "changie_get_version",
		Description: "Get the current version from git tags.",
	}, mcp.GetVersion)

	// Log server start
	log.Printf("%s v%s starting...", serverName, serverVersion)
	log.Printf("Registered tools: changie_bump_version, changie_add_changelog, changie_init, changie_get_version")

	// Run server with stdio transport (standard for MCP)
	ctx := context.Background()
	transport := &mcpsdk.StdioTransport{}

	if err := server.Run(ctx, transport); err != nil {
		log.Printf("Server failed: %v", err)
		os.Exit(1)
	}
}
