package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
)

// EnsureConfigDir creates the config directory if it doesn't exist.
// Creates ~/.config/gdocs-cli/ with 0700 permissions (user read/write/execute only).
func EnsureConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "gdocs-cli")

	// Create directory with 0700 permissions (rwx------)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create config directory %s: %w", configDir, err)
	}

	return configDir, nil
}

// LoadToken reads an OAuth2 token from a file.
func LoadToken(path string) (*oauth2.Token, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open token file: %w", err)
	}
	defer file.Close()

	token := &oauth2.Token{}
	if err := json.NewDecoder(file).Decode(token); err != nil {
		return nil, fmt.Errorf("failed to decode token: %w", err)
	}

	return token, nil
}

// SaveToken writes an OAuth2 token to a file with 0600 permissions.
// The file is created with read/write permissions for the owner only (rw-------).
func SaveToken(path string, token *oauth2.Token) error {
	// Create the file with 0600 permissions
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create token file: %w", err)
	}
	defer file.Close()

	// Encode token as JSON
	if err := json.NewEncoder(file).Encode(token); err != nil {
		return fmt.Errorf("failed to encode token: %w", err)
	}

	return nil
}
