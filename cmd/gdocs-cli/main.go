package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/famasya/gdocs-cli/internal/auth"
	"github.com/famasya/gdocs-cli/internal/gdocs"
	"github.com/famasya/gdocs-cli/internal/markdown"
)

//go:embed instruction.txt
var instructionText string

func main() {
	// Define flags
	urlFlag := flag.String("url", "", "Google Docs URL (required for normal operation)")
	configFlag := flag.String("config", "", "Path to OAuth credentials JSON file (defaults to ~/.config/gdocs-cli/config.json)")
	initFlag := flag.Bool("init", false, "Initialize OAuth and save token to default location")
	cleanFlag := flag.Bool("clean", false, "Clean output (suppress all logs, only output markdown)")
	instructionFlag := flag.Bool("instruction", false, "Print integration instructions for AI coding agents")
	flag.Parse()

	// Handle instruction mode - print instructions and exit
	if *instructionFlag {
		fmt.Print(instructionText)
		return
	}

	// Handle clean mode - suppress all logs
	if *cleanFlag {
		log.SetOutput(io.Discard)
	}

	// Determine config path (use default if not specified)
	configPath := *configFlag
	if configPath == "" {
		defaultPath, err := getDefaultConfigPath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		configPath = defaultPath
	}

	// Handle init mode
	if *initFlag {
		if err := initAuth(configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Validate flags for normal operation
	if *urlFlag == "" {
		fmt.Fprintln(os.Stderr, "Error: --url flag is required")
		fmt.Fprintln(os.Stderr)
		flag.Usage()
		os.Exit(1)
	}

	// Run the main logic
	if err := run(*urlFlag, configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(docURL, credPath string) error {
	ctx := context.Background()

	// Extract document ID from URL
	docID, err := gdocs.ExtractDocumentID(docURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Extract tab ID from URL (may be empty)
	tabID := gdocs.ExtractTabID(docURL)

	// Create authenticator
	authenticator, err := auth.NewAuthenticator(credPath)
	if err != nil {
		return fmt.Errorf("authentication setup failed: %w", err)
	}

	// Get authenticated HTTP client
	httpClient, err := authenticator.GetClient(ctx)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create Google Docs API client
	client, err := gdocs.NewClient(ctx, httpClient)
	if err != nil {
		return fmt.Errorf("failed to create Docs client: %w", err)
	}

	// Fetch document
	log.Printf("Fetching document %s...", docID)
	doc, err := client.FetchDocument(docID)
	if err != nil {
		return fmt.Errorf("failed to fetch document: %w", err)
	}

	// Convert to markdown
	var converter *markdown.Converter
	if tabID != "" {
		// Find the specific tab
		tab := gdocs.FindTab(doc, tabID)
		if tab == nil {
			return fmt.Errorf("tab '%s' not found in document", tabID)
		}
		log.Printf("Using tab: %s", tab.TabProperties.Title)
		converter = markdown.NewConverterFromTab(doc, tab)
	} else {
		converter = markdown.NewConverter(doc)
	}

	markdownOutput, err := converter.Convert()
	if err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	// Print to stdout
	fmt.Print(markdownOutput)

	return nil
}

// initAuth initializes OAuth authentication and saves the token.
func initAuth(credPath string) error {
	ctx := context.Background()

	fmt.Println("Initializing OAuth authentication...")
	fmt.Println()

	// Create authenticator
	authenticator, err := auth.NewAuthenticator(credPath)
	if err != nil {
		return fmt.Errorf("authentication setup failed: %w", err)
	}

	// Get authenticated HTTP client (this will trigger OAuth flow)
	_, err = authenticator.GetClient(ctx)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	fmt.Println()
	fmt.Println("✓ Authentication successful!")
	fmt.Println("✓ Token saved to ~/.config/gdocs-cli/token.json")
	fmt.Println()
	fmt.Println("You can now use the CLI without the --init flag:")
	fmt.Println("  ./gdocs-cli --url=\"https://docs.google.com/document/d/DOC_ID/edit\"")

	return nil
}

// getDefaultConfigPath returns the default path for the config file.
func getDefaultConfigPath() (string, error) {
	configDir, err := auth.EnsureConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %w", err)
	}
	return configDir + "/config.json", nil
}
