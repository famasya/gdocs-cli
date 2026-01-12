package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestCLIHelp tests the --help flag
func TestCLIHelp(t *testing.T) {
	// Build the binary first
	buildCmd := exec.Command("go", "build", "-o", "gdocs-cli-test", ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}
	defer os.Remove("gdocs-cli-test")

	// Run with --help
	cmd := exec.Command("./gdocs-cli-test", "--help")
	output, err := cmd.CombinedOutput()

	// --help returns exit code 2 in Go's flag package
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 2 {
			t.Errorf("Expected exit code 2 for --help, got %d", exitErr.ExitCode())
		}
	}

	outputStr := string(output)
	expectedStrings := []string{
		"Usage of",
		"-url",
		"-config",
		"-init",
		"-clean",
		"Google Docs URL",
		"OAuth credentials JSON file",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("Help output missing expected string: %q\nGot: %s", expected, outputStr)
		}
	}
}

// TestCLIMissingFlags tests error handling for missing required flags
func TestCLIMissingFlags(t *testing.T) {
	buildCmd := exec.Command("go", "build", "-o", "gdocs-cli-test", ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}
	defer os.Remove("gdocs-cli-test")

	tests := []struct {
		name     string
		args     []string
		wantErr  string
		exitCode int
	}{
		{
			name:     "no flags",
			args:     []string{},
			wantErr:  "Error: --url flag is required",
			exitCode: 1,
		},
		{
			name:     "only --config",
			args:     []string{"--config=credentials.json"},
			wantErr:  "Error: --url flag is required",
			exitCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./gdocs-cli-test", tt.args...)
			output, err := cmd.CombinedOutput()

			if err == nil {
				t.Error("Expected error, got nil")
				return
			}

			exitErr, ok := err.(*exec.ExitError)
			if !ok {
				t.Errorf("Expected exec.ExitError, got %T", err)
				return
			}

			if exitErr.ExitCode() != tt.exitCode {
				t.Errorf("Expected exit code %d, got %d", tt.exitCode, exitErr.ExitCode())
			}

			outputStr := string(output)
			if !strings.Contains(outputStr, tt.wantErr) {
				t.Errorf("Expected error message containing %q, got: %s", tt.wantErr, outputStr)
			}
		})
	}
}

// TestCLIInvalidURL tests error handling for invalid URLs
func TestCLIInvalidURL(t *testing.T) {
	buildCmd := exec.Command("go", "build", "-o", "gdocs-cli-test", ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}
	defer os.Remove("gdocs-cli-test")

	// Create a dummy credentials file
	dummyCreds := "dummy-creds.json"
	if err := os.WriteFile(dummyCreds, []byte(`{"type": "authorized_user"}`), 0600); err != nil {
		t.Fatalf("Failed to create dummy credentials: %v", err)
	}
	defer os.Remove(dummyCreds)

	tests := []struct {
		name    string
		url     string
		wantErr string
	}{
		{
			name:    "invalid URL format",
			url:     "not-a-valid-url",
			wantErr: "invalid URL",
		},
		{
			name:    "non-Google Docs URL",
			url:     "https://example.com/document/123",
			wantErr: "invalid URL",
		},
		{
			name:    "Google Sheets URL",
			url:     "https://docs.google.com/spreadsheets/d/123/edit",
			wantErr: "invalid URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./gdocs-cli-test",
				"--url="+tt.url,
				"--config="+dummyCreds,
			)
			output, err := cmd.CombinedOutput()

			if err == nil {
				t.Error("Expected error, got nil")
				return
			}

			exitErr, ok := err.(*exec.ExitError)
			if !ok {
				t.Errorf("Expected exec.ExitError, got %T", err)
				return
			}

			if exitErr.ExitCode() != 1 {
				t.Errorf("Expected exit code 1, got %d", exitErr.ExitCode())
			}

			outputStr := string(output)
			if !strings.Contains(outputStr, tt.wantErr) {
				t.Errorf("Expected error message containing %q, got: %s", tt.wantErr, outputStr)
			}
		})
	}
}

// TestCLIMissingCredentialsFile tests error handling for missing credentials file
func TestCLIMissingCredentialsFile(t *testing.T) {
	buildCmd := exec.Command("go", "build", "-o", "gdocs-cli-test", ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}
	defer os.Remove("gdocs-cli-test")

	cmd := exec.Command("./gdocs-cli-test",
		"--url=https://docs.google.com/document/d/123abc/edit",
		"--config=/nonexistent/credentials.json",
	)
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected error for missing credentials file, got nil")
		return
	}

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Errorf("Expected exec.ExitError, got %T", err)
		return
	}

	if exitErr.ExitCode() != 1 {
		t.Errorf("Expected exit code 1, got %d", exitErr.ExitCode())
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "failed to read credentials file") {
		t.Errorf("Expected error about credentials file, got: %s", outputStr)
	}
}

// TestCLICleanFlag tests that the --clean flag is recognized
func TestCLICleanFlag(t *testing.T) {
	buildCmd := exec.Command("go", "build", "-o", "gdocs-cli-test", ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}
	defer os.Remove("gdocs-cli-test")

	// Test that --clean flag is accepted (even if other required flags are missing)
	cmd := exec.Command("./gdocs-cli-test", "--clean")
	output, err := cmd.CombinedOutput()

	// Should fail due to missing required flags, not unrecognized flag
	outputStr := string(output)
	if strings.Contains(outputStr, "flag provided but not defined: -clean") {
		t.Error("--clean flag not recognized")
	}

	// Should see the normal error about missing flags
	if !strings.Contains(outputStr, "Error: --url flag is required") {
		t.Errorf("Expected missing flags error, got: %s", outputStr)
	}

	if err == nil {
		t.Error("Expected error for missing flags")
	}
}
