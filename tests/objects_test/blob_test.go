package objects_test

import (
	"os"
	"testing"

	"nenio/internal/objects"
)

func TestCreateBlob(t *testing.T) {
	tempDir := t.TempDir()
	content := []byte("Hello, Nenio!")

	hash, err := objects.CreateBlob(tempDir, content)
	if err != nil {
		t.Fatalf("failed to create blob: %v", err)
	}

	blobPath := tempDir + "/" + hash
	if _, err := os.Stat(blobPath); os.IsNotExist(err) {
		t.Fatalf("Blob file was not created: %s", blobPath)
	}
}