# Prompt: Add Changelog Entry

Use this prompt when you need to add an entry to the project changelog.

## Task

Add a properly categorized changelog entry to the CHANGELOG.md file.

## Steps

1. **Determine Entry Category**
   Based on the change type, select appropriate section:
   - `added` - New features
   - `changed` - Changes to existing functionality
   - `deprecated` - Soon-to-be removed features
   - `removed` - Removed features
   - `fixed` - Bug fixes
   - `security` - Security fixes

2. **Format Entry Text**
   - Write clear, user-focused description
   - Start with verb (Added, Fixed, Changed, etc.)
   - Be specific about what changed
   - Avoid technical jargon when possible

3. **Execute Command**
   ```bash
   changie changelog <section> "<entry text>"
   ```

4. **Verify Addition**
   - Confirm entry was added to Unreleased section
   - Check proper markdown formatting

## Decision Tree for Categories

```
What type of change?
├─ New functionality added → added
├─ Existing feature modified → changed
├─ Feature marked for removal → deprecated
├─ Feature completely removed → removed
├─ Bug or issue fixed → fixed
└─ Security vulnerability patched → security
```

## Examples

### Adding New Feature
```bash
changie changelog added "OAuth2 authentication support for enterprise users"
```

### Fixing Bug
```bash
changie changelog fixed "Memory leak in background job processor"
```

### Security Patch
```bash
changie changelog security "Updated dependencies to address CVE-2024-1234"
```

### Breaking Change
```bash
changie changelog changed "API endpoint /users now requires authentication (breaking change)"
```

## Best Practices

1. **Be User-Focused**
   - Bad: "Refactored auth module"
   - Good: "Improved login performance by 50%"

2. **Be Specific**
   - Bad: "Fixed bug"
   - Good: "Fixed crash when uploading files larger than 100MB"

3. **Include Context for Breaking Changes**
   - Always mention if it's a breaking change
   - Provide migration hints if applicable

4. **Group Related Changes**
   - Multiple related items can be separate entries
   - Or combine into one comprehensive entry

## Expected Output

```
Added to <Section> section: <entry text>
```

The entry will appear in CHANGELOG.md under:
```markdown
## [Unreleased]

### <Section>

- <entry text>
```

## Common Mistakes to Avoid

❌ Wrong section:
- "Added bug fix" - Should be `fixed`, not `added`

❌ Too technical:
- "Refactored UserService to use Repository pattern"
- Better: "Improved user data handling performance"

❌ Missing context:
- "Fixed issue"
- Better: "Fixed authentication timeout on slow networks"

## Example Usage

**User Request:** "I just added a dark mode feature, update the changelog"

**Agent Actions:**
1. Determine category: new feature → `added`
2. Format entry: "Dark mode support for better viewing in low-light environments"
3. Run: `changie changelog added "Dark mode support for better viewing in low-light environments"`
4. Confirm: "Added changelog entry under Added section"
