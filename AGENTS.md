# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a CLI tool that fetches Google Docs via the Google Docs API and converts them to Markdown with YAML frontmatter. It uses OAuth2 authentication with token caching for user authorization.

**Primary Use Case**: This tool is designed for AI coding agents to fetch Google Docs documentation and convert it to markdown for analysis, context gathering, or integration into AI workflows. The clean markdown output and stdout design make it ideal for piping into AI systems.

**📖 For comprehensive integration guide, see [INSTRUCTION.md](INSTRUCTION.md)** - This file contains detailed instructions for integrating this tool with various AI coding agents (Claude Code, Aider, Cursor, MCP servers, etc.).

## Build and Test Commands

```bash
# Build the CLI binary
go build -o gdocs-cli cmd/gdocs-cli/main.go

# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run tests with verbose output
go test ./... -v

# Run specific package tests
go test ./internal/gdocs -v
go test ./internal/markdown -v
go test ./internal/auth -v
go test ./cmd/gdocs-cli -v

# Run a single test
go test ./internal/gdocs -run TestExtractDocumentID -v

# Clean up dependencies
go mod tidy
```

## Architecture

### Data Flow

```
User Input (URL + credentials)
    ↓
main.go (CLI parsing & orchestration)
    ↓
auth/oauth.go (OAuth2 flow & token caching)
    ↓
gdocs/client.go (Google Docs API client)
    ↓
gdocs/url.go (Extract doc ID from URL)
    ↓
client.FetchDocument() → Returns *docs.Document
    ↓
markdown/converter.go (Main orchestrator)
    ↓
    ├→ markdown/frontmatter.go (YAML metadata)
    ├→ markdown/structure.go (headings, lists, tables)
    └→ markdown/text.go (bold, italic, links)
    ↓
Markdown output to stdout
```

### Key Architectural Decisions

1. **OAuth2 Token Caching**: Tokens are cached at `~/.config/gdocs-cli/token.json` to avoid repeated browser authentication. The `auth` package handles both initial authentication and automatic token refresh.

2. **Structural Element Processing**: Google Docs API returns a tree of `StructuralElement` objects (paragraphs, tables, etc.). The converter iterates through `doc.Body.Content` and dispatches to specialized functions based on element type.

3. **Text Style Application**: Text formatting is applied at the `TextRun` level. Each paragraph contains multiple `ParagraphElement` objects, which may contain `TextRun` objects with `TextStyle` properties (bold, italic, link, etc.).

4. **List Handling**: Lists in Google Docs API are complex - each list item is a paragraph with a `Bullet` field containing `ListId` and `NestingLevel`. The converter uses indentation (2 spaces per level) to represent nesting.

5. **Clean Flag**: The `--clean` flag suppresses logs by setting `log.SetOutput(io.Discard)`. Errors still go to stderr, but informational logs are silenced.

### Package Responsibilities

- **`cmd/gdocs-cli`**: CLI entry point, flag parsing, orchestration. No business logic here.
- **`internal/auth`**: OAuth2 flow, token caching, credential loading. Handles all authentication concerns.
- **`internal/gdocs`**: Google Docs API client wrapper and URL parsing. Fetches documents from the API.
- **`internal/markdown`**: Conversion logic split into:
  - `converter.go`: Main orchestrator that drives the conversion
  - `frontmatter.go`: YAML frontmatter generation
  - `text.go`: Text-level formatting (bold, italic, links)
  - `structure.go`: Document structure (headings, lists, tables, paragraphs)

### Google Docs API Structure

Understanding the Google Docs API structure is critical:

```
docs.Document
  ├─ Title (string)
  ├─ Body
  │   └─ Content []StructuralElement
  │       ├─ Paragraph
  │       │   ├─ Elements []ParagraphElement
  │       │   │   └─ TextRun
  │       │   │       ├─ Content (string)
  │       │   │       └─ TextStyle (Bold, Italic, Link, etc.)
  │       │   ├─ ParagraphStyle (NamedStyleType: HEADING_1, NORMAL_TEXT, etc.)
  │       │   └─ Bullet (ListId, NestingLevel)
  │       └─ Table
  │           └─ TableRows []TableRow
  │               └─ TableCells []TableCell
  │                   └─ Content []StructuralElement (recursive)
  └─ Lists map[string]*List (list metadata by ListId)
```

When adding support for new Google Docs features:
1. Check the `StructuralElement` type in `converter.go`
2. Add a new handler function in `structure.go` if it's structural
3. Add formatting logic in `text.go` if it's text-level
4. Update tests to cover the new feature

### Testing Strategy

- **Unit tests** (`*_test.go` in each package): Test individual functions with mocked Google Docs API structures
- **Integration tests** (`cmd/gdocs-cli/main_test.go`): Build and execute the binary as a subprocess to test end-to-end flows
- Integration tests don't show in coverage reports because they run the binary externally, but they ensure the CLI works correctly from a user perspective

### CLI Flags

- `--url`: Google Docs URL (required for normal operation)
- `--config`: Path to OAuth2 credentials JSON (defaults to `~/.config/gdocs-cli/config.json`)
- `--init`: Initialize OAuth and save token (doesn't require --url)
- `--clean`: Suppress all logs, only output markdown
- `--instruction`: Print integration instructions for AI coding agents

**Default Config Path**: If `--config` is not provided, the tool automatically looks for credentials at `~/.config/gdocs-cli/config.json`. This is particularly useful for AI agents that can set up credentials once and then use the tool without specifying the path on every invocation.

The `--clean` flag is particularly important for AI agents and scripting - it ensures only markdown goes to stdout, making it easy to pipe the output directly into AI systems or other tools.

### Token Refresh

OAuth2 tokens expire. The `auth` package uses `oauth2.TokenSource` which automatically refreshes tokens when they expire. The refreshed token is saved back to the cache file. If refresh fails (e.g., token revoked), the user must re-authenticate with `--init`.

### Error Handling

Errors are wrapped with context using `fmt.Errorf("context: %w", err)` to provide clear error chains. The main function prints errors to stderr and exits with non-zero codes. When adding new features, maintain this pattern for user-friendly error messages.
