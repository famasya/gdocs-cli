# Using gdocs-cli with AI Coding Agents

This guide explains how to integrate `gdocs-cli` into AI coding agent workflows to fetch and process Google Docs documentation.

## Why This Tool for AI Agents?

This CLI is specifically designed for AI coding agents that need to:
- Fetch technical documentation from Google Docs
- Convert documentation to markdown for analysis
- Integrate Google Docs content into their knowledge base
- Process documentation as part of their coding workflow

## Quick Setup for AI Agents

### 1. One-Time Setup

```bash
# Install the tool
go install github.com/famasya/gdocs-cli/cmd/gdocs-cli@latest

# Or build from source
git clone https://github.com/famasya/gdocs-cli.git
cd gdocs-cli
go build -o gdocs-cli cmd/gdocs-cli/main.go
sudo mv gdocs-cli /usr/local/bin/  # Make it globally available

# Set up credentials in default location
mkdir -p ~/.config/gdocs-cli
cp your-oauth-credentials.json ~/.config/gdocs-cli/config.json

# Initialize OAuth authentication (opens browser once)
gdocs-cli --init
```

After this one-time setup, the agent can use the tool without any authentication steps.

### 2. Basic Usage in Agent Context

Once configured, agents can simply run:

```bash
gdocs-cli --url="https://docs.google.com/document/d/DOC_ID/edit" --clean
```

This outputs clean markdown to stdout with:
- No log messages (thanks to `--clean`)
- YAML frontmatter with document metadata
- Properly formatted markdown content
- All text formatting preserved (bold, italic, links, etc.)

## Integration Patterns

### Pattern 1: Direct Context Injection

Fetch documentation and use it directly in the agent's context:

```bash
# Fetch and store in variable
DOCS=$(gdocs-cli --url="https://docs.google.com/document/d/DOC_ID/edit" --clean)

# Use in agent prompt
echo "Here's the documentation: $DOCS"
```

### Pattern 2: Save to File for Processing

```bash
# Fetch and save to file
gdocs-cli --url="https://docs.google.com/document/d/DOC_ID/edit" --clean > docs/api-spec.md

# Agent can now read the markdown file
```

### Pattern 3: Multiple Documents

```bash
# Fetch multiple related documents
gdocs-cli --url="https://docs.google.com/document/d/DOC1/edit" --clean > architecture.md
gdocs-cli --url="https://docs.google.com/document/d/DOC2/edit" --clean > api-reference.md
gdocs-cli --url="https://docs.google.com/document/d/DOC3/edit" --clean > deployment-guide.md

# Agent now has complete documentation set
```

### Pattern 4: Piped Processing

```bash
# Fetch and pipe to other tools
gdocs-cli --url="..." --clean | grep "TODO" | wc -l  # Count TODOs
gdocs-cli --url="..." --clean | head -50            # Get first 50 lines
gdocs-cli --url="..." --clean | pandoc -f markdown -t html  # Convert to HTML
```

## MCP Server Integration

For Model Context Protocol (MCP) servers, this tool can be wrapped as a resource provider:

```json
{
  "mcpServers": {
    "gdocs": {
      "command": "gdocs-cli",
      "args": ["--url", "${url}", "--clean"],
      "env": {
        "HOME": "${HOME}"
      }
    }
  }
}
```

The MCP server can then expose Google Docs as readable resources to the agent.

## Claude Code Integration

### As a Tool in Claude Code Context

When working with Claude Code (claude.ai/code), you can add this to your project's context:

**In `.claudeignore` or project instructions:**
```
This project uses gdocs-cli to fetch documentation from Google Docs.

To fetch documentation:
gdocs-cli --url="<google-docs-url>" --clean

Common documentation sources:
- Architecture: https://docs.google.com/document/d/DOC1/edit
- API Specs: https://docs.google.com/document/d/DOC2/edit
- Deployment: https://docs.google.com/document/d/DOC3/edit
```

### Example Claude Code Workflow

