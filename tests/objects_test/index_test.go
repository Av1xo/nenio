package objects_test

import (
	"nenio/internal/objects"
	"os"
	"testing"
)

func TestUpdateIndex(t *testing.T) {
	tempDir := t.TempDir()
	indexPath := tempDir + "/index"

	hash := "dummyhash"
	filePath := "test.txt"

	if err := objects.UpdateIndex(indexPath, hash, filePath); err != nil {
		t.Fatalf("Failed to update index: %v", err)
	}

	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("failed to read index: %v", err)
	}

	expected := hash + " " + filePath + "\n"
	if string(content) != expected {
		t.Errorf("Unexpected index content:\nGot: %s\nWant: %s", string(content), expected)
	}
}