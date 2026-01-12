# Google Docs CLI

A command-line tool to fetch Google Docs content and convert it to Markdown with YAML frontmatter.

**Designed for AI Coding Agents**: This tool is specifically built to help AI coding agents fetch Google Docs documentation and convert it to clean markdown for analysis, context gathering, or integration into AI workflows. The stdout-based output and `--clean` flag make it ideal for piping into AI systems.

**📖 See [AGENTS.md](AGENTS.md)** for a comprehensive guide on integrating this tool with AI coding agents (Claude Code, Aider, Cursor, MCP servers, etc.).

## Features

- OAuth2 authentication with automatic token caching
- Converts Google Docs to clean Markdown format
- YAML frontmatter with document metadata
- Supports text formatting: bold, italic, strikethrough, links
- Supports document structure: headings, lists (bullet and numbered), tables
- Output to stdout for easy piping to files or other commands

## Prerequisites

- Go 1.24.1 or later
- A Google Cloud project with Google Docs API enabled
- OAuth 2.0 credentials (Desktop application type)

## Installation

### From Source

```bash
git clone https://github.com/famasya/gdocs-cli.git
cd gdocs-cli
go build -o gdocs-cli cmd/gdocs-cli/main.go
```

### Using Go Install

```bash
go install github.com/famasya/gdocs-cli/cmd/gdocs-cli@latest
```

## Google Cloud Setup

Before using this tool, you need to set up OAuth2 credentials:

### 1. Create a Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Navigate to "APIs & Services" > "Library"
4. Search for "Google Docs API" and enable it

### 2. Create OAuth 2.0 Credentials

1. Go to "APIs & Services" > "Credentials"
2. Click "Create Credentials" > "OAuth client ID"
3. If prompted, configure the OAuth consent screen:
   - Choose "External" user type
   - Fill in required fields (app name, user support email)
   - Add your email as a test user
   - Save and continue
4. Choose "Desktop application" as the application type
5. Give it a name (e.g., "gdocs-cli")
6. Click "Create"
7. Download the credentials JSON file
8. Save it as `credentials.json` (or any name you prefer)

**Important:** Keep this file secure and never commit it to version control.

## Usage

### Initialize Authentication (Recommended First Step)

Before using the CLI for the first time, set up your credentials and initialize OAuth authentication:

**Option 1: Use default config location (recommended for AI agents)**
```bash
# Copy your credentials to the default location
mkdir -p ~/.config/gdocs-cli
cp credentials.json ~/.config/gdocs-cli/config.json

# Initialize OAuth (will use default config automatically)
./gdocs-cli --init
```

**Option 2: Specify config path**
```bash
./gdocs-cli --init --config="./credentials.json"
```

This will:
1. Open your browser for Google OAuth consent
2. Ask you to authorize the application
3. Save the token to `~/.config/gdocs-cli/token.json`

After initialization, you can use the CLI without re-authenticating.

### Basic Usage

**Using default config location:**
```bash
./gdocs-cli --url="https://docs.google.com/document/d/YOUR_DOC_ID/edit"
```

**Specifying config path:**
```bash
./gdocs-cli --url="https://docs.google.com/document/d/YOUR_DOC_ID/edit" --config="./credentials.json"
```

The tool will automatically use the cached token - no browser interaction needed.

**Note**: If `--config` is not provided, the tool looks for credentials at `~/.config/gdocs-cli/config.json` by default. This makes it easy for AI agents to use the tool without specifying the config path every time.

### Output to File

```bash
./gdocs-cli --url="https://docs.google.com/document/d/YOUR_DOC_ID/edit" > output.md
```

### Piping to Other Commands

```bash
./gdocs-cli --url="..." | less
./gdocs-cli --url="..." | grep "keyword"
```

### Clean Output (Suppress Logs)

Use the `--clean` flag to suppress all log output and only show the markdown:

```bash
./gdocs-cli --url="..." --clean
```

This is useful when:
- Piping output to AI systems or other tools
- Saving to a file without log messages
- Using in scripts where only the markdown is needed

**Example:**
```bash
# Without --clean: shows logs like "Fetching document..." to stderr
./gdocs-cli --url="..." > output.md

# With --clean: only markdown to stdout, no logs
./gdocs-cli --url="..." --clean > output.md

# Perfect for AI agents: clean output piped to processing
./gdocs-cli --url="..." --clean | your-ai-tool
```

## Supported Google Docs Features

