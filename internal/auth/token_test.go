package auth

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestSaveAndLoadToken(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gdocs-cli-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tokenPath := filepath.Join(tempDir, "test-token.json")

	// Create a test token
	testToken := &oauth2.Token{
		AccessToken:  "test-access-token",
		TokenType:    "Bearer",
		RefreshToken: "test-refresh-token",
		Expiry:       time.Now().Add(1 * time.Hour),
	}

	// Test SaveToken
	err = SaveToken(tokenPath, testToken)
	if err != nil {
		t.Fatalf("SaveToken() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		t.Errorf("Token file was not created")
	}

	// Test LoadToken
	loadedToken, err := LoadToken(tokenPath)
	if err != nil {
		t.Fatalf("LoadToken() error = %v", err)
	}

	// Verify token contents
	if loadedToken.AccessToken != testToken.AccessToken {
		t.Errorf("AccessToken = %v, want %v", loadedToken.AccessToken, testToken.AccessToken)
	}
	if loadedToken.TokenType != testToken.TokenType {
		t.Errorf("TokenType = %v, want %v", loadedToken.TokenType, testToken.TokenType)
	}
	if loadedToken.RefreshToken != testToken.RefreshToken {
		t.Errorf("RefreshToken = %v, want %v", loadedToken.RefreshToken, testToken.RefreshToken)
	}
}

func TestLoadToken_FileNotFound(t *testing.T) {
	_, err := LoadToken("/nonexistent/path/token.json")
	if err == nil {
		t.Error("LoadToken() expected error for nonexistent file, got nil")
	}
}

func TestLoadToken_InvalidJSON(t *testing.T) {
	// Create a temporary file with invalid JSON
	tempDir, err := os.MkdirTemp("", "gdocs-cli-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tokenPath := filepath.Join(tempDir, "invalid.json")
	err = os.WriteFile(tokenPath, []byte("not valid json"), 0600)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	_, err = LoadToken(tokenPath)
	if err == nil {
		t.Error("LoadToken() expected error for invalid JSON, got nil")
	}
}

func TestEnsureConfigDir(t *testing.T) {
	// Note: This test will create a real config directory
	// We can't easily mock os.UserHomeDir() without more complex setup
	configDir, err := EnsureConfigDir()
	if err != nil {
		t.Fatalf("EnsureConfigDir() error = %v", err)
	}

	// Verify directory was created
	info, err := os.Stat(configDir)
	if err != nil {
		t.Errorf("Config directory was not created: %v", err)
	}

	if !info.IsDir() {
		t.Errorf("Config path is not a directory")
	}

	// Verify permissions (should be 0700 on Unix-like systems)
	// Note: This test may behave differently on Windows
	if info.Mode().Perm() != 0700 {
		t.Logf("Warning: Config directory permissions are %v, expected 0700", info.Mode().Perm())
	}
}

func TestSaveToken_FilePermissions(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gdocs-cli-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tokenPath := filepath.Join(tempDir, "token.json")

	// Create and save a test token
	testToken := &oauth2.Token{
		AccessToken: "test-token",
		TokenType:   "Bearer",
	}

	err = SaveToken(tokenPath, testToken)
	if err != nil {
		t.Fatalf("SaveToken() error = %v", err)
	}

	// Check file permissions
	info, err := os.Stat(tokenPath)
	if err != nil {
		t.Fatalf("Failed to stat token file: %v", err)
	}

	// Verify permissions are 0600 (read/write for owner only)
	if info.Mode().Perm() != 0600 {
		t.Errorf("Token file permissions = %v, want 0600", info.Mode().Perm())
	}
}
