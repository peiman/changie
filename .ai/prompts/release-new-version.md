# Prompt: Release New Version

Use this prompt when you need to release a new version of a project.

## Task

Release a new version of the project using changie, ensuring all prerequisites are met and the release is successful.

## Steps

1. **Verify Prerequisites**
   - Check git is installed and repository is initialized
   - Verify working directory is clean (no uncommitted changes)
   - Confirm on main or master branch (or use --allow-any-branch if intended)
   - Ensure CHANGELOG.md exists with unreleased entries

2. **Determine Bump Type**
   - Analyze changes in unreleased section of CHANGELOG.md
   - Determine if changes are:
     - Breaking/incompatible → major
     - New features/enhancements → minor
     - Bug fixes/patches → patch

3. **Execute Release**
   ```bash
   changie bump <type> --json
   ```

4. **Verify Success**
   - Parse JSON output
   - Check `success: true`
   - Extract `new_version` value
   - Confirm git tag was created

5. **Push Changes**
   ```bash
   git push origin main
   git push origin --tags
   ```

   Or use `--auto-push` flag in step 3

## Expected Input

- Clean git working directory
- CHANGELOG.md with unreleased entries
- Confirmation of bump type (major/minor/patch)

## Expected Output

```json
{
  "success": true,
  "old_version": "1.2.3",
  "new_version": "1.3.0",
  "tag": "v1.3.0",
  "changelog_file": "CHANGELOG.md",
  "pushed": false,
  "bump_type": "minor"
}
```

## Error Handling

If any prerequisite fails:
- Report specific issue
- Provide solution (e.g., "Run git commit first")
- Abort release until resolved

## Success Criteria

- ✅ New version tagged in git
- ✅ CHANGELOG.md updated (Unreleased → versioned section)
- ✅ Changelog committed with "Release vX.Y.Z" message
- ✅ Changes optionally pushed to remote

## Example Usage

**User Request:** "Release a new minor version"

**Agent Actions:**
1. Run `git status` → verify clean
2. Run `git branch --show-current` → verify main/master
3. Check CHANGELOG.md has unreleased entries
4. Run `changie bump minor --json`
5. Parse output, confirm success
6. Report: "Successfully released v1.3.0"
