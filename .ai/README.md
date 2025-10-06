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

This directory structure is designed to support MCP integration:

### MCP Tool Definitions

Based on the content here, changie can be exposed as MCP tools:

**Tool: `changie_init`**
- Creates CHANGELOG.md
- Returns: `{created: bool, file: string}`

**Tool: `changie_add_entry`**
- Params: `{section: enum, content: string}`
- Returns: `{success: bool, section: string, content: string}`

**Tool: `changie_bump_version`**
- Params: `{type: enum[major,minor,patch], auto_push: bool}`
- Returns: `{success: bool, old_version: string, new_version: string, error?: string}`

### Context for MCP Servers

MCP servers should:
1. Load `context.md` as tool context
2. Use prompts as tool documentation
3. Reference workflows for complex operations
4. Parse JSON output from changie commands

### Example MCP Server Usage

```typescript
// MCP server using changie
import { Tool } from '@modelcontextprotocol/sdk';

const changieTools: Tool[] = [
  {
    name: 'changie_bump_version',
    description: 'Release a new version (read from .ai/context.md)',
    inputSchema: {
      type: 'object',
      properties: {
        bump_type: {
          type: 'string',
          enum: ['major', 'minor', 'patch'],
          description: 'major=breaking, minor=features, patch=fixes'
        },
        auto_push: { type: 'boolean', default: false }
      }
    },
    handler: async (input) => {
      const result = await exec(
        `changie bump ${input.bump_type} ${input.auto_push ? '--auto-push' : ''} --json`
      );
      return JSON.parse(result);
    }
  }
];
```

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
