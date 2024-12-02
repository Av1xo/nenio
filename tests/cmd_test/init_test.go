package cmd_test

import (
	"os"
	"testing"

	"nenio/cmd"
)

func TestInitializeRepo(t *testing.T) {
	// create temporary directory for testing
	tempDir := t.TempDir()
	defer os.RemoveAll(tempDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	if err := cmd.InitializeRepo(); err != nil {
		t.Fatalf(".nenio directory was not created")
	}

	requiredFiles := []string{
		".nenio/HEAD",
		".nenio/config",
		".nenio/refs/heads",
		".nenio/refs/tags",
		".nenio/objects",
	}

	for _, file := range requiredFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Fatalf("Missing required file/directory: %s", file)
		}
	}
}