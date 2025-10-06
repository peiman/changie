# .ai/ - AI Agent Resources

This directory contains resources specifically designed for AI agents, LLMs, and automation tools working with changie.

## Structure

```
.ai/
├── README.md                          # This file
├── context.md                         # Project context and agent guidance
├── prompts/                           # Task-specific prompts
│   ├── release-new-version.md        # Guide for releasing versions
│   └── add-changelog-entry.md        # Guide for adding changelog entries
└── workflows/                         # Multi-step workflows
    └── complete-release-workflow.md  # Full release process
```

## Purpose

The `.ai/` directory follows emerging 2025 patterns for AI-augmented development, providing:

1. **Structured Context** - Help agents understand the project quickly
2. **Task Prompts** - Reusable prompts for common operations
3. **Workflows** - Step-by-step guides for complex processes

## Files

### `context.md`

**Purpose:** Provide essential project context to AI agents

**Contains:**
- What changie does and how it works
- Core capabilities and commands
- When to use each version bump type
- Important constraints and prerequisites
- JSON output structure
- Common workflows and error handling
- Configuration options
- Best practices

**Use this:** When an agent needs to understand changie's capabilities and constraints

### `prompts/release-new-version.md`

**Purpose:** Guide agents through version release process

**Contains:**
- Step-by-step release instructions
- Prerequisites verification
- Bump type determination logic
- Success criteria
- Error handling

**Use this:** When releasing a new version

### `prompts/add-changelog-entry.md`

**Purpose:** Guide agents in adding changelog entries

**Contains:**
- Category selection logic
- Entry formatting guidelines
- Decision tree for section choice
- Best practices and examples
- Common mistakes to avoid

**Use this:** When adding changelog entries during development

### `workflows/complete-release-workflow.md`

**Purpose:** Complete end-to-end release workflow

**Contains:**
- Multi-step release process
- Pre-release validation
- Post-release actions
- Error recovery procedures
- Rollback strategy
- Automation scripts

**Use this:** For comprehensive release management

## For AI Agents

### Quick Start

1. **Read `context.md` first** - Understand changie capabilities
2. **Use appropriate prompt** - Select based on task
3. **Follow workflow if complex** - Use workflows for multi-step tasks
4. **Always use `--json` flag** - For reliable output parsing
5. **Check success field** - Before proceeding with results

### Integration Pattern

```python
# Example: Agent using changie

# 1. Load context
context = read_file(".ai/context.md")

# 2. Understand constraints
# - Clean working directory required
# - Must be on main/master or use --allow-any-branch
# - Need unreleased changelog entries

# 3. Execute with JSON
result = run_command("changie bump patch --json")
data = json.loads(result)

# 4. Verify success
if data["success"]:
    version = data["new_version"]
    print(f"Released {version}")
else:
    error = data["error"]
    # Error includes solution hints
    handle_error(error)
```

### Common Agent Tasks

**Task: Release a version**
→ Use: `prompts/release-new-version.md`

**Task: Add changelog entry**
→ Use: `prompts/add-changelog-entry.md`

**Task: Complete release with deployment**
→ Use: `workflows/complete-release-workflow.md`

**Task: Understand changie capabilities**
→ Use: `context.md`

## For MCP (Model Context Protocol)

Changie includes an **official MCP server** that exposes changelog operations as tools for AI agents.

### MCP Server (Go Implementation)

**Built with:** Official MCP Go SDK v1.0.0 (`github.com/modelcontextprotocol/go-sdk`)

**Available Tools:**

1. **`changie_bump_version`**
   - Bump semantic version (major, minor, or patch)
   - Params: `{type: "major"|"minor"|"patch", auto_push?: bool}`
   - Returns: `{success, old_version, new_version, tag, pushed, error?}`

2. **`changie_add_changelog`**
   - Add entry to changelog
   - Params: `{section: "added"|"changed"|"deprecated"|"removed"|"fixed"|"security", content: string}`
   - Returns: `{success, section, content, changelog_file, added, error?}`

3. **`changie_init`**
   - Initialize new changie project
   - Params: `{changelog_file?: string}`
   - Returns: `{success, changelog_file, created, error?}`

4. **`changie_get_version`**
   - Get current version from git tags
   - Params: `{}`
   - Returns: `{success, version, error?}`

### Running the MCP Server

**Option 1: Build from source**
```bash
go build -o changie-mcp-server ./cmd/mcp-server
./changie-mcp-server
```

**Option 2: Docker**
```bash
# Build image
docker build -f Dockerfile.mcp -t changie-mcp .

# Run server
docker run -i changie-mcp
```

**Option 3: With Claude Desktop**

Add to Claude Desktop config (`claude_desktop_config.json`):
```json
{
  "mcpServers": {
    "changie": {
      "command": "/path/to/changie-mcp-server"
    }
  }
}
```

### Architecture

- **Transport**: stdio (standard for MCP)
- **Protocol**: JSON-RPC 2.0
- **Implementation**: Go using official SDK
- **CLI Integration**: Calls `changie` binary with `--json` flag

### For Developers

**Source:**
- Server: `cmd/mcp-server/main.go`
- Tool handlers: `internal/mcp/tools.go`
- Dockerfile: `Dockerfile.mcp`

**Adding New Tools:**
1. Add handler function in `internal/mcp/tools.go`
2. Register tool in `cmd/mcp-server/main.go`
3. Ensure changie command supports `--json` output

### Context for MCP Integration

MCP servers and agents should:
1. Load `context.md` for tool context
2. Use prompts as tool documentation
3. Reference workflows for complex operations
4. Always check `success` field in responses
5. Use error messages for actionable hints

## Design Principles

1. **Self-Documenting** - Each file explains its purpose
2. **Machine-Readable** - Structured for parsing
3. **Human-Readable** - Clear markdown format
4. **Actionable** - Focus on what to do, not just what it is
5. **Error-Aware** - Include error handling and recovery

## Maintenance

When updating changie:
- Update `context.md` with new capabilities
- Add prompts for new common tasks
- Update workflows if processes change
- Keep examples current

## Related Resources

- `/llms.txt` - Comprehensive LLM-optimized documentation
- `/examples/` - Practical usage examples and scripts
- `/README.md` - Human-focused project documentation
- `/CLAUDE.md` - Architecture guide for contributors

## Contributing

When adding AI resources:
1. Keep content concise and actionable
2. Use consistent structure
3. Include examples
4. Focus on agent needs, not human needs
5. Test with actual AI agents if possible

## Future Enhancements

Potential additions:
- `prompts/hotfix-release.md` - Hotfix-specific workflow
- `prompts/rollback-release.md` - Release rollback procedures
- `workflows/monorepo-release.md` - Monorepo-specific workflows
- `context-extended.md` - Detailed technical internals