```bash
# 1. Fetch latest documentation
gdocs-cli --url="https://docs.google.com/document/d/API-SPEC/edit" --clean > api-spec.md

# 2. Ask Claude Code to analyze
# "Read api-spec.md and implement the POST /users endpoint according to the spec"

# 3. Update docs when implementation is done
# "Update the Google Doc at <url> with the new endpoint details"
# (Note: This tool is read-only, manual update needed)
```

## Aider Integration

For [Aider](https://github.com/paul-gauthier/aider) users:

```bash
# Start Aider with fetched documentation
gdocs-cli --url="..." --clean > docs.md
aider --read docs.md main.go

# Or inline
aider --read <(gdocs-cli --url="..." --clean) main.go
```

## Cursor / Windsurf Integration

Add to `.cursorrules` or `.windsurfrules`:

```markdown
# Documentation Sources

This project's documentation is stored in Google Docs. Fetch using:

```bash
gdocs-cli --url="<doc-url>" --clean
```

Key documents:
- Architecture: [DOC_ID_1]
- API Reference: [DOC_ID_2]
```

## Environment Variables (Optional)

For advanced setups, you can wrap the tool with environment variables:

```bash
# Create a wrapper script: ~/bin/fetch-docs
#!/bin/bash
export GDOCS_DEFAULT_URL="https://docs.google.com/document/d/YOUR-MAIN-DOC/edit"

if [ -z "$1" ]; then
  gdocs-cli --url="$GDOCS_DEFAULT_URL" --clean
else
  gdocs-cli --url="$1" --clean
fi
```

Then agents can simply run:
```bash
fetch-docs  # Uses default doc
fetch-docs "https://docs.google.com/document/d/OTHER-DOC/edit"  # Custom doc
```

## Troubleshooting for Agents

### Token Expired
If authentication fails, re-run:
```bash
gdocs-cli --init
```

### Permission Denied
Ensure the Google Doc is shared with the Google account used for OAuth authentication.

### Rate Limiting
Google Docs API has rate limits. For high-volume usage, implement caching:

```bash
# Cache docs with timestamp
CACHE_DIR="~/.cache/gdocs"
DOC_ID="extracted-from-url"
CACHE_FILE="$CACHE_DIR/$DOC_ID.md"
CACHE_TIME=3600  # 1 hour

if [ ! -f "$CACHE_FILE" ] || [ $(find "$CACHE_FILE" -mmin +60) ]; then
  gdocs-cli --url="..." --clean > "$CACHE_FILE"
fi

cat "$CACHE_FILE"
```

## Best Practices

1. **Use `--clean` flag always** - Ensures only markdown output for easy parsing
2. **Cache frequently accessed docs** - Avoid hitting API rate limits
3. **Set up default config** - Use `~/.config/gdocs-cli/config.json` for seamless access
4. **Document your doc URLs** - Keep a list of important documentation sources
5. **Version control output** - Commit fetched markdown files to track doc changes over time
6. **Handle errors gracefully** - Check exit codes and stderr for error messages

## Output Format

The tool outputs markdown with YAML frontmatter:

```markdown
---
title: Document Title
author:
created:
modified:
---

# Main Heading

Content with **bold**, *italic*, and [links](https://example.com).

- Bullet lists
- Nested lists
  - Sub-items

## Tables

| Header 1 | Header 2 |
|----------|----------|
| Cell 1   | Cell 2   |
```

This format is easily parseable by AI agents and integrates well with existing markdown tooling.

## Security Considerations

- **Credentials**: The OAuth credentials file contains sensitive information. Ensure it's properly secured.
- **Token**: The cached token at `~/.config/gdocs-cli/token.json` provides access to Google Docs. Protect it with file permissions (0600).
- **Scope**: This tool only requests `documents.readonly` scope - no write access.
- **Sharing**: Ensure agents only access documents they should have permission to read.

## Support

For issues or questions:
- GitHub Issues: https://github.com/famasya/gdocs-cli/issues
- Documentation: See README.md and CLAUDE.md in the repository
