package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	// OAuth2 redirect URI for local server
	redirectURI = "http://localhost:8080/callback"

	// Google Docs API scope for read-only access
	docsScope = "https://www.googleapis.com/auth/documents.readonly"
	// Google Drive API scope for read-only access (used for fetching comments)
	driveReadonlyScope = "https://www.googleapis.com/auth/drive.readonly"
)

// Authenticator handles OAuth2 authentication for Google Docs API.
type Authenticator struct {
	config    *oauth2.Config
	tokenPath string
}

// NewAuthenticator creates a new Authenticator by loading OAuth2 credentials from a file.
func NewAuthenticator(credPath string) (*Authenticator, error) {
	// Read credentials file
	credBytes, err := os.ReadFile(credPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %w\n\nTo fix this:\n1. Go to https://console.cloud.google.com/\n2. Create OAuth 2.0 credentials for a Desktop application\n3. Download the credentials JSON file\n4. Provide the path using --config flag", err)
	}

	// Parse credentials and create OAuth2 config
	config, err := google.ConfigFromJSON(credBytes, docsScope, driveReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %w", err)
	}

	// Set redirect URI for local server
	config.RedirectURL = redirectURI

	// Determine token path
	configDir, err := EnsureConfigDir()
	if err != nil {
		return nil, err
	}
	tokenPath := filepath.Join(configDir, "token.json")

	return &Authenticator{
		config:    config,
		tokenPath: tokenPath,
	}, nil
}

// GetClient returns an authenticated HTTP client.
// It first checks for a cached token. If none exists or if it's expired,
// it triggers the OAuth2 flow.
func (a *Authenticator) GetClient(ctx context.Context) (*http.Client, error) {
	// Try to load cached token
	token, err := LoadToken(a.tokenPath)
	if err != nil {
		// No cached token or error loading it - get new token
		log.Println("No cached token found. Starting OAuth2 flow...")
		token, err = a.getTokenFromWeb(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get token from web: %w", err)
		}

		// Save token for future use
		if err := SaveToken(a.tokenPath, token); err != nil {
			log.Printf("Warning: failed to save token: %v", err)
		} else {
			log.Printf("Token saved to %s", a.tokenPath)
		}
	}

	// Create HTTP client with token
	// The client will automatically refresh the token if it's expired
	tokenSource := a.config.TokenSource(ctx, token)
	return oauth2.NewClient(ctx, tokenSource), nil
}

// getTokenFromWeb starts a local HTTP server and initiates the OAuth2 flow.
func (a *Authenticator) getTokenFromWeb(ctx context.Context) (*oauth2.Token, error) {
	// Channel to receive the authorization code
	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	// Create HTTP server to handle OAuth callback
	server := &http.Server{Addr: "localhost:8080"}
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errChan <- fmt.Errorf("no authorization code received")
			fmt.Fprintf(w, "Error: No authorization code received")
			return
		}

		codeChan <- code
		fmt.Fprintf(w, "Authorization successful! You can close this window and return to the terminal.")
	})

	// Start server in background
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("failed to start callback server: %w", err)
		}
	}()

	// Generate authorization URL
	authURL := a.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	// Open browser
	fmt.Printf("Opening browser for authentication...\n")
	fmt.Printf("If the browser doesn't open automatically, visit this URL:\n%s\n\n", authURL)
	if err := openBrowser(authURL); err != nil {
		log.Printf("Failed to open browser: %v", err)
	}

	// Wait for authorization code or error
	var code string
	select {
	case code = <-codeChan:
		// Got the code, continue
	case err := <-errChan:
		server.Shutdown(ctx)
		return nil, err
	case <-ctx.Done():
		server.Shutdown(ctx)
		return nil, ctx.Err()
	}

	// Shutdown server
	server.Shutdown(ctx)

	// Exchange authorization code for token
	token, err := a.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange authorization code: %w", err)
	}

	return token, nil
}

// openBrowser opens the specified URL in the default browser.
func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}
