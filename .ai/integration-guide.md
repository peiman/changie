# Changie MCP Server - Integration Guide

This guide shows how to integrate changie's MCP server with popular AI coding tools.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Claude Desktop](#claude-desktop)
- [Claude Code CLI](#claude-code-cli)
- [Cursor IDE](#cursor-ide)
- [Verification](#verification)
- [Usage Examples](#usage-examples)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

### 1. Build the MCP Server

**Option A: Build from source**
```bash
# Clone changie repository
git clone https://github.com/peiman/changie.git
cd changie

# Build MCP server binary
task build:mcp

# Note the full path to the binary
CHANGIE_MCP_PATH=$(pwd)/changie-mcp-server
echo "MCP server path: $CHANGIE_MCP_PATH"
```

**Option B: Use Docker (Recommended)**
```bash
# Build Docker image
task docker:build:mcp

# Image will be tagged as changie-mcp:latest
```

### 2. Ensure changie CLI is installed

The MCP server calls the `changie` binary with `--json` flag, so it must be in your PATH:

```bash
# Install changie
go install github.com/peiman/changie@latest

# Verify installation
changie --version
```

---

## Claude Desktop

Claude Desktop has **native MCP support** and is the easiest to configure.

### Method 1: Manual Configuration (Recommended)

**1. Open Claude Desktop configuration:**

**macOS:**
```bash
# Config location
open ~/Library/Application\ Support/Claude/claude_desktop_config.json
```

**Windows:**
```bash
# Config location
%APPDATA%\Claude\claude_desktop_config.json
```

**Linux:**
```bash
# Config location
~/.config/Claude/claude_desktop_config.json
```

**2. Add changie MCP server to config:**

```json
{
  "mcpServers": {
    "changie": {
      "command": "/absolute/path/to/changie-mcp-server",
      "args": []
    }
  }
}
```

**Example (macOS):**
```json
{
  "mcpServers": {
    "changie": {
      "command": "/Users/yourname/go/bin/changie-mcp-server",
      "args": []
    }
  }
}
```

**3. Restart Claude Desktop**

Close and reopen Claude Desktop completely.

**4. Verify:**
- Look for a 🔨 hammer icon in the bottom-right corner of the chat input
- Click it to see available tools
- You should see: `changie_bump_version`, `changie_add_changelog`, `changie_init`, `changie_get_version`

### Method 2: Using Docker

If you prefer Docker:

```json
{
  "mcpServers": {
    "changie": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "--network=host",
        "-v", "/path/to/your/project:/workspace",
        "-w", "/workspace",
        "changie-mcp:latest"
      ]
    }
  }
}
```

**Note:** Replace `/path/to/your/project` with your actual project directory.

---

## Claude Code CLI

Claude Code CLI has **built-in MCP support** with a wizard.

### Method 1: CLI Wizard (Easiest)

Unfortunately, changie is not yet in the official MCP registry, so we'll use manual configuration.

### Method 2: Manual Configuration (Current Method)

**1. Locate config file:**

```bash
# macOS/Linux
~/.config/Claude/claude_code_config.json

# Or check with:
claude mcp list
# This will show config file location
```

**2. Edit config file:**

```bash
# Open in your editor
code ~/.config/Claude/claude_code_config.json
# or
vim ~/.config/Claude/claude_code_config.json
```

**3. Add changie server:**

```json
{
  "mcpServers": {
    "changie": {
      "command": "/absolute/path/to/changie-mcp-server",
      "args": [],
      "scope": "user"
    }
  }
}
```

**4. Verify:**

```bash
# List configured MCP servers
claude mcp list

# Should show changie in the list
```

**5. Restart Claude Code:**

```bash
# Exit and restart your terminal or
# Kill and restart Claude Code process
```

### Method 3: Project-Scoped Configuration

For project-specific setup:

**1. Create `.claude/` directory in your project:**

```bash
cd /path/to/your/project
mkdir -p .claude
```

**2. Create `.claude/mcp_config.json`:**

```json
{
  "mcpServers": {
    "changie": {
      "command": "/absolute/path/to/changie-mcp-server",
      "args": []
    }
  }
}
```

**3. Claude Code will auto-detect this configuration when working in the project.**

---

## Cursor IDE

Cursor has **full MCP support** since early 2025.

### Method 1: Project-Scoped Configuration (Recommended)

**1. Create `.cursor/` directory in your project:**

```bash
cd /path/to/your/project
mkdir -p .cursor
```

**2. Create `.cursor/mcp.json`:**

```json
{
  "mcpServers": {
    "changie": {
      "command": "/absolute/path/to/changie-mcp-server",
      "args": [],
      "transport": "stdio"
    }
  }
}
```

**Example:**
```json
{
  "mcpServers": {
    "changie": {
      "command": "/Users/yourname/dev/changie/changie-mcp-server",
      "args": [],
      "transport": "stdio"
    }
  }
}
```

**3. Restart Cursor IDE**

Cursor will auto-detect the configuration.

### Method 2: Global Configuration

**1. Locate Cursor's global config:**

```bash
# macOS
~/Library/Application Support/Cursor/mcp.json

# Linux
~/.config/Cursor/mcp.json

# Windows
%APPDATA%\Cursor\mcp.json
```

**2. Add changie server (same format as project-scoped)**

### Method 3: Using Docker with Cursor

```json
{
  "mcpServers": {
    "changie": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "-v", "${workspaceFolder}:/workspace",
        "-w", "/workspace",
        "changie-mcp:latest"
      ],
      "transport": "stdio"
    }
  }
}
```

**Note:** `${workspaceFolder}` is a Cursor variable that expands to your current project path.

---

## Verification

After configuration, verify the MCP server is working:

### Claude Desktop

1. Open a new conversation
2. Look for 🔨 hammer icon (bottom-right of input)
3. Click it to see tools
4. Type: "What tools do you have access to?"
5. Claude should mention changie tools

### Claude Code CLI

```bash
# List MCP servers
claude mcp list

# Should show:
# changie - /path/to/changie-mcp-server
```

In a conversation:
```bash
claude "What MCP tools do you have?"
```

### Cursor IDE

1. Open Cursor
2. Open AI chat (Cmd+L or Ctrl+L)
3. Ask: "What tools do you have available?"
4. Should list changie tools

### Test the Integration

Ask the AI assistant:

```
"Use the changie tools to get the current version of this project"
```

The AI should use `changie_get_version` tool and return your project's version.

---

## Usage Examples

Once configured, you can ask your AI assistant to perform changelog operations:

### Example 1: Bump Version

**You:**
```
Please bump the patch version and push the changes
```

**AI will:**
1. Call `changie_bump_version` with type="patch", auto_push=true
2. Return the new version number
3. Confirm the changes were pushed

### Example 2: Add Changelog Entry

**You:**
```
Add a changelog entry under "fixed" that says "Resolved memory leak in parser"
```

**AI will:**
1. Call `changie_add_changelog` with section="fixed", content="Resolved memory leak in parser"
2. Confirm the entry was added to CHANGELOG.md

### Example 3: Initialize Project

**You:**
```
Initialize changie in this project with a custom changelog file named HISTORY.md
```

**AI will:**
1. Call `changie_init` with changelog_file="HISTORY.md"
2. Create HISTORY.md with proper structure
3. Confirm initialization

### Example 4: Check Current Version

**You:**
```
What's the current version of this project?
```

**AI will:**
1. Call `changie_get_version`
2. Return version from git tags

### Example 5: Complex Workflow

**You:**
```
I just fixed a critical bug. Please:
1. Add it to the changelog under "fixed"
2. Bump the patch version
3. Push the changes
```

**AI will:**
1. Call `changie_add_changelog`
2. Call `changie_bump_version` with auto_push=true
3. Report the new version number

---

## Troubleshooting

### "MCP server not found" or "Command not found"

**Problem:** Path to `changie-mcp-server` is incorrect

**Solution:**
```bash
# Find the binary
which changie-mcp-server

# Or if built locally
find ~ -name "changie-mcp-server" -type f 2>/dev/null

# Use the absolute path in config
```

### "changie: command not found"

**Problem:** The `changie` CLI is not in PATH

**Solution:**
```bash
# Install changie
go install github.com/peiman/changie@latest

# Verify
changie --version

# If still not found, add Go bin to PATH
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

### MCP server starts but tools don't work

**Problem:** Working directory or git repository issues

**Solution:**
- Ensure you're running AI assistant in a git repository
- Check that CHANGELOG.md exists (or run `changie init` first)
- For bump commands, ensure working directory is clean

### "Transport error" or "Connection refused"

**Problem:** MCP server crashed or isn't using stdio correctly

**Solution:**
```bash
# Test MCP server manually
echo '{"jsonrpc":"2.0","method":"initialize","params":{},"id":1}' | /path/to/changie-mcp-server

# Should return JSON response
# If it crashes, check logs
```

### Docker version not working

**Problem:** Volume mounts or network issues

**Solution:**
```bash
# Test Docker image manually
docker run --rm -i -v "$(pwd):/workspace" -w /workspace changie-mcp:latest

# Verify image exists
docker images | grep changie-mcp

# Rebuild if needed
task docker:build:mcp
```

### Changes not detected after editing config

**Solution:**
- **Claude Desktop:** Completely quit (Cmd+Q / Alt+F4) and restart
- **Claude Code:** Restart terminal session
- **Cursor:** Restart IDE (Cmd+Q / Ctrl+Q, then reopen)

### Tools show up but commands fail

**Problem:** Project setup issues

**Checklist:**
```bash
# 1. Is git initialized?
git status

# 2. Is changie initialized?
ls -la CHANGELOG.md

# 3. Are there git tags?
git tag -l

# 4. Is working directory clean?
git status --short

# 5. Can changie run with --json?
changie bump patch --json --allow-any-branch
```

---

## Advanced Configuration

### Environment Variables

Pass environment variables to MCP server:

```json
{
  "mcpServers": {
    "changie": {
      "command": "/path/to/changie-mcp-server",
      "args": [],
      "env": {
        "APP_LOG_LEVEL": "debug",
        "APP_CHANGELOG_FILE": "HISTORY.md"
      }
    }
  }
}
```

### Multiple Projects with Different Configs

Use project-scoped configs:

```
my-project-1/
  .cursor/mcp.json          # Cursor config for this project
  .claude/mcp_config.json   # Claude Code config for this project

my-project-2/
  .cursor/mcp.json          # Different config here
  .claude/mcp_config.json
```

### Debugging MCP Communication

Enable debug logging:

**Claude Desktop:** Check logs:
```bash
# macOS
tail -f ~/Library/Logs/Claude/mcp*.log

# Windows
# Check %APPDATA%\Claude\logs\
```

**Test manually:**
```bash
# Send test request to MCP server
echo '{"jsonrpc":"2.0","method":"tools/list","params":{},"id":1}' | /path/to/changie-mcp-server
```

---

## Next Steps

1. **Read `.ai/context.md`** - Understand changie capabilities and constraints
2. **Try examples** - Test with simple commands first
3. **Check `examples/`** - See real-world usage patterns
4. **Read `llms.txt`** - Reference for AI assistants

---

## Support

- **Documentation:** See `.ai/README.md` for MCP architecture details
- **Issues:** https://github.com/peiman/changie/issues
- **MCP Docs:** https://modelcontextprotocol.io/

---

**Last Updated:** 2025-10-06
**Compatible With:** Claude Desktop, Claude Code CLI, Cursor IDE, and any MCP-compatible client
