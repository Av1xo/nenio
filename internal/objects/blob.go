package objects

import (
	"fmt"
	"os"
	"path/filepath"

	"lukechampine.com/blake3"
)


func CreateBlob(objectsDir string, content []byte) (string, error) {
	hash := blake3.Sum256(content)
	hashHex := fmt.Sprintf("%x", hash)

	if err := os.MkdirAll(objectsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create objects directory: %v", err)
	}

	blobPath := filepath.Join(objectsDir, hashHex)
	if err := os.WriteFile(blobPath, content, 0644); err != nil {
		return "", fmt.Errorf("failed to write blob: %v", err)
	}

	return hashHex, nil
}