### Text Formatting
- **Bold** text
- *Italic* text
- ***Bold and italic***
- ~~Strikethrough~~
- [Links](https://example.com)

### Document Structure
- Headings (H1 through H6)
- Bullet lists
- Numbered lists
- Nested lists
- Paragraphs
- Tables

### YAML Frontmatter
The tool adds YAML frontmatter with document metadata:
```yaml
---
title: Document Title
author: (if available)
created: (if available)
modified: (if available)
---
```

**Note:** The Google Docs API v1 doesn't provide author or date information. These fields may be empty unless fetched from Google Drive API.

## Known Limitations

- **Tables:** Complex tables with merged cells may not convert perfectly to Markdown
- **Images:** Inline images are not currently supported
- **Drawings:** Not supported - will be skipped
- **Equations:** Not supported - will be skipped
- **Comments:** Not included in output (not in API response by default)
- **Metadata:** Author and dates require Google Drive API (not implemented in this version)

## Troubleshooting

### Error: Failed to read credentials file

**Cause:** The credentials file path is incorrect or the file doesn't exist.

**Solution:**
1. Verify the file path in the `--config` flag
2. Ensure you've downloaded the OAuth credentials JSON from Google Cloud Console
3. Use an absolute path or relative path from your current directory

### Error: Unable to access document

**Possible causes:**
1. The document is private and you don't have permission
2. The document doesn't exist
3. The document ID is incorrect

**Solutions:**
- Ensure the document is shared with your Google account
- Verify you're authenticated with the correct Google account
- Check that the URL is correct
- Try opening the document in your browser first

### Browser doesn't open during OAuth

**Solution:**
The authorization URL will be printed in the terminal. Copy and paste it into your browser manually.

### Token expired or invalid

**Solution:**
Delete the cached token and re-authenticate using the `--init` flag:
```bash
rm ~/.config/gdocs-cli/token.json
./gdocs-cli --init --config="credentials.json"
```

### Permission denied when creating config directory

**Solution:**
Ensure you have write permissions to `~/.config/`. Try creating it manually:
```bash
mkdir -p ~/.config/gdocs-cli
chmod 700 ~/.config/gdocs-cli
```

## Development

### Project Structure

```
gdocs-cli/
├── cmd/gdocs-cli/main.go              # CLI entry point
├── internal/
│   ├── auth/
│   │   ├── oauth.go                   # OAuth2 flow implementation
│   │   └── token.go                   # Token caching
│   ├── gdocs/
│   │   ├── client.go                  # Docs API client
│   │   └── url.go                     # URL parsing
│   └── markdown/
│       ├── converter.go               # Main converter
│       ├── text.go                    # Text formatting
│       ├── structure.go               # Structure conversion
│       └── frontmatter.go             # YAML frontmatter
├── go.mod
├── go.sum
└── README.md
```

### Building from Source

```bash
go build -o gdocs-cli cmd/gdocs-cli/main.go
```

### Running Tests

The project includes comprehensive unit tests for all core functionality:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run tests for a specific package
go test ./internal/gdocs -v
go test ./internal/markdown -v
go test ./internal/auth -v
```

**Test Coverage:**
- **CLI Integration** (`cmd/gdocs-cli/main_test.go`): End-to-end tests for CLI flags, error handling, and user flows
  - Help flag functionality
  - Missing required flags validation
  - Invalid URL handling
  - Missing credentials file errors
  - Clean flag recognition
- **URL Parsing** (`internal/gdocs/url_test.go`): Tests for extracting document IDs from various URL formats
- **Text Formatting** (`internal/markdown/text_test.go`): Tests for bold, italic, links, and text style conversion
- **Structure Conversion** (`internal/markdown/structure_test.go`): Tests for headings, lists, tables, and paragraph conversion
- **Token Handling** (`internal/auth/token_test.go`): Tests for token saving, loading, and file permissions

**Total: 45+ test cases** covering both unit and integration testing. All tests pass successfully and ensure the reliability of the CLI tool.

## Security Considerations

- **Credentials file:** Never commit your `credentials.json` to version control
- **Token cache:** Tokens are stored in `~/.config/gdocs-cli/token.json` with 0600 permissions (read/write for owner only)
- **OAuth scope:** The tool only requests `documents.readonly` scope - no write access
- **Config directory:** Created with 0700 permissions (accessible only by owner)

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Credits

Built with:
- [Google Docs API](https://developers.google.com/docs/api)
- [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2)
- [gopkg.in/yaml.v3](https://gopkg.in/yaml.v3)
